package models

import "time"

type RefreshTokenData struct {
	UserID       string    `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	IPAddress    string    `json:"ip_adress"`
}

type CreateTokensRequest struct {
	UsedID string `json:"user_id"`
}

type RefreshTokenRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type CreateTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
