package grpc_server

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
)

func IsTenantRequired(ctx context.Context, callMeta interceptors.CallMeta) bool {
	if callMeta.Method == "ResolveLink" {
		return false // 不需要解析多租户
	}
	return true
}
