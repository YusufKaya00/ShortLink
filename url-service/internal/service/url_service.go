package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/urlshortener/url-service/internal/models"
	"github.com/urlshortener/url-service/internal/repository"
	"github.com/urlshortener/url-service/pkg/redis"
)

type URLService struct {
	repo  *repository.URLRepository
	redis *redis.RedisClient
}

func NewURLService(repo *repository.URLRepository, redis *redis.RedisClient) *URLService {
	return &URLService{repo: repo, redis: redis}
}

func (s *URLService) CreateURL(req *models.CreateURLRequest, userID *uuid.UUID) (*models.URLResponse, error) {
	var shortCode string

	if req.CustomCode != "" {
		// Use custom code if provided
		if s.repo.ShortCodeExists(req.CustomCode) {
			return nil, errors.New("custom code already exists")
		}
		shortCode = req.CustomCode
	} else {
		// Generate random short code
		var err error
		shortCode, err = s.generateShortCode()
		if err != nil {
			return nil, err
		}
	}

	url := &models.URL{
		ShortCode:   shortCode,
		OriginalURL: req.OriginalURL,
		UserID:      userID,
	}

	if req.ExpiresIn > 0 {
		expiresAt := time.Now().Add(time.Duration(req.ExpiresIn) * time.Hour)
		url.ExpiresAt = &expiresAt
	}

	if err := s.repo.Create(url); err != nil {
		return nil, err
	}

	return s.toURLResponse(url), nil
}

func (s *URLService) GetURL(shortCode string) (*models.URL, error) {
	url, err := s.repo.FindByShortCode(shortCode)
	if err != nil {
		return nil, err
	}

	// Check if expired
	if url.ExpiresAt != nil && url.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("URL has expired")
	}

	return url, nil
}

func (s *URLService) RecordClick(shortCode string, userAgent, ip, referer string) error {
	// Increment click count in database
	if err := s.repo.IncrementClickCount(shortCode); err != nil {
		return err
	}

	// Publish click event to Redis for Stats Service
	event := models.ClickEvent{
		ShortCode: shortCode,
		UserAgent: userAgent,
		IP:        ip,
		Referer:   referer,
		Timestamp: time.Now(),
	}

	ctx := context.Background()
	return s.redis.Publish(ctx, "url:click", event)
}

func (s *URLService) GetUserURLs(userID uuid.UUID, limit, offset int) ([]models.URLResponse, int64, error) {
	urls, total, err := s.repo.FindByUserID(userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var responses []models.URLResponse
	for _, url := range urls {
		responses = append(responses, *s.toURLResponse(&url))
	}

	return responses, total, nil
}

func (s *URLService) DeleteURL(id uuid.UUID, userID uuid.UUID) error {
	return s.repo.Delete(id, userID)
}

func (s *URLService) GetAllURLs(limit, offset int) ([]models.URLResponse, int64, error) {
	urls, total, err := s.repo.GetAll(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var responses []models.URLResponse
	for _, url := range urls {
		responses = append(responses, *s.toURLResponse(&url))
	}

	return responses, total, nil
}

func (s *URLService) generateShortCode() (string, error) {
	for i := 0; i < 10; i++ {
		bytes := make([]byte, 6)
		if _, err := rand.Read(bytes); err != nil {
			return "", err
		}

		code := base64.URLEncoding.EncodeToString(bytes)
		code = strings.TrimRight(code, "=")
		code = code[:6]

		if !s.repo.ShortCodeExists(code) {
			return code, nil
		}
	}
	return "", errors.New("failed to generate unique short code")
}

func (s *URLService) toURLResponse(url *models.URL) *models.URLResponse {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8082"
	}

	return &models.URLResponse{
		ID:          url.ID,
		ShortCode:   url.ShortCode,
		ShortURL:    baseURL + "/" + url.ShortCode,
		OriginalURL: url.OriginalURL,
		ClickCount:  url.ClickCount,
		ExpiresAt:   url.ExpiresAt,
		CreatedAt:   url.CreatedAt,
	}
}
