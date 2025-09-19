package schemas

type UserLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}