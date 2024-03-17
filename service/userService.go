package service

import (
	"LarsWebV0/config"
	"LarsWebV0/dao"
	"LarsWebV0/middleware"
	"LarsWebV0/model"
	"LarsWebV0/utils"
	"crypto/md5"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	logger "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

// @Summary 校验token合法性
// @Tags 用户
// @version 1.0
// @Accept application/x-json-stream
// @Param Authorization header string true "Authorization"
// @Router /user/checkToken [post]
func CheckToken(context *gin.Context) {
	response := model.Response{Context: context}
	response.Success("valid token")
}

// @Summary 用户注册
// @Tags 用户
// @version 1.0
// @Accept application/x-json-stream
// @Param user body model.User true "user"
// @Router /user/register [post]
func Register(context *gin.Context) {
	response := model.Response{Context: context}
	var user model.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		logger.Errorf("Unmarshal user fails: %v", err)
		response.Fails("Unmarshal user fails", err)
		return
	}
	if user.UserName == "" || user.PassWord == "" {
		logger.Errorf("illegal user")
		response.Fails("illegal user", err)
		return
	}
	encodePassWord := fmt.Sprintf("%x", md5.Sum([]byte(user.PassWord)))
	user.PassWord = encodePassWord
	user.Image = "https://larslarslar-laurus.oss-cn-beijing.aliyuncs.com/LarsWebv0/image/users/7cf1d100f77de32c4812ff254175164.jpg?Expires=1710334473&OSSAccessKeyId=TMP.3KdjBWq499CVRFCvoSw4nonNDJMxTxy7Dz8cgjuVdzqCATL5uNbpnbGKETJaJG5u9CUX7nsLmYofUExGhPWTDRsGpjaaBq&Signature=sdnXzOBXjA7Wiuv8JdqztnuCRg8%3D"
	err = dao.Register(user)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			logger.Errorf("账号已存在，请直接登录")
			response.Fails("账号已存在，请直接登录", nil)
		} else {
			logger.Errorf("Register user fails: %v", err)
			response.Fails("Register user fails", err)
		}
		return
	}
	response.Success(nil)
}

// @Summary 用户登录
// @Tags 用户
// @version 1.0
// @Accept application/x-json-stream
// @Param user body model.User true "user"
// @Router /user/login [post]
func Login(context *gin.Context) {
	response := model.Response{Context: context}
	var user model.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		logger.Errorf("Unmarshal user fails: %v", err)
		response.Fails("Unmarshal user fails", err)
		return
	}
	encodePassWord := fmt.Sprintf("%x", md5.Sum([]byte(user.PassWord)))
	userDB, err := dao.UserLogin(user.UserName, encodePassWord)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Errorf("用户名或密码错误")
			response.Fails("用户名或密码错误", nil)
			return
		}
		logger.Errorf("login fails: %v", err)
		response.Fails("login fails", err)
		return
	}
	//generate token
	expiresTime := time.Now().Unix() + int64(config.OneDayOfHours)
	claims := jwt.StandardClaims{
		Audience:  userDB.UserName,
		ExpiresAt: expiresTime,
		Id:        strconv.Itoa(int(userDB.ID)),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "LarsWebv0",
		NotBefore: time.Now().Unix(),
		Subject:   "login",
	}
	var jwtSecret = []byte(config.Secret)
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	if err != nil {
		logger.Errorf("generate token fails: %v", err)
		response.Fails("generate token fails", err)
		return
	}
	response.Success(token)
}

// ShowAccount godoc
// @Summary      用户修改个人信息
// @Description  get string by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Success      200  {object}  model.Account
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /user/UpdateUserrInfo [post]
func UpdateUserInfo(context *gin.Context) {
	response := model.Response{Context: context}
	var user model.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		logger.Errorf("Unmarshal user fails: %v", err)
		response.Fails("Unmarshal user fails", err)
		return
	}
	if user.ID <= 0 {
		logger.Errorf("illegal user")
		response.Fails("illegal user", err)
		return
	}
	userDB, err := dao.QueryUserById(user.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Errorf("illegal user")
			response.Fails("illegal user", err)
			return
		} else {
			logger.Errorf("query user fails")
			response.Fails("query user fails", err)
			return
		}
	}
	userHandler(&user, &userDB)
	err = dao.UpdateUserInfo(userDB)
	if err != nil {
		logger.Errorf("update user fails")
		response.Fails("update user fails", err)
		return
	}
	response.Success(nil)
}
func userHandler(user *model.User, userDB *model.User) {
	if user.UserName != userDB.UserName {
		userDB.UserName = user.UserName
	}
	if user.IsStation != userDB.IsStation {
		userDB.IsStation = user.IsStation
	}
}

// @Summary 修改用户个人头像
// @Tags 用户
// @version 1.0
// @Accept application/x-json-stream
// @Param image formData file true "image"
// @Param Authorization header string true "Authorization"
// @Router /user/UpdateUserImage [post]
func UpdateUserImage(context *gin.Context) {
	response := model.Response{Context: context}
	image, err := context.FormFile("image")
	id := middleware.GetIdInToken(context)
	if err != nil {
		logger.Errorf("illegal id: %v", err.Error())
		response.Fails("illegal id", err)
		return
	}
	if err != nil {
		logger.Errorf("upload image fail: %v", err.Error())
		response.Fails("upload image fail", err)
		return
	}
	imageUrl, err := utils.UploadImage(image)
	if err != nil {
		logger.Errorf("upload image to oss fail: %v", err.Error())
		response.Fails("upload image to oss fail", err)
		return
	}
	err = dao.UpdateUserImage(imageUrl, uint(id))
	if err != nil {
		logger.Errorf("update image fail: %v", err.Error())
		response.Fails("update image fail", err)
		return
	}
	response.Success(nil)
}

// @Summary 获取用户个人信息
// @Tags 用户
// @version 1.0
// @Accept application/x-json-stream
// @Param Authorization header string true "Authorization"
// @Router /user/getUserInfo [get]
func GetUserInfo(context *gin.Context) {
	response := model.Response{Context: context}
	userId := middleware.GetIdInToken(context)
	user, err := dao.QueryUserById(userId)
	if err != nil {
		logger.Errorf("get user info fails: %v", err)
		response.Fails("get user info fails", err)
		return
	}
	response.Success(user)
}
