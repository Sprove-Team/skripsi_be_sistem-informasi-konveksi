package auth

type Login struct {
	Username string `json:"username" validate:"required,printascii"`
	Password string `json:"password" validate:"required,printascii"`
}

type GetNewAccessToken struct {
	RefreshToken string `json:"refresh_token" validate:"required,jwt"`
}
