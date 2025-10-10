package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/cespare/xxhash/v2"
	"github.com/google/uuid"
	"github.com/turbo514/shortenurl-v2/link/usecase/base62"
	"github.com/turbo514/shortenurl-v2/shared/events"
	mytrace "github.com/turbo514/shortenurl-v2/shared/trace"
	"github.com/turbo514/shortenurl-v2/shared/zerr"
	"go.opentelemetry.io/otel/codes"
	"log/slog"
	"time"

	"github.com/turbo514/shortenurl-v2/link/entity"
)

type IShortLinkRepository interface {
	FindByCode(ctx context.Context, code string) (*entity.ShortLink, error)
	CreateLink(ctx context.Context, shortLink *entity.ShortLink) error
}

type IEventPublisher interface {
	PublishClickEvent(ctx context.Context, event events.ClickEvent) error
	PublishCreateEvent(ctx context.Context, event events.CreateEvent) error
}

//type ILinkUseCase interface {
//}

type LinkUseCase struct {
	repo      IShortLinkRepository
	publisher IEventPublisher
}

func NewLinkUseCase(repo IShortLinkRepository, publisher IEventPublisher) *LinkUseCase {
	return &LinkUseCase{
		repo:      repo,
		publisher: publisher,
	}
}

func (uc *LinkUseCase) Resolve(ctx context.Context, req ResolveRequest) (string, error) {
	tr := mytrace.GetTracer()
	ctx, span := tr.Start(ctx, "LinkUseCase.Resolve")
	defer span.End()

	ctx, repoSpan := tr.Start(ctx, "LinkRepo.FindByCode")
	link, err := uc.repo.FindByCode(ctx, req.Code)
	repoSpan.End()

	if err != nil {
		span.SetStatus(codes.Error, "查找短链接失败")
		return "", fmt.Errorf("查找短链接失败: %w", err)
	}

	// 播报点击事件
	clickEvent := events.ClickEvent{
		LinkID:    link.ID,
		TenantID:  link.TenantID,
		ClickTime: req.ClickTime,
		IpAddress: req.IpAddress,
		UserAgent: req.UserAgent,
		Referrer:  req.Referrer,
	}

	ctx, publishSpan := tr.Start(ctx, "LinkUseCase.PublishClickEvent")
	defer publishSpan.End()

	if err := uc.publisher.PublishClickEvent(ctx, clickEvent); err != nil {
		// TODO: 收件箱模式
		// TODO: 异步发送+连接池
		publishSpan.SetStatus(codes.Error, "发送点击事件失败")
		slog.Warn("发送点击事件失败", "err", err.Error(), "clickEvent", clickEvent)
	}
	return link.OriginalURL, nil
}

func (uc *LinkUseCase) Shorten(ctx context.Context, originalUrl string, tenantID, userID string, expiration time.Duration) (*entity.ShortLink, error) {
	tr := mytrace.GetTracer()
	ctx, span := tr.Start(ctx, "LinkUseCase.Shorten")
	defer span.End()

	// 参数解析
	tenantId, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, fmt.Errorf("解析TenantID失败: %w", err)
	}
	userId, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("解析UserID失败: %w", err)
	}

	var expiresAt *time.Time = nil
	if expiration > 0 {
		tmp := time.Now().UTC().Add(expiration)
		expiresAt = &tmp
	}

	shortLink := &entity.ShortLink{
		OriginalURL: originalUrl,
		TenantID:    tenantId,
		UserID:      userId,
		ExpireAt:    expiresAt,
	}

	// 开始创建短链
	ctx, repoSpan := tr.Start(ctx, "LinkRepo.Shorten")
	defer repoSpan.End()

	success := false
	for retryTimes := 0; retryTimes < 3; retryTimes++ {
		shortLink.ID, err = uuid.NewV7()
		if err != nil {
			continue
		}
		code := uc.generateShortCode(ctx, shortLink.ID[:])
		shortLink.ShortCode = code

		if err := uc.repo.CreateLink(ctx, shortLink); err != nil {
			if !errors.Is(err, zerr.ErrDuplicateEntryDB) {
				repoSpan.SetStatus(codes.Error, "插入数据库失败")
				return nil, fmt.Errorf("创建短链接失败: %w", err)
			}
		} else {
			success = true
			break
		}
	}
	if !success {
		repoSpan.SetStatus(codes.Error, "创建失败次数过多")
		return nil, fmt.Errorf("创建失败次数过多")
	}

	// 发送到消息队列,让数据分析库保留一份
	createEvent := events.CreateEvent{
		TenantID:    tenantId,
		LinkID:      shortLink.ID,
		OriginalURL: originalUrl,
	}

	ctx, publishSpan := tr.Start(ctx, "LinkUseCase.PublishCreateEvent")
	defer publishSpan.End()

	if err := uc.publisher.PublishCreateEvent(ctx, createEvent); err != nil {
		// TODO: 收件箱模式
		// TODO: 异步发送+连接池
		publishSpan.SetStatus(codes.Error, "发送CreateLink事件失败")
		slog.Warn("发送CreateLink事件失败", "err", err.Error(), "createEvent", createEvent)
	}

	return shortLink, nil
}

func (uc *LinkUseCase) generateShortCode(ctx context.Context, input []byte) string {
	hashed := xxhash.Sum64(input)
	code := base62.EncodeBase62(hashed)
	return code
}
