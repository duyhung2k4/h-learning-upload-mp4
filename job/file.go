package job

import (
	"app/config"
	"app/model"
	"fmt"
	"log"
	"os"
	"strings"

	"gorm.io/gorm"
)

type fileJob struct {
	psql *gorm.DB
}

type FileJob interface {
	DeleteFileMp4()
}

func (j *fileJob) DeleteFileMp4() {
	listFile, err := os.ReadDir("video")
	if err != nil {
		log.Println("error get list file: ", err)
		return
	}

	listUuid := []string{}
	for _, f := range listFile {
		if f.IsDir() {
			continue
		}
		uuid := strings.Split(f.Name(), ".")[0]
		listUuid = append(listUuid, uuid)
	}

	var listVideoLession []model.VideoLession
	err = j.psql.
		Model(&model.VideoLession{}).
		Where(`
			code IN ?
			AND url360p IS NOT NULL
		`, listUuid).
		Find(&listVideoLession).Error
	if err != nil {
		log.Println("error get listVideoLession: ", err)
		return
	}

	listError := []error{}
	for _, v := range listVideoLession {
		path := fmt.Sprintf("video/%s.mp4", v.Code)
		err := os.RemoveAll(path)
		if err != nil {
			listError = append(listError, err)
		}
	}

	if len(listError) > 0 {
		for _, e := range listError {
			log.Println("error delete file mp4: ", e)
		}

		return
	}
}

func NewFileJob() FileJob {
	return &fileJob{
		psql: config.GetPsql(),
	}
}
