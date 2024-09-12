package main

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/mikromolekula2002/Testovoe/internal/config"
	"github.com/mikromolekula2002/Testovoe/internal/jwt_service"
	"github.com/mikromolekula2002/Testovoe/internal/logger"
	"github.com/mikromolekula2002/Testovoe/internal/mail_send"
	"github.com/mikromolekula2002/Testovoe/internal/repo"
	"github.com/mikromolekula2002/Testovoe/internal/server"
	"github.com/mikromolekula2002/Testovoe/internal/service"
)

func main() {
	//Загрузка конфига
	cfg := config.LoadConfig("./config/config.yaml")
	// Инициализация базы данных
	logger := logger.Init(cfg.Logger.Level, cfg.Logger.FilePath, cfg.Logger.Output)

	//инициализация репозитория
	repo, err := repo.InitDB(cfg)
	if err != nil {
		log.Fatalf("Ошибка инициализации Базы Данных: \n%v", err)
	}

	//инициализация смпт сервера для отправки Email Warning
	smtp := mail_send.NewGmailSender(cfg.SMTP.Name,
		cfg.SMTP.Email_address,
		cfg.SMTP.Email_password,
		cfg.SMTP.Server_address,
		cfg.SMTP.Auth_address)

	//инициализация jwt сервиса
	jwtService := jwt_service.InitJWT()

	//инициализация сервисного слоя
	service := service.ServiceInit(logger.Logrus,
		repo,
		jwtService,
		smtp,
		[]byte(cfg.Jwt.JwtKey),
		cfg.Token.AccessTokenDuration,
		cfg.Token.RefreshTokenDuration)

	// инициализация сервера
	server := server.Init(cfg.Server.Port, service)
	// обработка енд поинтов
	server.Routes()

	// Запуск сервера на порту 8080
	if err := server.Echo.Start(cfg.Server.Port); err != nil {
		log.Fatal("Ошибка при запуске сервера:", err)
	}
}
