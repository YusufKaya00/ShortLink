package service

import (
	"strings"

	"github.com/urlshortener/stats-service/internal/models"
	"github.com/urlshortener/stats-service/internal/repository"
)

type StatsService struct {
	repo *repository.StatsRepository
}

func NewStatsService(repo *repository.StatsRepository) *StatsService {
	return &StatsService{repo: repo}
}

func (s *StatsService) RecordClick(event *models.ClickEvent) error {
	click := &models.Click{
		ShortCode: event.ShortCode,
		UserAgent: event.UserAgent,
		IP:        event.IP,
		Referer:   event.Referer,
		Device:    s.parseDevice(event.UserAgent),
		Browser:   s.parseBrowser(event.UserAgent),
		CreatedAt: event.Timestamp,
	}

	return s.repo.RecordClick(click)
}

func (s *StatsService) GetURLStats(shortCode string) (*models.URLStats, error) {
	totalClicks, err := s.repo.GetTotalClicks(shortCode)
	if err != nil {
		return nil, err
	}

	byDay, _ := s.repo.GetClicksByDay(shortCode, 30)
	byDevice, _ := s.repo.GetClicksByDevice(shortCode)
	byBrowser, _ := s.repo.GetClicksByBrowser(shortCode)
	byReferer, _ := s.repo.GetClicksByReferer(shortCode)

	return &models.URLStats{
		ShortCode:   shortCode,
		TotalClicks: totalClicks,
		ByDay:       byDay,
		ByDevice:    byDevice,
		ByBrowser:   byBrowser,
		ByReferer:   byReferer,
	}, nil
}

func (s *StatsService) GetOverallStats() (*models.OverallStats, error) {
	return s.repo.GetOverallStats()
}

func (s *StatsService) GetRecentClicks(limit int) ([]models.Click, error) {
	return s.repo.GetRecentClicks(limit)
}

func (s *StatsService) parseDevice(userAgent string) string {
	ua := strings.ToLower(userAgent)
	
	if strings.Contains(ua, "mobile") || strings.Contains(ua, "android") {
		return "Mobile"
	}
	if strings.Contains(ua, "tablet") || strings.Contains(ua, "ipad") {
		return "Tablet"
	}
	return "Desktop"
}

func (s *StatsService) parseBrowser(userAgent string) string {
	ua := strings.ToLower(userAgent)
	
	switch {
	case strings.Contains(ua, "chrome") && !strings.Contains(ua, "edge"):
		return "Chrome"
	case strings.Contains(ua, "firefox"):
		return "Firefox"
	case strings.Contains(ua, "safari") && !strings.Contains(ua, "chrome"):
		return "Safari"
	case strings.Contains(ua, "edge"):
		return "Edge"
	case strings.Contains(ua, "opera"):
		return "Opera"
	default:
		return "Other"
	}
}
