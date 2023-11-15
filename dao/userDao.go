package dao

import (
	"LarsWebV0/model"
	"errors"
)

func RegisterUser(user model.User) error {
	e := FindByUsername(user.UserName)
	if e != nil {
		return errors.New("name already exists")
	}
	err := db.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}
func UserLogin(userName, userPassword string) (model.User, error) {
	var user model.User
	err := db.Where(model.User{UserName: userName, PassWord: userPassword}).First(&user).Error
	if err != nil {
		return model.User{}, errors.New("login failed")
	}
	return user, nil
}

func FindByUsername(username string) error {
	var user model.User
	_ = db.Where(model.User{UserName: username}).First(&user)
	if user.ID != 0 {
		return errors.New("invalid username")
	}
	return nil
}
