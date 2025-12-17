package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/urlshortener/stats-service/internal/models"
	"github.com/urlshortener/stats-service/internal/service"
)

type ClickConsumer struct {
	redisClient *redis.Client
	service     *service.StatsService
}

func NewClickConsumer(redisClient *redis.Client, service *service.StatsService) *ClickConsumer {
	return &ClickConsumer{
		redisClient: redisClient,
		service:     service,
	}
}

func (c *ClickConsumer) Start(ctx context.Context) {
	pubsub := c.redisClient.Subscribe(ctx, "url:click")
	defer pubsub.Close()

	log.Println("Stats Consumer: Listening for click events...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Stats Consumer: Shutting down...")
			return
		default:
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				log.Printf("Stats Consumer: Error receiving message: %v", err)
				continue
			}

			var event models.ClickEvent
			if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
				log.Printf("Stats Consumer: Error unmarshaling event: %v", err)
				continue
			}

			if err := c.service.RecordClick(&event); err != nil {
				log.Printf("Stats Consumer: Error recording click: %v", err)
				continue
			}

			log.Printf("Stats Consumer: Recorded click for %s", event.ShortCode)
		}
	}
}
