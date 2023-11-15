package dao

import (
	"LarsWebV0/model"
	"errors"
	"time"
)

// 增
func AddArticle(user model.User, article model.Article) error {
	article.UserId = user.ID
	article.CreateTime = time.Now()
	err := db.Create(&article).Error
	if err != nil {
		return errors.New("add failed")
	}
	return nil
}

// 删
func DeleteArticle(user model.User, articleId string) error {
	var article model.Article
	article.UserId = user.ID
	article.ID = articleId
	err := db.Delete(&article).Error
	if err != nil {
		return errors.New("delete failed")
	}
	return nil
}

//改

// 批量查user_id
func FindAll(id string) ([]model.Article, error) {
	res := []model.Article{}
	err := db.Where("user_id = ?", id).Find(&res)
	if err.Error != nil {
		return []model.Article{}, errors.New("search failed")
	}
	return res, nil
}

// 查user_id指定article_id
/*func FindById(userId string, id model.Article) (model.Article, error) {
	var article model.Article
	err := db.Where("")
}
*/
