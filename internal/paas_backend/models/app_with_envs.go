package models

// Virtual model, to have an application with its environments
type AppWithEnvs struct {
	DBApplication
	Envs []DBEnvironment `gorm:"foreignKey:ApplicationID"`
}

func (AppWithEnvs) TableName() string {
	return DBApplication{}.TableName()
}
