package middleware

import (
	"LarsWebV0/config"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		auth := context.Request.Header.Get("Authorization")
		if len(auth) == 0 {
			context.Abort()
			context.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "无法认证，重新登录",
			})
			return
		}
		//校验token
		_, err := parseToken(auth)
		if err != nil {
			context.Abort()
			context.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "token 过期" + err.Error(),
			})
			return
		} else {
			logger.Info("auto pass")
		}
		context.Next()
	}
}

func parseToken(token string) (*jwt.StandardClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Secret), nil
	})
	if err != nil && jwtToken != nil {
		if claims, ok := jwtToken.Claims.(*jwt.StandardClaims); ok && jwtToken.Valid {
			return claims, nil
		}
	}
	return nil, err
}

func GetIdInToken(context *gin.Context) uint {
	token := context.Request.Header.Get("Authorization")
	jwtToken, _ := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Secret), nil
	})
	idStr := jwtToken.Claims.(*jwt.StandardClaims).Id
	id, _ := strconv.ParseInt(idStr, 10, 32)
	return uint(id)
}
