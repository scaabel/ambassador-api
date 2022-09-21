package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	Id           uint
	Name         string
	Email        string
	Password     []byte
	IsAmbassador bool
}

func (user *User) SetPassword(password string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)

	user.Password = hashedPassword
}

func (user *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}
