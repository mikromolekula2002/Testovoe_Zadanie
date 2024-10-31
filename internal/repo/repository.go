package repo

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mikromolekula2002/Testovoe/internal/config"
	"github.com/mikromolekula2002/Testovoe/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type RepoPostgre struct {
	DB *gorm.DB
}

type TokenRepo interface {
	SaveRefreshToken(refreshToken *models.RefreshTokenData) error
	GetRefreshToken(UserID string) (*models.RefreshTokenData, error)
	DeleteRefreshToken(userID string) error
}

type RefreshToken struct {
	UserID           uuid.UUID `gorm:"type:uuid;not null"`                // UUID для уникального идентификатора пользователя
	RefreshTokenHash string    `gorm:"type:varchar(255);not null;unique"` // Хешированный refresh token
	CreatedAt        time.Time `gorm:"autoCreateTime"`                    // Время создания записи
	ExpiresAt        time.Time `gorm:"not null"`                          // Время истечения токена
	IPAdress         string    `gorm:"type:varchar(45);not null"`         // IP адрес пользователя (45 символов для IPv6)
}

// initDB инициализирует соединение с базой данных PostgreSQL
func InitDB(config *config.Config) (*RepoPostgre, error) {
	// Конфигурация для подключения к базе данных PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		config.Database.Host,
		config.Database.User,
		config.Database.Password,
		config.Database.DBName,
		config.Database.Port,
		config.Database.Sslmode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&RefreshToken{})
	if err != nil {
		return nil, err
	}

	database := RepoPostgre{
		DB: db,
	}
	return &database, nil
}

// Сохранение хеша рефреш токена
func (r *RepoPostgre) SaveRefreshToken(refreshToken *models.RefreshTokenData) error {
	op := "repo.SaveRefreshToken"

	tokenHash, err := bcrypt.GenerateFromPassword([]byte(refreshToken.RefreshToken), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s - Ошибка хеширования RefreshToken: \n%v", op, err)
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	userID, err := uuid.Parse(refreshToken.UserID)
	if err != nil {
		return fmt.Errorf("%s - Ошибка хеширования RefreshToken: \n%v", op, err) //
	}

	tokenData := &RefreshToken{
		UserID:           userID,
		RefreshTokenHash: string(tokenHash),
		ExpiresAt:        expirationTime,
		IPAdress:         refreshToken.IPAddress,
	}

	result := r.DB.Create(tokenData)
	if result.Error != nil {
		return fmt.Errorf("%s - Ошибка сохранения RefreshToken: \n%v", op, result.Error)
	}
	return nil
}

// Получение рефреш токена
func (r *RepoPostgre) GetRefreshToken(UserID string) (*models.RefreshTokenData, error) {
	op := "repo.GetRefreshToken"
	var refreshToken RefreshToken

	result := r.DB.Where("user_id = ?", UserID).First(&refreshToken)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("%s - RefreshToken не найден: \n%v", op, result.Error)
		}
		return nil, fmt.Errorf("%s - Ошибка получения RefreshToken: \n%v", op, result.Error)
	}

	response := &models.RefreshTokenData{
		UserID:       refreshToken.UserID.String(),
		RefreshToken: refreshToken.RefreshTokenHash,
		ExpiresAt:    refreshToken.ExpiresAt,
		IPAddress:    refreshToken.IPAdress,
	}
	return response, nil
}

// Удаление данных рефреш токена
func (r *RepoPostgre) DeleteRefreshToken(userID string) error {
	op := "repo.DeleteRefreshToken"

	// Попытка удалить запись, где user_id совпадает
	result := r.DB.Where("user_id = ?", userID).Delete(RefreshToken{})
	if result.Error != nil {
		return fmt.Errorf("%s - Ошибка удаления RefreshToken: \n%v", op, result.Error)
	}

	// Проверка, что хотя бы одна строка была удалена
	if result.RowsAffected == 0 {
		return fmt.Errorf("%s - Запись для удаления не найдена", op)
	}

	return nil
}
