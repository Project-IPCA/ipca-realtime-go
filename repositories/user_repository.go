package repositories

import (
	"gorm.io/gorm"

	"github.com/Project-IPCA/ipca-realtime-go/models"
)

type UserRepositoryQ interface{
	GetOnlineStudents(
		user *models.User,
		group_id string,
	)
	GetOnlineStudentsOld(
		userIDs *[]int,
		group_id string,
	)
	UpdateUserStatus(
		userID string,
		status bool,
	)
}

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (userRepository *UserRepository)GetOnlineStudentsOld(userIDs *[]string,group_id string){
	userRepository.DB.Table("students").
	Select("students.stu_id").
	Joins("JOIN users ON students.stu_id = users.user_id").
	Where("students.group_id = ? AND users.is_online = ?", group_id, true).
	Find(userIDs)
}

func (userRepository *UserRepository)UpdateUserStatus(userID string,status bool){
	userRepository.DB.Model(&models.User{}).Where("user_id = ?", userID).Update("is_online", status)
}