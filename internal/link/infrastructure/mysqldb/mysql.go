package mysqldb

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/turbo514/shortenurl-v2/link/adapter/db"
	"github.com/turbo514/shortenurl-v2/link/entity"
	model "github.com/turbo514/shortenurl-v2/link/infrastructure/mysqldb/model"
	mytrace "github.com/turbo514/shortenurl-v2/shared/trace"
	"github.com/turbo514/shortenurl-v2/shared/zerr"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"go.opentelemetry.io/otel/trace"
	"time"
)

var _ db.IShortLinkDB = (*MysqlShortLinkDB)(nil)

// 一定程度上的反模式无伤大雅(可能并非反模式)

type MysqlShortLinkDB struct {
	queries *model.Queries
}

func NewMysqlShortLinkDB(queries *model.Queries) *MysqlShortLinkDB {
	return &MysqlShortLinkDB{queries: queries}
}

func (m *MysqlShortLinkDB) FindByCode(ctx context.Context, code string) (*entity.ShortLink, error) {
	tr := mytrace.GetTracer()
	ctx, span := tr.Start(ctx, "MySQL.FindByCode")
	defer span.End()

	adjectSpan(span, "SELECT", model.GetOriginalUrlByCode)
	span.SetAttributes(
		attribute.String("code", code),
	)

	now := time.Now()
	result, err := m.queries.GetOriginalUrlByCode(ctx, model.GetOriginalUrlByCodeParams{
		ShortCode: code,
		ExpiresAt: &now,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "数据库查询失败")
		return nil, err
	}

	id, err := uuid.FromBytes(result.ID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "id解析失败")
		return nil, fmt.Errorf("id parse error: %s", err.Error())
	}
	tenantId, err := uuid.FromBytes(result.TenantID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "tenantId解析失败")
		return nil, fmt.Errorf("tenantId parse error: %s", err.Error())
	}
	userId, err := uuid.FromBytes(result.UserID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "userId解析失败")
		return nil, fmt.Errorf("userId parse error: %s", err.Error())
	}

	return &entity.ShortLink{
		ID:          id,
		ShortCode:   code,
		TenantID:    tenantId,
		OriginalURL: result.OriginalUrl,
		UserID:      userId,
		ExpireAt:    result.ExpiresAt,
	}, nil
}

func (m *MysqlShortLinkDB) CreateLink(ctx context.Context, shortLink *entity.ShortLink) error {
	tr := mytrace.GetTracer()
	ctx, span := tr.Start(ctx, "MySQL.CreateLink")
	defer span.End()

	adjectSpan(span, "INSERT", model.CreateShortLink)

	// TODO: 添加attribute

	err := m.queries.CreateShortLink(ctx, model.CreateShortLinkParams{
		ID:          shortLink.ID[:],
		TenantID:    shortLink.TenantID[:],
		UserID:      shortLink.UserID[:],
		ShortCode:   shortLink.ShortCode,
		OriginalUrl: shortLink.OriginalURL,
		CreatedAt:   time.Now(),
		ExpiresAt:   shortLink.ExpireAt,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "数据库插入失败")
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 { // 违反唯一性约束
				return fmt.Errorf("%w: %w", zerr.ErrDuplicateEntryDB, err)
			}
			return err
		}
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
