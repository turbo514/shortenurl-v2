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
	"github.com/turbo514/shortenurl-v2/shared/zerr"
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
	now := time.Now()
	result, err := m.queries.GetOriginalUrlByCode(ctx, model.GetOriginalUrlByCodeParams{
		ShortCode: code,
		ExpiresAt: &now,
	})
	if err != nil {
		return nil, err
	}

	id, err := uuid.FromBytes(result.ID)
	if err != nil {
		return nil, fmt.Errorf("id parse error: %s", err.Error())
	}
	tenantId, err := uuid.FromBytes(result.TenantID)
	if err != nil {
		return nil, fmt.Errorf("tenantId parse error: %s", err.Error())
	}
	userId, err := uuid.FromBytes(result.UserID)
	if err != nil {
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
