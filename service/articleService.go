package service

import (
	"LarsWebV0/dao"
	"LarsWebV0/middleware"
	"LarsWebV0/model"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"time"
)

// @Summary 创建文章
// @Tags 文章
// @version 1.0
// @Accept application/x-json-stream
// @Param article body model.Article true "article"
// @Param Authorization header string true "Authorization"
// @Router /article/insert [post]
func InsertArticle(context *gin.Context) {
	response := model.Response{Context: context}
	var err error
	var article model.Article
	err = context.ShouldBindJSON(&article)
	if err != nil {
		logger.Errorf("Unmarshal article fails: %v", err)
		response.Fails("Unmarshal article fails", err)
		return
	}
	userId := middleware.GetIdInToken(context)
	article.User.ID = userId
	article.CreateTime = time.Now()
	err = dao.InsertArticle(article)
	if err != nil {
		logger.Errorf("Insert Article fails: %v", err)
		response.Fails("Insert Article fails", err)
		return
	}
	response.Success(nil)
}
