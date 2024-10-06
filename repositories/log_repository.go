package repositories

import (
	"gorm.io/gorm"

	"github.com/Project-IPCA/ipca-realtime-go/models"
)

type ClassLogRepositoryQ interface{
	GetActivityLog(
		activityLog *[]models.ActivityLog,
		group_id string,
	)
	GetActivityLogOld(
		activityLog *[]models.ActivityLog,
		group_id string,
	)
}

type ClassLogRepository struct {
	DB *gorm.DB
}

func NewClassLogRepository(db *gorm.DB) *ClassLogRepository {
	return &ClassLogRepository{DB: db}
}

func (classLogRepository *ClassLogRepository)GetActivityLogOld(activityLog *[]models.ActivityLog,group_id string){
	classLogRepository.DB.Where("group_id",group_id).Find(&activityLog)
}