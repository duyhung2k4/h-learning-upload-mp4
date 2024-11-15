package job

import (
	"app/config"

	"gorm.io/gorm"
)

type fileJob struct {
	psql *gorm.DB
}

type FileJob interface{}

func NewFileJob() FileJob {
	return &fileJob{
		psql: config.GetPsql(),
	}
}
