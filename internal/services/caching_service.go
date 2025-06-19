package services

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"randomMeetsProject/internal/models/validators"
	"randomMeetsProject/pkg/database"
	"time"
)

type CachingServiceInterface interface {
	GetAllMeetUps(ctx context.Context) []validators.MeetUp
	CachingMeetUps(ctx context.Context, meets []validators.MeetUp)
}

type CachingService struct {
	redisClient redis.Cmdable
}

func NewCachingService() CachingServiceInterface {
	redisClient, _ := database.RedisClient()
	return &CachingService{redisClient: redisClient}
}

func (service *CachingService) GetAllMeetUps(ctx context.Context) []validators.MeetUp {
	val, err := service.redisClient.Get(ctx, "meetups").Result()
	if err != nil {
		return nil
	}

	var meetups []validators.MeetUp
	err = json.Unmarshal([]byte(val), &meetups)
	if err != nil {
		return nil
	}
	return meetups
}

func (service *CachingService) CachingMeetUps(ctx context.Context, meets []validators.MeetUp) {
	meetUps, _ := json.Marshal(meets)
	service.redisClient.Set(ctx, "meetups", meetUps, time.Minute)
}
