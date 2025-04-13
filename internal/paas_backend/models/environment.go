package models

type DBEnvironment struct {
	BaseModel
	ApplicationID uint `gorm:"uniqueIndex:idx_env_name"`
	Domain        string
	Name          string `gorm:"uniqueIndex:idx_env_name;not null;default:null"`
	Branch        string
	Running       bool
}
