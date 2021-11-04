package papers

import (
	"time"

	"gorm.io/gorm"
)

type Paper struct {
	gorm.Model
	Id            string `gorm:"primaryKey"`
	Name          string
	DatePublished time.Time
	DateUpdated   time.Time
	DateRead      time.Time
	Description   string
	Url           string
	Type          string
}
