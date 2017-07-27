package models

import (
	"github.com/go-pg/pg/orm"
)

type User struct {
	ID                int
	Email             string
	Admin             bool
	APIAccess         bool
	APIToken          string
	EncryptedPassword string
}

func ValidAPIUser(db orm.DB, token string) bool {
	var user User
	err := db.Model(&user).Where("api_access=? AND api_token=?", true, token).First()
	return err == nil && user.ID != 0
}
