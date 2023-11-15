package dao

import (
	"LarsWebV0/model"
	"fmt"
	"testing"
)

func TestRegisterUser(T *testing.T) {
	SetupDB()
	err := RegisterUser(model.User{UserName: "Lars", PassWord: "0421"})
	err1 := RegisterUser(model.User{UserName: "Laurus", PassWord: "0121"})
	err2 := RegisterUser(model.User{UserName: "kiwi", PassWord: "0121"})
	if err != nil {
		fmt.Sprintf("err")
	}
	if err1 != nil {
		fmt.Sprintf("err1")
	}
	if err2 != nil {
		fmt.Sprintf("failed")
	} else {
		fmt.Sprintf("success")
	}
}

func TestUserLogin(t *testing.T) {
	SetupDB()
	user, err := UserLogin("yyy", "0421")
	user1, err1 := UserLogin("Lars", "0421")
	if err != nil {
		fmt.Sprintf(err.Error())
	} else {
		fmt.Printf(user.UserName)
	}
	if err != nil {
		fmt.Sprintf(err1.Error())
	} else {
		fmt.Printf(user1.UserName)
	}
}
