package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Click struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	ShortCode string    `gorm:"index;not null" json:"short_code"`
	UserAgent string    `json:"user_agent"`
	IP        string    `json:"ip"`
	Referer   string    `json:"referer"`
	Country   string    `json:"country"`
	Device    string    `json:"device"`
	Browser   string    `json:"browser"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}

func (c *Click) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

type ClickEvent struct {
	ShortCode string    `json:"short_code"`
	UserAgent string    `json:"user_agent"`
	IP        string    `json:"ip"`
	Referer   string    `json:"referer"`
	Timestamp time.Time `json:"timestamp"`
}

type URLStats struct {
	ShortCode   string         `json:"short_code"`
	TotalClicks int64          `json:"total_clicks"`
	ByDay       []DayStats     `json:"by_day"`
	ByDevice    []DeviceStats  `json:"by_device"`
	ByBrowser   []BrowserStats `json:"by_browser"`
	ByReferer   []RefererStats `json:"by_referer"`
}

type DayStats struct {
	Date   string `json:"date"`
	Clicks int64  `json:"clicks"`
}

type DeviceStats struct {
	Device string `json:"device"`
	Count  int64  `json:"count"`
}

type BrowserStats struct {
	Browser string `json:"browser"`
	Count   int64  `json:"count"`
}

type RefererStats struct {
	Referer string `json:"referer"`
	Count   int64  `json:"count"`
}

type OverallStats struct {
	TotalURLs    int64 `json:"total_urls"`
	TotalClicks  int64 `json:"total_clicks"`
	TodayClicks  int64 `json:"today_clicks"`
	ActiveURLs   int64 `json:"active_urls"`
}
