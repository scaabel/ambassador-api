package models

type User struct {
	Id           uint
	Name         string
	Email        string
	Password     []byte
	IsAmbassador bool
}
