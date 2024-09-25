package jwt_service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct{}

type JWTService interface {
	GenerateAccessToken(userID, ipAddress string, secretKey []byte, timeDuration int) (string, error)
	CheckoutAccessToken(accessToken string, secretKey []byte) (*Claims, error)
}

// JWT claims с IPAdress и UserID
type Claims struct {
	UserID    string `json:"user_id"`
	IPAddress string `json:"ip_address"`
	jwt.RegisteredClaims
}

// инициализация сервиса JWT
func InitJWT() *JWTManager {
	return &JWTManager{} // Возвращаем конкретную реализацию интерфейса
}

// generateAccessToken создает новый Access токен
func (j *JWTManager) GenerateAccessToken(userID, ipAddress string, secretKey []byte, timeDuration int) (string, error) {
	op := "jwt_service.GenerateAccessToken"

	claims := &Claims{
		UserID:    userID,
		IPAddress: ipAddress,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(timeDuration) * time.Minute)), // Токен истекает через 15 минут
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	accessToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("%s: \n%v", op, err)
	}
	return accessToken, nil
}

func (j *JWTManager) CheckoutAccessToken(accessToken string, secretKey []byte) (*Claims, error) {
	op := "jwt_service.CheckoutAccessToken"

	token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: \n%v", op, err)
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("%s: \n%v", op, err)
	}
	return claims, nil
}
