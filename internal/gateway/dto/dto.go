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
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
