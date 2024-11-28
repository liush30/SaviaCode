package routes

import (
	"eldercare_health/app/controllers"
	"eldercare_health/app/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterPatientRoutes(r *gin.Engine)  {
	patientGroup := r.Group("/patient")
	patientGroup.Use(middleware.AuthMiddleware())
	{
		patientGroup.GET("/register/pending", controllers.GetPatientRegisterInfo) //获取患者就诊记录
		patientGroup.GET("/register/active", controllers.)  //获取患者就诊中记录
		patientGroup.GET("/register/accept", controllers.)     //获取患者已就诊记录
		patientGroup.GET("/register/end", controllers.)           //获取患者结束就诊

}
