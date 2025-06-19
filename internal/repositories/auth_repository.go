package repositories

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"randomMeetsProject/internal/models/sql_models"
	"randomMeetsProject/internal/utils"
	"randomMeetsProject/pkg/database"
)

type AuthRepositoryInterface interface {
	CreateUser(username, email, password string) (uuid.UUID, error)
	CheckUserExists(username, email string) bool
	CheckUser(username, email, password string) (sql_models.User, error)
	GetUserByUUID(uuid uuid.UUID) (*sql_models.User, *sql_models.MeetUp, error)
	AddPhotoUrl(userUUID uuid.UUID, url string) error
}

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository() (AuthRepositoryInterface, error) {
	db, err := database.DB()
	if err != nil {
		return nil, err
	}
	return &AuthRepository{db: db}, nil
}

func (repo *AuthRepository) CreateUser(username, email, password string) (uuid.UUID, error) {
	hashedPassword, err := utils.HashedPassword(password)
	if err != nil {
		return uuid.Nil, err
	}

	newUser := &sql_models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}
	err = repo.db.Create(newUser).Error
	if err != nil {
		return uuid.Nil, err
	}

	return newUser.ID, nil
}

func (repo *AuthRepository) CheckUserExists(username, email string) bool {
	var user sql_models.User
	repo.db.Where("username = ? AND email = ?", username, email).First(&user)
	return user.ID != uuid.Nil
}

func (repo *AuthRepository) CheckUser(username, email, password string) (sql_models.User, error) {
	var user sql_models.User
	err := repo.db.Where("username = ? AND email = ?", username, email).First(&user).Error
	if err != nil {
		return sql_models.User{}, err
	}

	err = utils.VerifyPassword(string(user.Password), password)
	if err != nil {
		return sql_models.User{}, errors.New("invalid credentials")
	}

	return user, nil
}

func (repo *AuthRepository) GetUserByUUID(uuid uuid.UUID) (*sql_models.User, *sql_models.MeetUp, error) {
	var user sql_models.User
	var meetUp sql_models.MeetUp
	repo.db.Where("id = ?", uuid).First(&user)
	err := repo.db.Where("user_id = ?", uuid).First(&meetUp).Error
	if err != nil {
		return &user, nil, err
	}
	return &user, &meetUp, nil
}

func (repo *AuthRepository) AddPhotoUrl(userUUID uuid.UUID, url string) error {
	var user sql_models.User
	repo.db.Where("id = ?", userUUID).First(&user)
	user.PhotoURL = url
	err := repo.db.Save(&user).Error
	if err != nil {
		return err
	}
	return nil
}
