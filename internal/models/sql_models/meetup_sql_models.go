package sql_models

import (
	"github.com/google/uuid"
	"time"
)

type MeetUp struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title       string    `gorm:"size:255;not null"`
	Description string    `gorm:"type:text"`
	Date        time.Time `gorm:"type:timestamp;not null"`
	PeopleCount int       `gorm:"type:int;not null"`
	CreatedAt   time.Time `gorm:"type:timestamp;default:now()"`
	UpdatedAt   time.Time `gorm:"type:timestamp;default:now()"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index"`
	User        User      `gorm:"foreignKey:UserID"`
}
