package repository

import (
	"github.com/google/uuid"
	"github.com/urlshortener/url-service/internal/models"
	"gorm.io/gorm"
)

type URLRepository struct {
	db *gorm.DB
}

func NewURLRepository(db *gorm.DB) *URLRepository {
	return &URLRepository{db: db}
}

func (r *URLRepository) Create(url *models.URL) error {
	return r.db.Create(url).Error
}

func (r *URLRepository) FindByShortCode(code string) (*models.URL, error) {
	var url models.URL
	err := r.db.Where("short_code = ?", code).First(&url).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *URLRepository) FindByID(id uuid.UUID) (*models.URL, error) {
	var url models.URL
	err := r.db.First(&url, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *URLRepository) FindByUserID(userID uuid.UUID, limit, offset int) ([]models.URL, int64, error) {
	var urls []models.URL
	var total int64

	r.db.Model(&models.URL{}).Where("user_id = ?", userID).Count(&total)
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&urls).Error

	return urls, total, err
}

func (r *URLRepository) ShortCodeExists(code string) bool {
	var count int64
	r.db.Model(&models.URL{}).Where("short_code = ?", code).Count(&count)
	return count > 0
}

func (r *URLRepository) IncrementClickCount(code string) error {
	return r.db.Model(&models.URL{}).
		Where("short_code = ?", code).
		UpdateColumn("click_count", gorm.Expr("click_count + ?", 1)).Error
}

func (r *URLRepository) Delete(id uuid.UUID, userID uuid.UUID) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.URL{}).Error
}

func (r *URLRepository) GetAll(limit, offset int) ([]models.URL, int64, error) {
	var urls []models.URL
	var total int64

	r.db.Model(&models.URL{}).Count(&total)
	err := r.db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&urls).Error

	return urls, total, err
}
