package service

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	MyErr "github.com/mikromolekula2002/Testovoe/internal/errors"
	"github.com/mikromolekula2002/Testovoe/internal/logger"
	"github.com/mikromolekula2002/Testovoe/internal/mocks"
	"github.com/mikromolekula2002/Testovoe/internal/models"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateTokens(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := logger.Init("debug", "", "stdout")

	mockRepo := mocks.NewMockTokenRepo(ctrl)
	mockJWT := mocks.NewMockJWTService(ctrl)
	mockSMTP := mocks.NewMockEmailSender(ctrl)

	serviceTest := ServiceInit(logger.Logrus,
		mockRepo,
		mockJWT,
		mockSMTP,
		[]byte("secretKey"),
		15,
		24)

	mockRepo.EXPECT().SaveRefreshToken(gomock.Any()).Return(nil)
	mockJWT.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("testToken", nil)

	access, refreshh, err := serviceTest.CreateTokens("testUser", "testIP")
	require.NoError(t, err)
	require.NotEmpty(t, access)
	require.NotEmpty(t, refreshh)
	require.Equal(t, "testToken", access, "Wrong actual accesToken, should be a `testToken`")
}

func TestCreateTokensWrongData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := logger.Init("debug", "", "stdout")
	mockRepo := mocks.NewMockTokenRepo(ctrl)
	mockJWT := mocks.NewMockJWTService(ctrl)
	mockSMTP := mocks.NewMockEmailSender(ctrl)

	serviceTest := ServiceInit(logger.Logrus,
		mockRepo,
		mockJWT,
		mockSMTP,
		[]byte("secretKey"),
		15,
		24)

	testTable := []struct {
		testName  string
		userID    string
		ipAddress string
	}{
		{
			testName:  "Wrong UserID",
			userID:    "",
			ipAddress: "testIP1",
		},
		{
			testName:  "Wrong IP Address",
			userID:    "user2",
			ipAddress: "",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.testName, func(t *testing.T) {
			_, _, err := serviceTest.CreateTokens(tt.userID, tt.ipAddress)
			require.Error(t, err)
		})
	}
}

func TestCreateTokensErrGenerateAccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := logger.Init("debug", "", "stdout")

	mockRepo := mocks.NewMockTokenRepo(ctrl)
	mockJWT := mocks.NewMockJWTService(ctrl)
	mockSMTP := mocks.NewMockEmailSender(ctrl)

	serviceTest := ServiceInit(logger.Logrus,
		mockRepo,
		mockJWT,
		mockSMTP,
		[]byte("secretKey"),
		15,
		24)

	expectedError := errors.New("test error")
	mockJWT.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", expectedError)

	access, refresh, err := serviceTest.CreateTokens("UserID", "IP")
	require.Error(t, err)
	require.Empty(t, access)
	require.Empty(t, refresh)
}

func TestCreateTokensErrSaveToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := logger.Init("debug", "", "stdout")

	mockRepo := mocks.NewMockTokenRepo(ctrl)
	mockJWT := mocks.NewMockJWTService(ctrl)
	mockSMTP := mocks.NewMockEmailSender(ctrl)

	serviceTest := ServiceInit(logger.Logrus,
		mockRepo,
		mockJWT,
		mockSMTP,
		[]byte("secretKey"),
		15,
		24)

	expectedError := errors.New("test error")
	mockJWT.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil)
	mockRepo.EXPECT().SaveRefreshToken(gomock.Any()).Return(expectedError)

	access, refreshh, err := serviceTest.CreateTokens("testUser", "testIP")
	require.Error(t, err)
	require.Empty(t, access)
	require.Empty(t, refreshh)
}

func TestRefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := logger.Init("debug", "", "stdout")

	mockRepo := mocks.NewMockTokenRepo(ctrl)
	mockJWT := mocks.NewMockJWTService(ctrl)
	mockSMTP := mocks.NewMockEmailSender(ctrl)

	serviceTest := ServiceInit(logger.Logrus,
		mockRepo,
		mockJWT,
		mockSMTP,
		[]byte("secretKey"),
		15,
		24)

	hashedToken, err := bcrypt.GenerateFromPassword([]byte("TestRefresh"), bcrypt.DefaultCost)
	if err != nil {
		logger.Logrus.Error("ошибка создания хеш пароля")
	}
	// Определите таблицу тестов
	testTable := struct {
		userID           string
		refreshToken     string
		ipAddress        string
		refreshTokenData *models.RefreshToken
	}{
		userID:       "UserID",
		refreshToken: "TestRefresh",
		ipAddress:    "testIP",
		refreshTokenData: &models.RefreshToken{
			ID:        24,
			UserID:    "UserID",
			TokenHash: string(hashedToken),
			Blocked:   false,
			CreatedAt: time.Now(),
			ExpiresAt: (time.Now().Add(time.Duration(serviceTest.refreshTokenDuration) * time.Hour)),
			IPAdress:  "testIP",
		},
	}

	mockRepo.EXPECT().GetRefreshToken(gomock.Any()).Return(testTable.refreshTokenData, nil)
	mockJWT.EXPECT().GenerateAccessToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("TestToken", nil)

	access, err := serviceTest.RefreshToken(testTable.userID, testTable.refreshToken, testTable.ipAddress)
	require.NoError(t, err)
	require.NotEmpty(t, access)
	require.Equal(t, "TestToken", access, "Access Token должен совпадать")
}

func TestRefreshTokenErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := logger.Init("debug", "", "stdout")

	mockRepo := mocks.NewMockTokenRepo(ctrl)
	mockJWT := mocks.NewMockJWTService(ctrl)
	mockSMTP := mocks.NewMockEmailSender(ctrl)

	serviceTest := ServiceInit(logger.Logrus,
		mockRepo,
		mockJWT,
		mockSMTP,
		[]byte("secretKey"),
		15,
		24)

	hashedToken, err := bcrypt.GenerateFromPassword([]byte("TestRefresh"), bcrypt.DefaultCost)
	if err != nil {
		logger.Logrus.Error("ошибка создания хеш пароля")
	}

	// Определите таблицу тестов
	testTable := []struct {
		testName         string
		userID           string
		refreshToken     string
		ipAddress        string
		refreshTokenData *models.RefreshToken
		expectedError    error
	}{
		{
			testName:     "Expired Token",
			userID:       "UserID",
			refreshToken: "TestRefresh",
			ipAddress:    "testIP",
			refreshTokenData: &models.RefreshToken{
				ID:        24,
				UserID:    "UserID",
				TokenHash: string(hashedToken),
				Blocked:   false,
				CreatedAt: time.Now(),
				ExpiresAt: time.Now().Add(-1 * time.Hour), // Устаревший токен
				IPAdress:  "testIP",
			},
			expectedError: MyErr.ErrExpiredToken,
		},
		{
			testName:     "Blocked Token",
			userID:       "UserID",
			refreshToken: "TestRefresh",
			ipAddress:    "testIP",
			refreshTokenData: &models.RefreshToken{
				ID:        24,
				UserID:    "UserID",
				TokenHash: string(hashedToken),
				Blocked:   true,
				CreatedAt: time.Now(),
				ExpiresAt: time.Now().Add(24 * time.Hour),
				IPAdress:  "testIP",
			},
			expectedError: MyErr.ErrBlockedToken,
		},
		{
			testName:     "Wrong Refresh Token",
			userID:       "UserID",
			refreshToken: "TestRefresh",
			ipAddress:    "testIP",
			refreshTokenData: &models.RefreshToken{
				ID:        24,
				UserID:    "UserID",
				TokenHash: "Wrong Refresh Token Test",
				Blocked:   false,
				CreatedAt: time.Now(),
				ExpiresAt: time.Now().Add(24 * time.Hour),
				IPAdress:  "testIP",
			},
			expectedError: MyErr.ErrDataToken,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.testName, func(t *testing.T) {
			mockRepo.EXPECT().GetRefreshToken(tt.userID).Return(tt.refreshTokenData, nil)

			if tt.expectedError == nil {
				mockJWT.EXPECT().GenerateAccessToken(tt.userID, tt.ipAddress, gomock.Any(), gomock.Any()).Return("testToken", nil)
			}

			actualToken, err := serviceTest.RefreshToken(tt.userID, tt.refreshToken, tt.ipAddress)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectedError, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, "testToken", actualToken)
			}
		})
	}
}

func TestCreateRefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := logger.Init("debug", "", "stdout")

	mockRepo := mocks.NewMockTokenRepo(ctrl)
	mockJWT := mocks.NewMockJWTService(ctrl)
	mockSMTP := mocks.NewMockEmailSender(ctrl)

	serviceTest := ServiceInit(logger.Logrus,
		mockRepo,
		mockJWT,
		mockSMTP,
		[]byte("secretKey"),
		15,
		24)

	_, _, err := serviceTest.CreateRefreshToken()
	require.NoError(t, err)
}
