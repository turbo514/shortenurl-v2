package mysql_repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/turbo514/shortenurl-v2/link/adapter"
	"github.com/turbo514/shortenurl-v2/link/domain"
	model "github.com/turbo514/shortenurl-v2/link/infrastructure/mysql_repository/model"
	"github.com/turbo514/shortenurl-v2/link/metrics"
	mytrace "github.com/turbo514/shortenurl-v2/shared/trace"
	"github.com/turbo514/shortenurl-v2/shared/zerr"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"go.opentelemetry.io/otel/trace"
	"time"
)

var _ adapter.IShortLinkRepository = (*MysqlShortLinkRepository)(nil)

type MysqlShortLinkRepository struct {
	queries *model.Queries
}

func NewMysqlShortLinkDB(queries *model.Queries) *MysqlShortLinkRepository {
	return &MysqlShortLinkRepository{queries: queries}
}

func (m *MysqlShortLinkRepository) FindByCode(ctx context.Context, code string) (*domain.ShortLink, error) {
	tr := mytrace.GetTracer()
	ctx, span := tr.Start(ctx, "MysqlShortLinkRepository.FindByCode")
	defer span.End()

	adjectSpan(span, "SELECT", model.GetOriginalUrlByCode)
	span.SetAttributes(
		attribute.String("code", code),
	)

	// 根据code来查找对应短链
	metrics.AddDbOperationsTotalQuery()
	//g := singleflight.Group{}
	start := time.Now()
	result, err := m.queries.GetOriginalUrlByCode(ctx, model.GetOriginalUrlByCodeParams{
		ShortCode: code,
		ExpiresAt: &start,
	})
	end := time.Now()
	metrics.ObserveDbDurationSecondsQuery(end.Sub(start))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, zerr.ErrNotExist
		} else {
			span.RecordError(err)
			span.SetStatus(codes.Error, "数据库查询失败")
			return nil, fmt.Errorf("failed to query FindByCode(%s): %w", code, err)
		}
	}

	// 解析为领域实体
	id, err := uuid.FromBytes(result.ID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "id解析失败")
		return nil, fmt.Errorf("mysql FindByCode parse id error: %w", err)
	}
	tenantId, err := uuid.FromBytes(result.TenantID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "tenantId解析失败")
		return nil, fmt.Errorf("mysql parse tenantId error: %w", err)
	}
	userId, err := uuid.FromBytes(result.UserID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "userId解析失败")
		return nil, fmt.Errorf("mysql FindByCode parse userId error: %w", err)
	}

	return &domain.ShortLink{
		ID:          id,
		ShortCode:   code,
		TenantID:    tenantId,
		OriginalURL: result.OriginalUrl,
		UserID:      userId,
		ExpireAt:    result.ExpiresAt,
	}, nil
}

func (m *MysqlShortLinkRepository) CreateLink(ctx context.Context, shortLink *domain.ShortLink) error {
	tr := mytrace.GetTracer()
	ctx, span := tr.Start(ctx, "MysqlShortLinkRepository.CreateLink")
	defer span.End()

	adjectSpan(span, "INSERT", model.CreateShortLink)

	// TODO: 添加attribute

	// 向数据库插入对应短链
	metrics.AddDbOperationsTotalInsert()
	start := time.Now()
	err := m.queries.CreateShortLink(ctx, model.CreateShortLinkParams{
		ID:          shortLink.ID[:],
		TenantID:    shortLink.TenantID[:],
		UserID:      shortLink.UserID[:],
		ShortCode:   shortLink.ShortCode,
		OriginalUrl: shortLink.OriginalURL,
		CreatedAt:   start,
		ExpiresAt:   shortLink.ExpireAt,
	})
	end := time.Now()
	metrics.ObserveDbDurationSecondsInsert(end.Sub(start))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "数据库插入失败")
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return zerr.ErrDuplicateEntry
		}
		return fmt.Errorf("failed to execute insert CreateLink: %w", err)
	}
	return nil
}

func adjectSpan(s trace.Span, operation, statement string) {
	s.SetAttributes(
		semconv.DBSystemMySQL,
		semconv.DBName("link_dev"),
		semconv.DBOperation(operation),
		semconv.DBStatement(statement),
	)
}
