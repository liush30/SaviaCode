package main

import (
	"github.com/gin-gonic/gin"
	"lyods-fabric-demo/enroll/routes"
)

func main() {
	r := gin.Default()

	// 注册路由
	routes.RegisterUserRoutes(r)

	// 启动服务
	if err := r.Run(":8081"); err != nil {
		panic(err)
	}
}
