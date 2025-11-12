package util

import (
	"github.com/google/uuid"
	"strings"
)

// ParseFullMethod 提取grpc服务名和grpc方法名
func ParseFullMethod(fullMethod string) (serviceName, methodName string) {
	// fullMethod 形如：/package.Service/Method
	parts := strings.Split(fullMethod, "/")
	if len(parts) == 3 {
		serviceName = parts[1] // package.Service
		methodName = parts[2]  // Method
	} else {
		serviceName = "unknown"
		methodName = "unknown"
	}
	return
}

func UUIDToRedisMember(id uuid.UUID) string {
	return id.String()
}

func RedisMemberToUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
