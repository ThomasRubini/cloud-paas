package models

import "gorm.io/gorm"

type DBEnvironment struct {
	gorm.Model
	ApplicationID uint          `gorm:"uniqueIndex:idx_env_name"`
	Application   DBApplication `gorm:"foreignKey:ApplicationID"`
	Domain        string
	Name          string `gorm:"uniqueIndex:idx_env_name"`
}
