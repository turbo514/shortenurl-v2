package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/cespare/xxhash/v2"
	"github.com/google/uuid"
	"github.com/turbo514/shortenurl-v2/link/adapter"
	"github.com/turbo514/shortenurl-v2/link/usecase/base62"
	"github.com/turbo514/shortenurl-v2/shared/dto"
	"github.com/turbo514/shortenurl-v2/shared/mylog"
	mytrace "github.com/turbo514/shortenurl-v2/shared/trace"
	"github.com/turbo514/shortenurl-v2/shared/zerr"
	"go.opentelemetry.io/otel/codes"
	"time"

	"github.com/turbo514/shortenurl-v2/link/domain"
)

//type ILinkUseCase interface {
//}

type LinkUseCase struct {
	repo      adapter.IShortLinkRepository
	publisher adapter.IEventPublisher
}

func NewLinkUseCase(repo adapter.IShortLinkRepository, publisher adapter.IEventPublisher) *LinkUseCase {
	return &LinkUseCase{
		repo:      repo,
		publisher: publisher,
	}
}

func (uc *LinkUseCase) Resolve(ctx context.Context, req ResolveRequest) (string, error) {
	tr := mytrace.GetTracer()
	ctx, span := tr.Start(ctx, "LinkUseCase.Resolve")
	defer span.End()

	eventId, err := uuid.NewV7()
	if err != nil {
		span.SetStatus(codes.Error, "生成事件id失败")
		span.RecordError(err)
		return "", fmt.Errorf("生成事件id失败: %w", err)
	}

	// 根据短链码获取短链
	link, err := uc.repo.FindByCode(ctx, req.Code)
	if err != nil {
		return "", err
	}

	// 构造点击事件并发送到消息队列
	clickEvent := dto.ClickEvent{
		EventId:   eventId,
		LinkID:    link.ID,
		TenantID:  link.TenantID,
		ClickTime: req.ClickTime,
		IpAddress: req.IpAddress,
		UserAgent: req.UserAgent,
		Referrer:  req.Referrer,
	}
	if err := uc.publisher.PublishClickEvent(ctx, clickEvent); err != nil {
		// TODO: 收件箱模式
		// TODO: 异步发送+连接池
		mylog.GetLogger().Warn("发送点击事件失败", "err", err.Error(), "clickEvent", clickEvent)
	}
	return link.OriginalURL, nil
}

func (uc *LinkUseCase) Shorten(ctx context.Context, originalUrl string, tenantID, userID string, expiration time.Duration) (*domain.ShortLink, error) {
	tr := mytrace.GetTracer()
	ctx, span := tr.Start(ctx, "LinkUseCase.Shorten")
	defer span.End()

	eventId, err := uuid.NewV7()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "生成事件id失败")
		return nil, fmt.Errorf("生成事件id失败: %w", err)
	}

	// 参数解析,获取租户ID,用户ID
	tenantId, err := uuid.Parse(tenantID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "解析TenantID失败")
		return nil, fmt.Errorf("解析TenantID失败: %w", err)
	}
	userId, err := uuid.Parse(userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "解析UserID失败")
		return nil, fmt.Errorf("解析UserID失败: %w", err)
	}

	var expiresAt *time.Time = nil // 如果expiration小于等于0,则该短链永不过期
	if expiration > 0 {
		tmp := time.Now().UTC().Add(expiration)
		expiresAt = &tmp
	}

	shortLink := &domain.ShortLink{
		OriginalURL: originalUrl,
		TenantID:    tenantId,
		UserID:      userId,
		ExpireAt:    expiresAt,
	}

	// 创建短链
	success := false
	for retryTimes := 0; retryTimes < 3; retryTimes++ {
		shortLink.ID, err = uuid.NewV7()
		if err != nil {
			continue
		}
		code := uc.generateShortCode(shortLink.ID)
		shortLink.ShortCode = code

		if err := uc.repo.CreateLink(ctx, shortLink); err != nil {
			span.RecordError(err)
			if !errors.Is(err, zerr.ErrDuplicateEntry) {
				span.SetStatus(codes.Error, "创建短链失败")
				return nil, fmt.Errorf("创建短链接失败: %w", err)
			}
		} else {
			success = true
			break
		}
	}
	if !success {
		span.SetStatus(codes.Error, "创建失败次数过多")
		return nil, fmt.Errorf("创建失败次数过多")
	}

	// 发送短链创建事件到消息队列
	createEvent := dto.CreateEvent{
		EventId:     eventId,
		TenantID:    tenantId,
		LinkID:      shortLink.ID,
		OriginalURL: originalUrl,
	}
	if err := uc.publisher.PublishCreateEvent(ctx, createEvent); err != nil {
		// TODO: 收件箱模式
		// TODO: 异步发送+连接池
		mylog.GetLogger().Warn("发送CreateLink事件失败", "err", err.Error(), "createEvent", createEvent)
	}

	return shortLink, nil
}

func (uc *LinkUseCase) generateShortCode(id uuid.UUID) string {
	hashed := xxhash.Sum64(id[:])
	code := base62.EncodeBase62(hashed)
	return code
}
