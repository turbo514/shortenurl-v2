package client

import (
	analyticspb "github.com/turbo514/shortenurl-v2/shared/gen/proto/analytics"
	linkpb "github.com/turbo514/shortenurl-v2/shared/gen/proto/link"
	tenantpb "github.com/turbo514/shortenurl-v2/shared/gen/proto/tenant"
	"google.golang.org/grpc"
)

type LinkConnection *grpc.ClientConn
type TenantConnection *grpc.ClientConn
type AnalyticsConnection *grpc.ClientConn

func NewLinkConn(addr string, opts ...grpc.DialOption) (LinkConnection, error) {
	conn, err := grpc.NewClient(addr, opts...)
	return LinkConnection(conn), err
}

func NewTenantConn(addr string, opts ...grpc.DialOption) (TenantConnection, error) {
	conn, err := grpc.NewClient(addr, opts...)
	return TenantConnection(conn), err
}

func NewAnalyticsConn(addr string, opts ...grpc.DialOption) (AnalyticsConnection, error) {
	conn, err := grpc.NewClient(addr, opts...)
	return AnalyticsConnection(conn), err
}

func NewLinkClient(connection LinkConnection) linkpb.LinkServiceClient {
	conn := (*grpc.ClientConn)(connection)
	client := linkpb.NewLinkServiceClient(conn)
	return client
}

func NewTenantClient(connection TenantConnection) tenantpb.TenantServiceClient {
	conn := (*grpc.ClientConn)(connection)
	client := tenantpb.NewTenantServiceClient(conn)
	return client
}

func NewAnalyticsClient(connection AnalyticsConnection) analyticspb.AnalyticsServiceClient {
	conn := (*grpc.ClientConn)(connection)
	client := analyticspb.NewAnalyticsServiceClient(conn)
	return client
}
