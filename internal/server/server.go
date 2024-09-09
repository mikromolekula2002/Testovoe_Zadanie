package server

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/mikromolekula2002/Testovoe/internal/service"
)

// Структура сервера с фреймворком и сервисным слоем
type Server struct {
	Echo *echo.Echo
	Svc  *service.Service
}

// Инициализация сервера
func Init(port string, Svc *service.Service) *Server {
	echo := echo.New()

	//
	echo.HideBanner = true
	//

	echo.Server = &http.Server{
		Addr:         port,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}

	return &Server{
		Echo: echo,
		Svc:  Svc,
	}
}

// Маршруты сервера
func (s *Server) Routes() {
	// Маршрут для проверки работы сервера
	s.Echo.POST("/getToken", s.GetTokens)
	s.Echo.POST("/refreshToken", s.RefreshToken)
}

// Обработчик маршрута /getToken
func (s *Server) GetTokens(c echo.Context) error {
	userID := c.QueryParam("user_id")
	ipAddress := c.RealIP()

	// Запуск сервисного слоя с обработкой токенов
	accessToken, refreshToken, err := s.Svc.CreateTokens(userID, ipAddress)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"access_token":  "nil",
			"refresh_token": "nil",
			"error":         err.Error()})
	}

	// Возвращаем токены
	return c.JSON(http.StatusOK, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"error":         "nil",
	})
}

// Обработчик маршрута /refreshToken
func (s *Server) RefreshToken(c echo.Context) error {
	refreshToken := c.FormValue("refresh_token")
	userID := c.FormValue("user_id")
	ipAddress := c.RealIP()

	// запуск сервисного слоя
	newAccessToken, err := s.Svc.RefreshToken(userID, refreshToken, ipAddress)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"access_token":  "nil",
			"refresh_token": "nil",
			"error":         err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"access_token": newAccessToken,
		"error":        "nil",
	})
}
