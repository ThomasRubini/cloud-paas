package models

type DBEnvironment struct {
	BaseModel
	ApplicationID uint          `gorm:"uniqueIndex:idx_env_name"`
	Application   DBApplication `gorm:"foreignKey:ApplicationID"`
	Domain        string
	Name          string `gorm:"uniqueIndex:idx_env_name;not null;default:null"`
	Branch        string
	Running       bool
}
