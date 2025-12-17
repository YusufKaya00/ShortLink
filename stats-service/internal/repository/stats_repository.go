package repository

import (
	"time"

	"github.com/urlshortener/stats-service/internal/models"
	"gorm.io/gorm"
)

type StatsRepository struct {
	db *gorm.DB
}

func NewStatsRepository(db *gorm.DB) *StatsRepository {
	return &StatsRepository{db: db}
}

func (r *StatsRepository) RecordClick(click *models.Click) error {
	return r.db.Create(click).Error
}

func (r *StatsRepository) GetTotalClicks(shortCode string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Click{}).Where("short_code = ?", shortCode).Count(&count).Error
	return count, err
}

func (r *StatsRepository) GetClicksByDay(shortCode string, days int) ([]models.DayStats, error) {
	var stats []models.DayStats
	
	startDate := time.Now().AddDate(0, 0, -days)
	
	err := r.db.Model(&models.Click{}).
		Select("DATE(created_at) as date, COUNT(*) as clicks").
		Where("short_code = ? AND created_at >= ?", shortCode, startDate).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&stats).Error
	
	return stats, err
}

func (r *StatsRepository) GetClicksByDevice(shortCode string) ([]models.DeviceStats, error) {
	var stats []models.DeviceStats
	
	err := r.db.Model(&models.Click{}).
		Select("device, COUNT(*) as count").
		Where("short_code = ?", shortCode).
		Group("device").
		Order("count DESC").
		Scan(&stats).Error
	
	return stats, err
}

func (r *StatsRepository) GetClicksByBrowser(shortCode string) ([]models.BrowserStats, error) {
	var stats []models.BrowserStats
	
	err := r.db.Model(&models.Click{}).
		Select("browser, COUNT(*) as count").
		Where("short_code = ?", shortCode).
		Group("browser").
		Order("count DESC").
		Scan(&stats).Error
	
	return stats, err
}

func (r *StatsRepository) GetClicksByReferer(shortCode string) ([]models.RefererStats, error) {
	var stats []models.RefererStats
	
	err := r.db.Model(&models.Click{}).
		Select("referer, COUNT(*) as count").
		Where("short_code = ?", shortCode).
		Group("referer").
		Order("count DESC").
		Limit(10).
		Scan(&stats).Error
	
	return stats, err
}

func (r *StatsRepository) GetOverallStats() (*models.OverallStats, error) {
	var stats models.OverallStats
	
	// Total clicks
	r.db.Model(&models.Click{}).Count(&stats.TotalClicks)
	
	// Today's clicks
	today := time.Now().Truncate(24 * time.Hour)
	r.db.Model(&models.Click{}).Where("created_at >= ?", today).Count(&stats.TodayClicks)
	
	// Total unique URLs
	r.db.Model(&models.Click{}).Distinct("short_code").Count(&stats.TotalURLs)
	
	// Active URLs (clicked today)
	r.db.Model(&models.Click{}).Where("created_at >= ?", today).Distinct("short_code").Count(&stats.ActiveURLs)
	
	return &stats, nil
}

func (r *StatsRepository) GetRecentClicks(limit int) ([]models.Click, error) {
	var clicks []models.Click
	err := r.db.Order("created_at DESC").Limit(limit).Find(&clicks).Error
	return clicks, err
}
