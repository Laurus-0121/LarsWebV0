package dao

import (
	"LarsWebV0/config"
	"LarsWebV0/model"
	"context"
	"encoding/json"
	"errors"
	"github.com/olivere/elastic/v7"
	logger "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const articleLikeCollectionName = "articleLikeMap"
const articleCollectCollectionName = "articleCollectMap"
const articleCommentCollectionName = "articleComment"

// 增
func AddArticle(user model.User, article model.Article) error {
	/*article.UserId = user.ID
	article.CreateTime = time.Now()
	err := db.Create(&article).Error
	if err != nil {
		return errors.New("add failed")
	}
	return nil*/

	//_, err := EsClient.Index().Index(config.EsIndex).BodyJson(article).Do(context.TODO())
	//if err != nil {
	//	logger.Errorf("insert artice fails: %v", err)
	//	return err
	//}
	return nil
}

func InsertArticle(article model.Article) error {
	_, err := EsClient.Index().Index(config.EsIndex).BodyJson(article).Do(context.TODO())
	if err != nil {
		logger.Errorf("insert artice fails: %v", err)
		return err
	}
	return nil
}

func UpdateArticleById(article model.Article) error {
	_, err := EsClient.Update().Index(config.EsIndex).Id(article.ID).Doc(article).Do(context.TODO())
	if err != nil {
		logger.Errorf("update article fails: %v", err)
		return err
	}
	return nil
}

/*
// 删

	func DeleteArticle(user model.User, articleId string) error {
		var article model.Article
		article.User.ID = user.ID
		article.ID = articleId
		err := db.Delete(&article).Error
		if err != nil {
			return errors.New("delete failed")
		}
		return nil
	}
*/
func DeletedArticleById(articleId string) error {
	_, err := EsClient.Delete().Index(config.EsIndex).Id(articleId).Do(context.TODO())
	if err != nil {
		return err
	}
	return nil
}

func GetArticleById(id string) (model.Article, error) {
	var article model.Article
	result, err := EsClient.Get().Index(config.EsIndex).Id(id).Do(context.TODO())
	if err != nil {
		logger.Errorf("get article by id fails: %v", err)
		return article, err
	}
	json.Unmarshal(result.Source, &article)
	article.ID = id
	return article, nil
}

// 批量查user_id
func GetArticleList(curPage, pageSize int) (ArticleListResponse, error) {
	var articleResponse ArticleListResponse
	sorts := []elastic.Sorter{
		elastic.NewFieldSort("_score").Desc(),
		elastic.NewFieldSort("view").Desc(),
	}
	res, err := EsClient.Search().Index(config.EsIndex).
		From(pageSize * (curPage - 1)).
		Size(pageSize).
		SortBy(sorts...).Do(context.TODO())
	if err != nil {
		return articleResponse, err
	}
	articles := make([]model.Article, 0)
	userIdUserMap := make(map[uint]model.User, 0)
	//遍历查询结果中的每一个命中文档
	for _, hit := range res.Hits.Hits {
		var article model.Article
		json.Unmarshal(hit.Source, &article)
		article.ID = hit.Id
		if _, ok := userIdUserMap[article.User.ID]; !ok {
			user, err := QueryUserById(article.User.ID)
			if err != nil {
				logger.Errorf("fix article user fails: %v", err)
			} else {
				userIdUserMap[article.User.ID] = user
			}
		}
		articles = append(articles, article)
	}
	for i := 0; i < len(articles); i++ {
		if _, ok := userIdUserMap[articles[i].User.ID]; ok {
			tmpUser := userIdUserMap[articles[i].User.ID]
			articles[i].User.UserName = tmpUser.UserName
			articles[i].User.Image = tmpUser.Image
		}
	}
	articleResponse.Articles = articles
	articleResponse.Total = res.TotalHits()
	return articleResponse, nil
}
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

func GetIsLikeArticle(userId uint, articleId string) (bool, error) {
	var articleLikeMaps []*ArticleLikeMap
	filter := bson.M{
		"user_id":    userId,
		"article_id": articleId,
	}
	result, err := database.Collection(articleLikeCollectionName).Find(context.TODO(), filter)
	if err != nil {
		return false, err
	}
	result.All(context.TODO(), &articleLikeMaps)
	return len(articleLikeMaps) > 0, nil
}

func LikeArticle(userId uint, articleId string) error {
	isLikeArticle, err := GetIsLikeArticle(userId, articleId)
	if err != nil {
		return err
	}
	//先查询是否点赞过，如果没有点赞，如果有取消点赞
	if !isLikeArticle {
		articleLikeMap := ArticleLikeMap{
			UserId:     userId,
			ArticleId:  articleId,
			CreateTime: time.Now(),
		}
		_, err = database.Collection(articleLikeCollectionName).InsertOne(context.TODO(), articleLikeMap)
		if err != nil {
			return err
		}
	} else {
		filter := bson.M{
			"user_id":    userId,
			"article_id": articleId,
		}
		_, err = database.Collection(articleLikeCollectionName).DeleteMany(context.TODO(), filter)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetMyLikeArticle(userId uint, curPage, pageSize int64) ([]model.Article, error) {
	var articles []model.Article
	var articleLikeMap []ArticleLikeMap
	filter := bson.M{
		"user_id": userId,
	}
	findOptions := options.Find().SetSkip((curPage - 1) * pageSize).SetLimit(pageSize)
	result, err := database.Collection(articleLikeCollectionName).Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	result.All(context.TODO(), &articleLikeMap)
	for _, articleLike := range articleLikeMap {
		article, _ := GetArticleById(articleLike.ArticleId)
		articles = append(articles, article)
	}
	return articles, nil
}

func CollectArticle(userId uint, articleId string) error {
	isCollectArticle, err := GetIsCollectArticle(userId, articleId)
	if err != nil {
		return err
	}
	if !isCollectArticle {
		articleCollectMap := ArticleCollectMap{
			UserId:     userId,
			ArticleId:  articleId,
			CreateTime: time.Now(),
		}
		_, err = database.Collection(articleCollectCollectionName).InsertOne(context.TODO(), articleCollectMap)
		if err != nil {
			return err
		}
	} else {
		filter := bson.M{
			"user_id":    userId,
			"article_id": articleId,
		}
		_, err = database.Collection(articleCollectCollectionName).DeleteMany(context.TODO(), filter)
		if err != nil {
			return err
		}
	}
	return nil
}
func GetIsCollectArticle(userId uint, articleId string) (bool, error) {
	var articleCollectMaps []*ArticleCollectMap
	filter := bson.M{
		"user_id":    userId,
		"article_id": articleId,
	}
	result, err := database.Collection(articleCollectCollectionName).Find(context.TODO(), filter)
	if err != nil {
		return false, err
	}
	result.All(context.TODO(), &articleCollectMaps)
	return len(articleCollectMaps) > 0, nil
}

func GetMyCollectArticle(userId uint, curPage, pageSize int64) ([]model.Article, error) {
	var articles []model.Article
	var articleCollectMap []ArticleCollectMap
	filter := bson.M{
		"user_id": userId,
	}
	findOptions := options.Find().SetSkip((curPage - 1) * pageSize).SetLimit(pageSize)
	result, err := database.Collection(articleCollectCollectionName).Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	result.All(context.TODO(), &articleCollectMap)
	for _, articleCollect := range articleCollectMap {
		article, _ := GetArticleById(articleCollect.ArticleId)
		articles = append(articles, article)
	}
	return articles, nil
}

func GetArticleCollectCount(articleId string) (int64, error) {
	filter := bson.M{
		"article_id": articleId,
	}
	count, err := database.Collection(articleCollectCollectionName).CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func CommentArticle(articleComment ArticleComment) error {
	_, err := database.Collection(articleCommentCollectionName).InsertOne(context.TODO(), articleComment)
	if err != nil {
		return err
	}
	return nil
}

func GetArticleCommentByArticleId(articleId string) ([]ArticleCommentDto, error) {
	var articleComments []ArticleCommentDto
	filter := bson.M{
		"article_id": articleId,
	}
	result, err := database.Collection(articleCommentCollectionName).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	result.All(context.TODO(), &articleComments)
	return articleComments, nil
}

func DeleteCommentById(objectID primitive.ObjectID) error {
	filter := bson.M{
		"_id": objectID,
	}
	_, err := database.Collection(articleCommentCollectionName).DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}

func GetCommentById(objectID primitive.ObjectID) (ArticleComment, error) {
	var articleComment ArticleComment
	filter := bson.M{
		"_id": objectID,
	}
	result := database.Collection(articleCommentCollectionName).FindOne(context.TODO(), filter)
	err := result.Decode(&articleComment)
	if err != nil {
		return articleComment, err
	}
	return articleComment, nil
}

func GetMyArticle(userId uint) ([]model.Article, error) {
	var articles []model.Article
	query := elastic.NewBoolQuery().Must(elastic.NewMatchQuery("user.ID", userId))
	result, err := EsClient.Search().Index(config.EsIndex).Query(query).Do(context.TODO())
	if err != nil {
		return nil, err
	}
	for _, hit := range result.Hits.Hits {
		var article model.Article
		json.Unmarshal(hit.Source, &article)
		articles = append(articles, article)
	}
	return articles, nil
}

func GetArticleSwiper() ([]string, error) {
	articleSwiperList := make([]string, 0)
	boolQuery := elastic.NewBoolQuery()
	sorts := []elastic.Sorter{
		elastic.NewFieldSort("_score").Desc(),
		elastic.NewFieldSort("view").Desc(),
	}
	result, err := EsClient.Search().
		Index(config.EsIndex).
		Query(boolQuery).
		SortBy(sorts...).
		From(0).
		Size(5).
		Do(context.TODO())
	if err != nil {
		logger.Errorf("get article swiper fails: %v", err)
		return articleSwiperList, err
	}
	for _, hit := range result.Hits.Hits {
		var article model.Article
		json.Unmarshal(hit.Source, &article)
		articleSwiperList = append(articleSwiperList, article.Image)
	}
	return articleSwiperList, nil
}

func GlobalSearchArticleOrderByView(keyword string, curPage, pageSize int) (ArticleListResponse, error) {
	var articleResponse ArticleListResponse
	boolQuery := elastic.NewBoolQuery().Should(elastic.NewMatchQuery("body", keyword))
	sorts := []elastic.Sorter{
		elastic.NewFieldSort("_score").Desc(),
		elastic.NewFieldSort("view").Desc(),
	}
	bodyHeightLight := elastic.NewHighlight().Fields(elastic.NewHighlighterField("body"))
	result, err := EsClient.Search().
		Index(config.EsIndex).
		Query(boolQuery).
		Highlight(bodyHeightLight).
		From(pageSize * (curPage - 1)).
		Size(pageSize).
		SortBy(sorts...).
		Do(context.TODO())
	if err != nil {
		logger.Errorf("global search article fails: %v", err)
		return articleResponse, err
	}
	articles := make([]model.Article, 0)
	userIdUserMap := make(map[uint]model.User, 0)
	for _, hit := range result.Hits.Hits {
		var article model.Article
		json.Unmarshal(hit.Source, &article)
		article.ID = hit.Id
		if hit.Highlight["body"] != nil {
			article.Content = hit.Highlight["body"][0]
		}
		if _, ok := userIdUserMap[article.User.ID]; !ok {
			user, err := QueryUserById(article.User.ID)
			if err != nil {
				logger.Errorf("fix article user fails: %v", err)
			} else {
				userIdUserMap[article.User.ID] = user
			}
		}
		articles = append(articles, article)
	}
	for i := 0; i < len(articles); i++ {
		if _, ok := userIdUserMap[articles[i].User.ID]; ok {
			tmpUser := userIdUserMap[articles[i].User.ID]
			articles[i].User.UserName = tmpUser.UserName
			articles[i].User.Image = tmpUser.Image
		}
	}
	articleResponse.Articles = articles
	articleResponse.Total = result.TotalHits()
	return articleResponse, nil
}

type ArticleListResponse struct {
	Articles []model.Article
	Total    int64
}

type ArticleLikeMap struct {
	UserId     uint      `json:"user_id" bson:"user_id"`
	ArticleId  string    `json:"article_id" bson:"article_id"`
	CreateTime time.Time `json:"create_time" bson:"create_time"`
}

type ArticleCollectMap struct {
	UserId     uint      `json:"user_id" bson:"user_id"`
	ArticleId  string    `json:"article_id" bson:"article_id"`
	CreateTime time.Time `json:"create_time" bson:"create_time"`
}

type ArticleComment struct {
	UserId    uint   `json:"user_id" bson:"user_id"`
	ArticleId string `json:"article_id" bson:"article_id"`
	Text      string `json:"text" bson:"text"`
}

type ArticleCommentDto struct {
	Id        string `json:"_id" bson:"_id"`
	UserId    uint   `json:"user_id" bson:"user_id"`
	ArticleId string `json:"article_id" bson:"article_id"`
	Text      string `json:"text" bson:"text"`
	Username  string `json:"username"`
	UserImage string `json:"user_image"`
}
