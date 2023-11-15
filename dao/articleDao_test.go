package dao

import (
	"LarsWebV0/model"
	"testing"
)

func TestFindAll(t *testing.T) {
	SetupDB()
	_, err := FindAll("1")
	if err != nil {
		return
	}

}

func TestAddArcicle(t *testing.T) {
	SetupDB()
	err := AddArticle(model.User{ID: 1}, model.Article{})
	if err != nil {
		return
	}
}
func TestDeleteArticle(t *testing.T) {
	SetupDB()
	err := DeleteArticle(model.User{ID: 2}, "54")
	if err != nil {
		return
	}
}
