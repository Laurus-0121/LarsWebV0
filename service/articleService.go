package service

import (
	"LarsWebV0/dao"
	"LarsWebV0/middleware"
	"LarsWebV0/model"
	"fmt"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"strconv"
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
	article.CreateTime = time.Now()
	err = context.ShouldBindJSON(&article)
	fmt.Println("error json:", article)
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

// @Summary 删除文章
// @Tags 文章
// @version 1.0
// @Accept application/x-json-stream
// @Param articleId path string true "articleId"
// @Param Authorization header string true "Authorization"
// @Router /article/deletedById/{articleId} [post]
func DeleteArticle(context *gin.Context) {
	response := model.Response{Context: context}
	userId := middleware.GetIdInToken(context)
	articleId := context.Param("articleId")
	article, err := dao.GetArticleById(articleId)
	if err != nil {
		logger.Errorf("find article fails when delete: %v", err)
		response.Fails("find article fails when delete", err)
		return
	}
	if article.User.ID != userId {
		logger.Errorf("no auth to delete article")
		response.Fails("no auth to delete article", nil)
		return
	}
	err = dao.DeletedArticleById(articleId)
	if err != nil {
		logger.Errorf("delete article fails: %v", err)
		response.Fails("delete article fails", err)
		return
	}
	response.Success(nil)
}

// @Summary 首页获取文章列表
// @Tags 文章
// @version 1.0
// @Accept application/x-json-stream
// @Param curPage path string true "curPage"
// @Param pageSize path string true "pageSize"
// @Router /article/getArticleList/{curPage}/{pageSize} [get]
func GetArticleList(context *gin.Context) {
	resposne := model.Response{Context: context}
	var err error
	var curPage, pageSize int64
	curPage, err = strconv.ParseInt(context.Param("curPage"), 10, 64)
	pageSize, err = strconv.ParseInt(context.Param("pageSize"), 10, 64)
	articleList, err := dao.GetArticleList(int(curPage), int(pageSize))
	if err != nil {
		logger.Errorf("get article list fals: %v", err)
		resposne.Fails("get article list fails", err)
		return
	}
	resposne.Success(articleList)
}

// @Summary 获取文章轮播图
// @Tags 文章
// @version 1.0
// @Accept application/x-json-stream
// @Router /article/getArticleSwiper [get]
func GetArticleSwiper(context *gin.Context) {
	response := model.Response{Context: context}
	articleSwiperList, err := dao.GetArticleSwiper()
	if err != nil {
		response.Fails("get article swiper fails: %v", err)
		return
	}
	response.Success(articleSwiperList)
}

// @Summary 根据id获取文章
// @Tags 文章
// @version 1.0
// @Accept application/x-json-stream
// @Param id path string true "id"
// @Param Authorization header string true "Authorization"
// @Router /article/getById/{id} [get]
func GetArticleById(context *gin.Context) {
	response := model.Response{Context: context}
	id := context.Param("id")
	article, err := dao.GetArticleById(id)
	if err != nil {
		logger.Errorf("get Article by id fails: %v", err)
		response.Fails("get Article by id fails", err)
		return
	}
	article.View += 1
	if err = dao.UpdateArticleById(article); err != nil {
		logger.Errorf("update Article view fails: %v", err)
	}
	var getArticleByIdDTO GetArticleByIdDTO
	if user, err := dao.QueryUserById(article.User.ID); err != nil {
		logger.Errorf("fix article user fails: %v", err)
	} else {
		article.User.UserName = user.UserName
		article.User.Image = user.Image
		//点赞
		likeCount, _ := dao.GetArticleLikeCount(article.ID)
		article.Like = likeCount
		//收藏
		collectCount, _ := dao.GetArticleCollectCount(article.ID)
		article.Collect = collectCount
		getArticleByIdDTO.Article = article
		curUserId := middleware.GetIdInToken(context)
		isLikeArticle, _ := dao.GetIsLikeArticle(curUserId, article.ID)
		isCollectArticle, _ := dao.GetIsCollectArticle(curUserId, article.ID)
		getArticleByIdDTO.IsLike = isLikeArticle
		getArticleByIdDTO.IsCollect = isCollectArticle
		//评论
		commentDtos, _ := dao.GetArticleCommentByArticleId(article.ID)
		for i := 0; i < len(commentDtos); i++ {
			commentUser, _ := dao.QueryUserById(commentDtos[i].UserId)
			commentDtos[i].Username = commentUser.UserName
			commentDtos[i].UserImage = commentUser.Image
		}
		getArticleByIdDTO.Comment = commentDtos
	}
	response.Success(getArticleByIdDTO)
}
