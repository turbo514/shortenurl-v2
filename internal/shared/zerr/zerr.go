package zerr

import "errors"

var (
	ErrNotFoundDB        = errors.New("db not found")        // 请求的资源不存在
	ErrDuplicateEntry    = errors.New("db duplicate entry")  // 违反唯一性约束
	ErrOperationFailedDB = errors.New("db operation failed") // 未分类的数据库错误
)

var (
	ErrNotFoundCache = errors.New("cache not found")
)

var (
	ErrNotExist = errors.New("not exist")
)
