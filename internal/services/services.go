package services

import (
	"github.com/mikromolekula2002/Testovoe/internal/config"
	"github.com/mikromolekula2002/Testovoe/internal/repo"
)

type Service struct {
	AuthService     *AuthService
	PostgreService  *PostgreService
	MailSendService *MailSendService
}

func NewService(db repo.TokenRepo, newConfig *config.Config) *Service {
	return &Service{
		AuthService:     NewAuthService(newConfig),
		PostgreService:  NewPostgreService(db), //записать сюда пул подключения к БД
		MailSendService: NewMailSendService(newConfig),
	}
}
