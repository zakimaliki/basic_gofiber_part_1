package models

import (
	"crypto/rand"
	"gofiber/src/configs"
	"math/big"
	"time"

	"gorm.io/gorm"
)

// type User struct {
// 	gorm.Model
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

type User struct {
	gorm.Model
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required,min=8,max=20"`
	Verify    string    `json:"verify"`
	UpdatedOn time.Time `json:updated_on`
}

type UserVerification struct {
	gorm.Model
	UserID    uint      `json:"user_id"`
	Token     string    `json:"token"`
	UpdatedOn time.Time `json:updated_on`
	User      User      `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

func PostUser(item *User) error {
	result := configs.DB.Create(&item)
	return result.Error
}

func PostUserVerify(item *UserVerification) error {
	result := configs.DB.Create(&item)
	return result.Error
}

func CheckUsersVerification(userID uint, token string) (*UserVerification, error) {
	var userVerification UserVerification
	if err := configs.DB.Where("user_id = ? AND token = ?", userID, token).First(&userVerification).Error; err != nil {
		return nil, err
	}
	return &userVerification, nil
}

func UpdateAccountVerification(ID uint) error {
	return configs.DB.Model(&User{}).Where("id = ?", ID).Update("verify", "true").Error
}

func UpdateUserVerify(id int) error {
	result := configs.DB.Model(&User{}).Where("id = ?", id).Update("verify", "true")
	return result.Error
}

func DeleteUsersVerification(ID uint, token string) error {
	return configs.DB.Where("id = ? AND token = ?", ID, token).Delete(&UserVerification{}).Error
}

func GenerateRandom6DigitID() (uint, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(899999))
	if err != nil {
		return 0, err
	}
	return uint(n.Int64() + 100000), nil
}

func FindEmail(input *User) []User {
	items := []User{}
	configs.DB.Raw("SELECT * FROM users WHERE email = ?", input.Email).Scan(&items)
	return items
}

func FindID(id uint) User {
	items := User{}
	configs.DB.Raw("SELECT * FROM users WHERE id = ?", id).Scan(&items)
	return items
}

// func CreateUser(user *User) (uint, error) {
// 	result := configs.DB.Create(&user)
// 	return user.ID, result.Error
// }

// func CreateUser(user )
