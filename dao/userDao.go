package dao

import (
	"LarsWebV0/model"
	"errors"
)

func Register(user model.User) error {
	e := QueryByUsername(user.UserName)
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

func QueryUserById(id uint) (model.User, error) {
	var user model.User
	err := db.First(&user, id).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func QueryByUsername(username string) error {
	var user model.User
	_ = db.Where(model.User{UserName: username}).First(&user)
	if user.ID != 0 {
		return errors.New("invalid username")
	}
	return nil
}

func UpdateUserInfo(user model.User) error {
	err := db.Model(&user).Updates(user).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateUserImage(image string, id uint) error {
	err := db.Model(model.User{ID: id}).Updates(model.User{Image: image}).Error
	if err != nil {
		return err
	}
	return nil
}
