package router

import (
	"LarsWebV0/middleware"
	"LarsWebV0/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.LoadHTMLFiles("templates/page/home.html", "templates/page/userInfo.html")
	router.GET("/pages/home", func(c *gin.Context) {
		c.HTML(200, "home.html", nil)
	})
	router.GET("/pages/userInfo", func(c *gin.Context) {
		c.HTML(200, "userInfo.html", nil)
	})
	// 提供 CSS 文件的静态文件服务
	router.Static("/css", "./templates/css")
	// 提供 JavaScript 文件的静态文件服务
	router.Static("/js", "./templates/js")

	url := ginSwagger.URL("http://127.0.0.1:8080/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	router.MaxMultipartMemory = 8 << 20
	router.Use(middleware.Cors())
	router.GET("/ping", service.Ping)

	user := router.Group("/user")
	user.Use(middleware.Cors())
	{
		user.POST("/login", service.Login)
		user.POST("/register", service.Register)
		user.POST("/UpdateUserInfo", service.UpdateUserInfo).Use(middleware.Auth())
		user.POST("/UpdateUserImage", service.UpdateUserImage).Use(middleware.Auth())
		user.GET("/GetUserInfo", service.GetUserInfo).Use(middleware.Auth())

	}

	article := router.Group("/article")
	article.Use(middleware.Cors())
	article.Use(middleware.Auth())
	{
		article.POST("/insert", service.InsertArticle)
	}
	return router
}
