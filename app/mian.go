//go:build pkcs11
// +build pkcs11

package main

import (
	"eldercare_health/app/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	if err := r.Run(":8089"); err != nil {
		panic(err)
	}
}
