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
	userRepository.DB.Raw(`
    SELECT user.id
    FROM user
    JOIN user_student ON user.id = user_student.stu_id
    WHERE user_student.stu_group = ? AND user.status = ?
	`, group_id , "online").Scan(&userIDs)
}

func (userRepository *UserRepository)UpdateUserStatus(userID string,status bool){
	userRepository.DB.Model(&models.User{}).Where("user_id = ?", userID).Update("is_online", status)
}