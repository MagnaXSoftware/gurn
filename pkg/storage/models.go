package storage

import (
	"time"

	"gorm.io/gorm"
)

type Urn struct {
	gorm.Model

	Name        string `gorm:"unique,index"`
	Destination string
}

type Access struct {
	ID         uint `gorm:"primarykey"`
	AccessedAt time.Time
	Urn        string `gorm:"index"`
	Result     string
}
