//go:build pkcs11
// +build pkcs11

package main

import (
	"github.com/gin-gonic/gin"
	"lyods-fabric-demo/app/routes"
)

func main() {
	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowHeaders = []string{"Authorization", "Content-Type"}
	r.Use(cors.New(corsConfig))
	routes.RegisterPharmacy(r)
	routes.RegisterDoctorRoutes(r)
	routes.RegisterPatientRoutes(r)
	routes.RegisterUserRoutes(r)

	// 启动服务
	r.Run()
}
