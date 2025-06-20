package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"randomMeetsProject/internal/models/validators"
	"randomMeetsProject/internal/repositories"
	"randomMeetsProject/internal/utils"
)

type AuthServiceInterface interface {
	RegisterUser(user *validators.LoginUser) (uuid.UUID, error)
	LoginUser(user *validators.LoginUser) (validators.JWTTokens, error)
	GetMe(userUUID uuid.UUID) (validators.User, error)
	AddPhotoUrl(userUUID uuid.UUID, url string) error
	ConfirmEmail(userUUID uuid.UUID) error
}

type AuthService struct {
	repository repositories.AuthRepositoryInterface
}

func NewAuthService() (AuthServiceInterface, error) {
	repository, err := repositories.NewAuthRepository()
	if err != nil {
		return nil, err
	}

	return AuthService{repository: repository}, nil
}

func (service AuthService) RegisterUser(user *validators.LoginUser) (uuid.UUID, error) {
	checkUserExists := service.repository.CheckUserExists(user.Username, user.Email)
	if checkUserExists {
		return uuid.Nil, errors.New("User already exists")
	}

	result, err := service.repository.CreateUser(user.Username, user.Email, user.Password)
	newPublisher := NewPublisher("email-confirm", result.String()+" "+user.Email)
	if newPublisher != nil {
		return uuid.Nil, err
	}
	if err != nil {
		return uuid.Nil, err
	}
	return result, nil
}

func (service AuthService) LoginUser(user *validators.LoginUser) (validators.JWTTokens, error) {
	checkUser, err := service.repository.CheckUser(user.Username, user.Email, user.Password)
	if err != nil {
		return validators.JWTTokens{}, err
	}
	jwt, err := utils.NewTokens(checkUser)
	if err != nil {
		return validators.JWTTokens{}, err
	}
	return jwt, nil
}

func (service AuthService) GetMe(userUUID uuid.UUID) (validators.User, error) {
	user, meet, err := service.repository.GetUserByUUID(userUUID)
	if err != nil && user != nil {
		return validators.User{
			Username: user.Username,
			Email:    user.Email,
			PhotoURL: user.PhotoURL,
			Meet:     validators.MeetUp{},
		}, err
	}
	meetUp := validators.MeetUp{
		Title:       meet.Title,
		Description: meet.Description,
		Date:        meet.Date.Format(time.RFC3339),
		PeopleCount: meet.PeopleCount,
	}
	result := validators.User{
		Username: user.Username,
		Email:    user.Email,
		PhotoURL: user.PhotoURL,
		Meet:     meetUp,
	}
	return result, nil
}

func (service AuthService) AddPhotoUrl(userUUID uuid.UUID, url string) error {
	err := service.repository.AddPhotoUrl(userUUID, url)
	if err != nil {
		return err
	}
	return nil
}

func (service AuthService) ConfirmEmail(userUUID uuid.UUID) error {
	err := service.repository.CreateUserActive(userUUID)
	if err != nil {
		return err
	}
	return nil
}
