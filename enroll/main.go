package main

import (
	"eldercare_health/enroll/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowHeaders = []string{"Authorization", "Content-Type"}
	r.Use(cors.New(corsConfig))
	// 注册路由
	routes.RegisterUserRoutes(r)

	// 启动服务
	if err := r.Run(":8088"); err != nil {
		panic(err)
	}
}
