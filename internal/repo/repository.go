package repo

import (
	"fmt"

	"github.com/mikromolekula2002/Testovoe/internal/config"
	"github.com/mikromolekula2002/Testovoe/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type RepoPostgre struct {
	DB *gorm.DB
}

type TokenRepo interface {
	SaveRefreshToken(refreshToken *models.RefreshToken) error
	GetRefreshToken(UserID string) (*models.RefreshToken, error)
	UpdateRefreshToken(refreshToken *models.RefreshToken) error
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
	err = db.AutoMigrate(&models.RefreshToken{})
	if err != nil {
		return nil, err
	}

	database := RepoPostgre{
		DB: db,
	}
	return &database, nil
}

// Сохранение хеша рефреш токена
func (r *RepoPostgre) SaveRefreshToken(refreshToken *models.RefreshToken) error {
	op := "repo.SaveRefreshToken"

	result := r.DB.Create(refreshToken)
	if result.Error != nil {
		return fmt.Errorf("%s - Ошибка сохранения RefreshToken: \n%v", op, result.Error)
	}
	return nil
}

// Получение рефреш токена
func (r *RepoPostgre) GetRefreshToken(UserID string) (*models.RefreshToken, error) {
	op := "repo.GetRefreshToken"
	var refreshToken models.RefreshToken

	result := r.DB.Where("user_id = ?", UserID).First(&refreshToken)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("%s - RefreshToken не найден: \n%v", op, result.Error)
		}
		return nil, fmt.Errorf("%s - Ошибка получения RefreshToken: \n%v", op, result.Error)
	}

	return &refreshToken, nil
}

// Обновление данных рефреш токена
func (r *RepoPostgre) UpdateRefreshToken(refreshToken *models.RefreshToken) error {
	op := "repo.UpdateRefreshToken"

	// Попытка обновить запись с существующим значением поля
	result := r.DB.Model(&models.RefreshToken{}).Where("token_hash = ?", refreshToken.TokenHash).Updates(refreshToken)
	if result.Error != nil {
		return fmt.Errorf("%s - Ошибка обновления RefreshToken: \n%v", op, result.Error)
	}

	// Проверка, что хотя бы одна строка была обновлена
	if result.RowsAffected == 0 {
		return fmt.Errorf("%s - Запись для обновления не найдена", op)
	}

	return nil
}
