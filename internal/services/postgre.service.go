package services

import (
	"github.com/mikromolekula2002/Testovoe/internal/models"
)

type TokenRepo interface {
	SaveRefreshToken(refreshToken *models.RefreshTokenData) error
	GetRefreshToken(UserID string) (*models.RefreshTokenData, error)
	DeleteRefreshToken(userID string) error
}

type PostgreService struct {
	TokenRepository TokenRepo
}

func NewPostgreService(tokenRepository TokenRepo) *PostgreService {
	return &PostgreService{tokenRepository}
}

// тут переименовать и написать логику
func (ms *PostgreService) SaveRefreshToken(refreshToken *models.RefreshTokenData) error {
	// связь с БД
	return nil
}

// тут переименовать и написать логику
func (ms *PostgreService) GetResfreshToken(userID string) (*models.RefreshTokenData, error) {
	// связь с БД
	return nil, nil
}

// тут переименовать и написать логику
func (ms *PostgreService) DeleteRefreshToken(userID string) error {
	// связь с БД
	return nil
}
