package appcontext

import (
	analyticspb "github.com/turbo514/shortenurl-v2/shared/gen/proto/analytics"
	linkpb "github.com/turbo514/shortenurl-v2/shared/gen/proto/link"
	tenantpb "github.com/turbo514/shortenurl-v2/shared/gen/proto/tenant"
)

type Services struct {
	Link      linkpb.LinkServiceClient
	Tenant    tenantpb.TenantServiceClient
	Analytics analyticspb.AnalyticsServiceClient
}

func NewServices(link linkpb.LinkServiceClient, tenant tenantpb.TenantServiceClient, analytics analyticspb.AnalyticsServiceClient) *Services {
	return &Services{
		Link:      link,
		Tenant:    tenant,
		Analytics: analytics,
	}
}
