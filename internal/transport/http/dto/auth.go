package httpdto

import "time"

const (
	AccessTokenDuration  = 1 * time.Hour
	RefreshTokenDuration = 7 * 24 * time.Hour
)

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
