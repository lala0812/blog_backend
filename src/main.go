// File: main.go
package main

import (
	"go_gin_example/router"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// 創建一個新的mux（路由器）
	server := gin.Default()

	// 創建CORS處理程序
	// server.Use(router.SetCors())
	// server.Use(router.CORSMiddleware())

	// 設置路由
	router.SetupRoutes(server)

	// 啟動伺服器
	http.ListenAndServe(":8080", server)
}
