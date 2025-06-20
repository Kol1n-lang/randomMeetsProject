package repositories

import (
	"randomMeetsProject/internal/models/validators"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"randomMeetsProject/internal/models/sql_models"
	"randomMeetsProject/pkg/database"
)

type MeetUpRepositoryInterface interface {
	CreateMeetUp(title, description, date string, peopleCount int, userID uuid.UUID) (uuid.UUID, error)
	CheckMeetUpExists(userID uuid.UUID) bool
	GetMeetByUserID(userID uuid.UUID) *sql_models.MeetUp
	GetMeetUps() ([]sql_models.MeetUp, error)
	DeleteMeetUp(userID uuid.UUID) error
	UpdateMeetUp(meetUp *validators.PutMeetUp, userUUID uuid.UUID) error
}

type MeetUpRepository struct {
	db *gorm.DB
}

func NewMeetUpRepository() (MeetUpRepositoryInterface, error) {
	db, err := database.DB()
	if err != nil {
		return nil, err
	}
	return &MeetUpRepository{db: db}, nil
}

func (repo *MeetUpRepository) CreateMeetUp(title, description, date string, peopleCount int, userID uuid.UUID) (uuid.UUID, error) {
	parsedDate, err := time.Parse("2006-01-02-15-04", date)
	if err != nil {
		return uuid.Nil, err
	}
	meetUp := sql_models.MeetUp{
		Title:       title,
		Description: description,
		Date:        parsedDate,
		PeopleCount: peopleCount,
		UserID:      userID,
	}

	err = repo.db.Create(&meetUp).Error
	if err != nil {
		return uuid.Nil, err
	}

	return meetUp.ID, nil
}

func (repo *MeetUpRepository) CheckMeetUpExists(userID uuid.UUID) bool {
	var meetUp sql_models.MeetUp
	repo.db.Where("user_id = ?", userID).First(&meetUp)
	return meetUp.ID != uuid.Nil
}

func (repo *MeetUpRepository) GetMeetByUserID(userID uuid.UUID) *sql_models.MeetUp {
	var meetUp sql_models.MeetUp
	repo.db.
		Where("user_id = ?", userID).
		Preload("User").
		First(&meetUp)
	return &meetUp
}

func (repo *MeetUpRepository) GetMeetUps() ([]sql_models.MeetUp, error) {
	var meetUps []sql_models.MeetUp
	err := repo.db.
		Preload("User").
		Where("date > ?", time.Now()).
		Find(&meetUps).
		Error
	return meetUps, err
}

func (repo *MeetUpRepository) DeleteMeetUp(userID uuid.UUID) error {
	meetUp := sql_models.MeetUp{
		UserID: userID,
	}
	err := repo.db.Where(&meetUp).Delete(&meetUp).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *MeetUpRepository) UpdateMeetUp(meetUp *validators.PutMeetUp, userUUID uuid.UUID) error {
	var meet sql_models.MeetUp
	repo.db.Where("user_id = ?", userUUID).First(&meet)
	if meetUp.Title != "" {
		meet.Title = meetUp.Title
	}
	if meetUp.Description != "" {
		meet.Description = meetUp.Description
	}
	if meetUp.Date != "" {
		parsedDate, _ := time.Parse("2006-01-02-15-04", meetUp.Date)
		meet.Date = parsedDate
	}
	if meetUp.PeopleCount > 0 {
		meet.PeopleCount = meetUp.PeopleCount
	}

	return repo.db.Save(&meet).Error
}
