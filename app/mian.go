//go:build pkcs11
// +build pkcs11

package main

import (
	"github.com/gin-gonic/gin"
	"lyods-fabric-demo/app/routes"
)

func main() {
	r := gin.Default()

	// 注册个人档案路由
	routes.RegisterPersonalRecordRoutes(r)

	// 注册取药单据路由
	routes.RegisterDispensingRoutes(r)

	// 注册就诊记录路由
	routes.RegisterMedicalRoutes(r)

	// 启动服务
	r.Run()
}
