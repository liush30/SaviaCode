//go:build pkcs11
// +build pkcs11

package routes

import (
	"eldercare_health/app/controllers"
	"eldercare_health/app/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterPharmacy(r *gin.Engine) {
	pharmacyGroup := r.Group("/pharmacy")
	pharmacyGroup.Use(middleware.AuthMiddleware())
	{
		pharmacyGroup.GET("/info", controllers.GetDispenseRecord) //获取指定状态的就诊记录信息
		pharmacyGroup.GET("/confirm", controllers.ConfirmSignature)
		pharmacyGroup.GET("/prescription/get", controllers.QueryPrescription)
	}

}
