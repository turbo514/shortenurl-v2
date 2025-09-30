package migrate_helper

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"path/filepath"
	"strings"
)

func Up(relativeUrl string, databaseURL string) error {
	migrationPath, err := filepath.Abs(relativeUrl)
	if err != nil {
		return fmt.Errorf("解析url失败: %w", err)
	}
	migrationPath = strings.Replace(migrationPath, "\\", "/", -1)

	m, err := migrate.New("file://"+migrationPath, databaseURL)
	if err != nil {
		return fmt.Errorf("迁移初始化失败: %w", err)
	}
	defer m.Close()
	
	err = m.Up()
	if err != nil {
		return fmt.Errorf("迁移失败: %w", err)
	}

	return nil
}
