package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"randomMeetsProject/internal/models/validators"
	"randomMeetsProject/internal/repositories"
	"time"
)

type MeetUpServiceInterface interface {
	CreateMeetUp(user *validators.MeetUp, userID uuid.UUID) (uuid.UUID, error)
	GetMeetUp(userID uuid.UUID) (*validators.MeetUp, error)
	GetAllMeetUps(ctx context.Context) ([]validators.MeetUp, error)
	DeleteMeetUp(userID uuid.UUID) error
}

type MeetUpService struct {
	repository repositories.MeetUpRepositoryInterface
}

func NewMeetUpService() (MeetUpServiceInterface, error) {
	repository, err := repositories.NewMeetUpRepository()
	if err != nil {
		return nil, err
	}
	return &MeetUpService{repository: repository}, nil
}

func (service *MeetUpService) CreateMeetUp(meet *validators.MeetUp, userID uuid.UUID) (uuid.UUID, error) {
	checkMeetExists := service.repository.CheckMeetUpExists(userID)
	if checkMeetExists {
		return uuid.Nil, errors.New("Meet Up already exists")
	}
	result, err := service.repository.CreateMeetUp(meet.Title, meet.Description, meet.Date, meet.PeopleCount, userID)
	if err != nil {
		return uuid.Nil, err
	}
	return result, nil
}

func (service *MeetUpService) GetMeetUp(userID uuid.UUID) (*validators.MeetUp, error) {
	meetEntity := service.repository.GetMeetByUserID(userID)
	if meetEntity == nil {
		return nil, errors.New("Meet Up not found")
	}
	dateString := meetEntity.Date.Format(time.RFC3339)
	result := validators.MeetUp{
		Title:       meetEntity.Title,
		Description: meetEntity.Description,
		Date:        dateString,
		PeopleCount: meetEntity.PeopleCount,
	}

	return &result, nil
}

func (service *MeetUpService) GetAllMeetUps(ctx context.Context) ([]validators.MeetUp, error) {
	cachingService := NewCachingService()
	cachingMeetUps := cachingService.GetAllMeetUps(ctx)
	if cachingMeetUps == nil {
		meetUps, err := service.repository.GetMeetUps()
		if err != nil {
			return nil, err
		}
		var result []validators.MeetUp
		for _, meetUp := range meetUps {
			meet := validators.MeetUp{
				Title:       meetUp.Title,
				Description: meetUp.Description,
				Date:        meetUp.Date.Format(time.RFC3339),
				PeopleCount: meetUp.PeopleCount,
				PhotoURL:    meetUp.User.PhotoURL,
			}
			result = append(result, meet)
		}
		cachingService.CachingMeetUps(ctx, result)
		return result, nil
	}

	return cachingMeetUps, nil
}

func (service *MeetUpService) DeleteMeetUp(userID uuid.UUID) error {
	err := service.repository.DeleteMeetUp(userID)
	if err != nil {
		return err
	}
	return nil
}
