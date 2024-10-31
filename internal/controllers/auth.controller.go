package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikromolekula2002/Testovoe/internal/models"
	"github.com/mikromolekula2002/Testovoe/internal/services"
	"github.com/sirupsen/logrus"
)

var (
	subject = "Email Warning"
	content = "Maybe your account was hacked, please check it"
	to      = []string{"kortymalik@gmail.com"}
)

type AuthController struct {
	authService     *services.AuthService
	postgreService  *services.PostgreService
	mailSendService *services.MailSendService
	log             *logrus.Logger
}

func NewAuthController(
	authService *services.AuthService,
	postgreService *services.PostgreService,
	mailSendService *services.MailSendService,
	log *logrus.Logger,
) *AuthController {
	return &AuthController{
		authService,
		postgreService,
		mailSendService,
		log,
	}
}

// @Summary Sending messages to recipients and save logs
// @Description Sending messages to recipients by smtp and saving the log of the sent mail
// @Tags mail
// @ID send-mail
// @Accept  application/json
// @Produce  application/json
// @Param subject body string true "Subject of the mail"
// @Param to body array true "Recipient mail addresses"
// @Param cc body array false "CC mail addresses"
// @Param bcc body array false "BCC mail addresses"
// @Param message body string true "Message content to be sent. Can be plain text or HTML-formatted text"
// @Success 200 {object} schemas.SendMailResponse
// @Router /mail [post]
func (mc *AuthController) CreateTokens(ctx *gin.Context) {
	// Получили с request параметр user uuid
	var req models.CreateTokensRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// получаем айпи адресс юзера
	clientIP := ctx.ClientIP()

	// создаем с помощью jwt новый access token с параметрами user id, exp time,
	accessToken, err := mc.authService.CreateAccessToken(req.UsedID, clientIP)
	if err != nil {
		fmt.Println() //
	}

	// создаем новый рандомный refresh token
	refreshToken := mc.authService.CreateRefreshToken(req.UsedID, clientIP)

	// отправляем рефреш токен в сервис с БД, хешируем рефреш токен и сохраняем вместе с айпи и user id, а также exp time
	err = mc.postgreService.SaveRefreshToken(refreshToken)
	if err != nil {
		fmt.Println("") //
	}

	// кодируем refresh token в BASE64
	tokenResponse := mc.authService.CreateResponse(accessToken, refreshToken.RefreshToken)
	//отправляем ответ юзеру с access token и refresh token

	ctx.JSON(http.StatusOK, tokenResponse)
}

func (mc *AuthController) RefreshTokens(ctx *gin.Context) {
	// ТУТ БУДЕТ ЛОГИКА ОБРАБОТКИ РОУТСОВ

	//получаем access и refresh tokens
	var req models.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Экстрактим access token из строки
	userID, err := mc.authService.ExtractAccessToken(req.AccessToken)
	if err != nil {
		fmt.Println("") //
	}

	// Получаем refresh token из БД и параллельно проверяем на совпадение по user id(GUID)
	refreshTokenData, err := mc.postgreService.GetResfreshToken(userID)
	if err != nil {
		fmt.Println("") //
	}

	// удаляем запись из БД нашего рефреш токена(сделано чтобы если зайдет хакер и проебется, мы сразу отключим такую сессию)
	err = mc.postgreService.DeleteRefreshToken(userID)
	if err != nil {
		fmt.Println("") //
	}

	// новое решение, будет метод в сервисе который получает реквест рефреш токен, декодирует его
	// возьмет рефреш токен из бд и все сверит с полученным от юзера, а сюда просто отдаст ошибку
	err = mc.authService.CheckRefreshToken(req.RefreshToken, userID, refreshTokenData)
	if err != nil {
		fmt.Println("") //
	}

	// получаем айпи адресс юзера
	clientIP := ctx.ClientIP()

	// проверка по айпи адрессу
	if clientIP != refreshTokenData.IPAddress {
		mc.mailSendService.SendMailWarning(subject, content, to)
	}

	accessToken, err := mc.authService.CreateAccessToken(userID, clientIP)
	if err != nil {
		fmt.Println("") //
	}

	refreshTokenData = mc.authService.CreateRefreshToken(userID, clientIP)

	// отправляем рефреш токен в сервис с БД, хешируем рефреш токен и сохраняем вместе с айпи и user id, а также exp time
	err = mc.postgreService.SaveRefreshToken(refreshTokenData)
	if err != nil {
		fmt.Println("") //
	}

	// отправляем ответ юзеру с новой парой токенов
	// кодируем refresh token в BASE64
	tokenResponse := mc.authService.CreateResponse(accessToken, refreshTokenData.RefreshToken)
	//отправляем ответ юзеру с access token и refresh token

	ctx.JSON(http.StatusOK, tokenResponse)
}
