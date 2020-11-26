package user

import (
	"goblog/app/models"
)

// User ?????
type User struct {
	models.BaseModel

	Name            string `gorm:"type:varchar(255);not null;unique" valid:"name"`
	Email           string `gorm:"type:varchar(255);unique;" valid:"email"`
	Password        string `gorm:"type:varchar(255)" valid:"password"`
	PasswordComfirm string `gorm:"-" valid:"password_comfirm"`
}

// ComparePassword ????????
func (u User) ComparePassword(password string) bool {
	return u.Password == password
}
