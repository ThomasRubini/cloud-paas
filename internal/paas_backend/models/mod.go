// this package contains database models (database tables are created from these structs)
package models

import (
	"time"
)

type BaseModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
