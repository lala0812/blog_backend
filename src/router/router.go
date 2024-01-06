package router

import (
	"fmt"
	"go_gin_example/controller"
	"go_gin_example/envconfig"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	fmt.Println("收到 GET 請求 - 測試訊息")
	c.JSON(200, gin.H{
		"users": envconfig.GetEnv("USERS"),
	})
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Private-Network", "true")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Category, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// func SetCors() gin.HandlerFunc {
// 	config := cors.DefaultConfig()
// 	config.AllowOrigins = []string{"*"} // 設置你的允許來源
// 	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
// 	config.AllowHeaders = []string{"Content-Type", "Authorization"}
// 	return cors.New(config)
// }

func SetupRoutes(router *gin.Engine) {
	router.Use(CORSMiddleware()) // 設定跨域
	// router.Use(SetCors()) // 設定跨域
	// 定義路由
	router.GET("/", GetUsers) // 讀取首頁

	router.GET("/test", GetUsers) // 測試連線

	router.GET("/articles", controller.GetArticle)                  // 讀取所有文章
	router.GET("/articles/:id", controller.GetArticleById)          // 讀取單篇文章
	router.POST("/articles", controller.CreateArticle)              // 新增文章
	router.PUT("/articles/:id", controller.UpdateArticleById)       // 更新單篇文章
	router.DELETE("/articles/:id", controller.DeleteArticleById)    // 刪除單篇文章
	router.GET("/articles/user/:id", controller.GetArticleByUserId) // 讀取單篇文章的所有留言

	router.GET("/users", controller.GetUser)         // 讀取所有使用者
	router.GET("/users/:id", controller.GetUserById) // 讀取單一使用者
	// router.POST("/users", controller.CreateUser)           // 新增使用者
	router.PUT("/users/:id", controller.UpdateUserById)    // 更新單一使用者
	router.DELETE("/users/:id", controller.DeleteUserById) // 刪除單一使用者

	router.GET("/comments", controller.GetComment)                         // 讀取所有留言
	router.GET("/comments/:id", controller.GetCommentById)                 // 讀取單篇留言
	router.POST("/comments", controller.CreateComment)                     // 新增留言
	router.PUT("/comments/:id", controller.UpdateCommentById)              // 更新單篇留言
	router.DELETE("/comments/:id", controller.DeleteCommentById)           // 刪除單篇留言
	router.GET("/comments/article/:id", controller.GetCommentsByArticleId) // 讀取單篇文章的所有留言

	router.PUT("/changePassword", controller.HandleChangePassword) // 更改密碼
	router.POST("/register", controller.HandleRegister)            // 註冊使用者
	router.POST("/login", controller.HandleLogin)                  // 登入使用者
}
