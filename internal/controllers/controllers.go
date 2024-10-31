package controllers

import (
	"github.com/mikromolekula2002/Testovoe/internal/services"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	AuthController *AuthController
}

func NewController(services *services.Service, log *logrus.Logger) *Controller {
	return &Controller{
		AuthController: NewAuthController(services.AuthService, services.PostgreService, services.MailSendService, log),
	}
}
