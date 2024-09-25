package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/mikromolekula2002/Testovoe/internal/errors"
	"github.com/mikromolekula2002/Testovoe/internal/jwt_service"
	"github.com/mikromolekula2002/Testovoe/internal/mail_send"
	"github.com/mikromolekula2002/Testovoe/internal/models"
	"github.com/mikromolekula2002/Testovoe/internal/repo"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const (
	Subject = "Email Warning"
	Content = "Ваш токен пытались получить с другого ip адресса, если это были не вы, свяжитесь с нами."
)

var (
	EmailRecipient = []string{"kortymalik@gmail.com"} //вводите свою почту для проверки рассылки
)

// Основная структура для работы сервисного слоя
type Service struct {
	log                  *logrus.Logger
	repo                 repo.TokenRepo
	jwtService           jwt_service.JWTService
	smtp                 mail_send.EmailSender
	jwtKey               []byte
	accessTokenDuration  int
	refreshTokenDuration int
}

// Инит нашей структуры с базой данных, смпт сервером и сервисом jwt
func ServiceInit(logger *logrus.Logger,
	repo repo.TokenRepo,
	jwt jwt_service.JWTService,
	smtp mail_send.EmailSender,
	jwtkey []byte,
	accesstokenDuration int,
	refreshtokenDuration int) *Service {
	return &Service{

		log:                  logger,
		repo:                 repo,
		jwtService:           jwt,
		smtp:                 smtp,
		jwtKey:               jwtkey,
		accessTokenDuration:  accesstokenDuration,
		refreshTokenDuration: refreshtokenDuration,
	}
}

// метод создающий access token и refresh token, а также записывающий в БД refresh token
func (s *Service) CreateTokens(UserID string, ipAdress string) (string, string, error) {
	// проверка на пустого пользователя из параметров запроса
	if UserID == "" {
		s.log.Error("CreateTokens - пустой UserID")
		return "", "", errors.ErrMissingUserID
	}

	if ipAdress == "" {
		s.log.Error("CreateTokens - пустой IP Address")
		return "", "", errors.ErrWrongIP
	}

	// создание access token
	AccessToken, err := s.jwtService.GenerateAccessToken(UserID, ipAdress, s.jwtKey, s.accessTokenDuration)
	if err != nil {
		s.log.Error(err)
		return "", "", errors.ErrServer
	}

	// создание refresh token
	refreshToken, hashedToken, err := s.CreateRefreshToken()
	if err != nil {
		s.log.Error(err)
		return "", "", errors.ErrServer
	}

	// Сохранение хешированного Refresh токена в базу данных
	token := models.RefreshToken{
		UserID:      UserID,
		TokenHash:   hashedToken,
		AccessToken: AccessToken,
		IPAdress:    ipAdress,
		CreatedAt:   (time.Now()),
		ExpiresAt:   (time.Now().Add(time.Duration(s.refreshTokenDuration) * time.Hour)),
	}

	// сохраняем хеш resfresh token в базу данных postgreSQL
	err = s.repo.SaveRefreshToken(&token)
	if err != nil {
		s.log.Error(err)
		return "", "", errors.ErrServer
	}

	return AccessToken, refreshToken, nil
}

// Проверка refresh token и выдача нового access token
func (s *Service) RefreshToken(UserID string, RefreshToken string, IpAdress string, accessToken string) (string, error) {
	op := "service.RefreshToken"

	// проверка на пустого пользователя из параметров запроса
	if UserID == "" {
		return "", errors.ErrMissingUserID
	}

	// Поиск хешированного токена в базе данных по userID
	tokenStruct, err := s.repo.GetRefreshToken(UserID)
	if err != nil {
		s.log.Error(err)
		return "", errors.ErrMissingToken
	}

	claimsAccess, err := s.jwtService.CheckoutAccessToken(accessToken, s.jwtKey)
	if err != nil {
		return "", errors.ErrValidToken
	} else if claimsAccess.UserID != UserID {
		return "", errors.ErrValidToken
	}

	if tokenStruct.AccessToken != accessToken {
		return "", errors.ErrValidToken
	}

	if time.Now().After(tokenStruct.ExpiresAt) {
		s.log.Errorf("%s: Срок действия refresh token истек", op)
		return "", errors.ErrExpiredToken
	}

	if tokenStruct.Blocked {
		s.log.Errorf("%s: Рефреш токен заблокирован", op)
		return "", errors.ErrBlockedToken
	}

	// Сравнение полученного Refresh токена с сохраненным хешем
	err = bcrypt.CompareHashAndPassword([]byte(tokenStruct.TokenHash), []byte(RefreshToken))
	if err != nil {
		s.log.Errorf("%s: \n%v", op, err)
		return "", errors.ErrDataToken
	}

	// Проверка IP адреса, если ip адресс отличается отправка варна на почту
	if tokenStruct.IPAdress != IpAdress {
		err := s.smtp.SendEmail(Subject, Content, EmailRecipient)
		if err != nil {
			s.log.Error(err)
		}
	}

	// Генерация нового Access токена
	newAccessToken, err := s.jwtService.GenerateAccessToken(UserID, IpAdress, s.jwtKey, s.accessTokenDuration)
	if err != nil {
		s.log.Error(err)
		return "", errors.ErrServer
	}

	// Сохранение хешированного Refresh токена в базу данных
	token := models.RefreshToken{
		UserID:      UserID,
		AccessToken: newAccessToken,
		TokenHash:   tokenStruct.TokenHash,
		CreatedAt:   (time.Now()),
	}

	// сохраняем хеш resfresh token в базу данных postgreSQL
	err = s.repo.UpdateRefreshToken(&token)
	if err != nil {
		s.log.Error(err)
		return "", errors.ErrServer
	}

	return newAccessToken, nil
}

// Создание refresh Token
func (s *Service) CreateRefreshToken() (string, string, error) {
	op := "service.CreateRefreshToken"
	// Генерация случайного Refresh токена
	refreshToken := make([]byte, 32)
	if _, err := rand.Read(refreshToken); err != nil {
		return "", "", fmt.Errorf("%s: \n%v", op, err)
	}
	refreshTokenString := base64.URLEncoding.EncodeToString(refreshToken)

	// Хеширование Refresh токена
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshTokenString), bcrypt.DefaultCost)
	if err != nil {
		return "", "", fmt.Errorf("%s: \n%v", op, err)
	}

	return refreshTokenString, string(hashedToken), nil
}
