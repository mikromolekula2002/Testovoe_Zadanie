package services

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/mikromolekula2002/Testovoe/internal/config"
	"github.com/mikromolekula2002/Testovoe/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// JWT claims с IPAdress и UserID
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type AuthService struct {
	accessTokenTimeDuration  int
	refreshTokenTimeDuration int
	secretKey                string
}

func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{
		accessTokenTimeDuration:  cfg.Token.AccessTokenTimeDuration,
		refreshTokenTimeDuration: cfg.Token.RefreshTokenTimeDuration,
		secretKey:                cfg.Jwt.SecretKey,
	}
}

// тут переименовать и написать логику
func (ms *AuthService) CreateAccessToken(userID string, ipAddress string) (string, error) {
	op := "jwt_service.GenerateAccessToken"

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(ms.accessTokenTimeDuration) * time.Minute)), // Токен истекает через 15 минут
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	accessToken, err := token.SignedString(ms.secretKey)
	if err != nil {
		return "", fmt.Errorf("%s: \n%v", op, err)
	}
	return accessToken, nil
}

// Создание refresh Token
func (ms *AuthService) CreateRefreshToken(userID, ipAddress string) *models.RefreshTokenData {
	refreshToken := uuid.New()

	refreshTokenData := &models.RefreshTokenData{
		UserID:       userID,
		RefreshToken: refreshToken.String(),
		ExpiresAt:    (time.Now().Add(time.Duration(ms.refreshTokenTimeDuration) * time.Minute)),
		IPAddress:    ipAddress,
	}

	return refreshTokenData
}

func (ms *AuthService) CreateResponse(accessToken, refreshToken string) *models.CreateTokenResponse {
	refreshTokenBytes := base64.StdEncoding.EncodeToString([]byte(refreshToken))

	resp := &models.CreateTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenBytes,
	}

	return resp
}

// проверка refresh token
func (ms *AuthService) CheckRefreshToken(
	refreshTokenReq string,
	userID string,
	refreshTokenData *models.RefreshTokenData) error {

	op := "service.RefreshToken"

	if time.Now().After(refreshTokenData.ExpiresAt) {
		return fmt.Errorf("%s: \n%v", op, "Expired refresh token")
	}

	decodedRefreshTokenReq, err := base64.StdEncoding.DecodeString(refreshTokenReq)
	if err != nil {
		return fmt.Errorf("%s: \n%v", op, err)
	}

	if userID != refreshTokenData.UserID {
		return fmt.Errorf("%s: \n%v", op, "Invalid User ID")
	}

	err = bcrypt.CompareHashAndPassword([]byte(refreshTokenData.RefreshToken), []byte(decodedRefreshTokenReq))
	if err != nil {
		return fmt.Errorf("%s: \n%v", op, err)
	}

	return nil
}

func (ms *AuthService) ExtractAccessToken(accessToken string) (string, error) {
	op := "jwt_service.CheckoutAccessToken"

	token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return ms.secretKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("%s: \n%v", op, err)
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return "", fmt.Errorf("%s: \n%v", op, err)
	}
	return claims.UserID, nil
}
