package dto

type CreateTenantRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateTenantResponse struct {
	TenantId string `json:"tenant_id"`
	ApiKey   string `json:"api_key"`
}

type CreateUserRequest struct {
	TenantId string `json:"tenant_id" binding:"required"`
	Apikey   string `json:"api_key" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateUserResponse struct{}

type LoginRequest struct {
	Name     string `json:"name" binding:"required"`
	TenantId string `json:"tenant_id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type CreateLinkRequest struct {
	OriginalUrl string `json:"original_url" binding:"required"`
	Expiration  int64  `json:"expiration" binding:"required"`
}

type CreateLinkResponse struct {
	OriginalUrl string `json:"original_url"`
	ShortCode   string `json:"short_code"`
}

type ResolveLinkRequest struct {
	ShortCode string `json:"short_code" binding:"required" uri:"short_code"`
}

type ResolveLinkResponse struct {
	OriginalUrl string `json:"original_url"`
}

type GetTopLinksTodayRequest struct {
	Num int64 `form:"num"`
}

type TopLinks struct {
	ID          string `json:"id"`
	OriginalUrl string `json:"original_url"`
	ClickTimes  int64  `json:"click_times"`
}
type GetTopLinksTodayResponse struct {
	TopLinks TopLinks `json:"top_links"`
}
