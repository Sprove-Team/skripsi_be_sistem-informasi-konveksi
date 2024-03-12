package req_auth

type Login struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type GetNewToken struct {
	RefreshToken string `json:"refresh_token" validate:"required,jwt"`
}
