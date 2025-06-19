package sql_models

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email    string    `gorm:"type:varchar(255);unique;not null" validate:"required,email"`
	Username string    `gorm:"type:varchar(100);not null" validate:"required"`
	Password []byte    `gorm:"type:bytea;not null"`
	MeetUps  []MeetUp  `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	PhotoURL string    `gorm:"type:varchar(2047);"`
}
