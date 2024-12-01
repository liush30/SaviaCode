package routes

import (
	"eldercare_health/app/controllers"
	"eldercare_health/app/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterPatientRoutes(r *gin.Engine) {
	patientGroup := r.Group("/patient")
	patientGroup.Use(middleware.AuthMiddleware())
	{
		//结束就诊
		patientGroup.GET("/medical/end", controllers.EndMedicalInfo)                 //结束就诊
		patientGroup.GET("/medical/info", controllers.GetPatientMedicalInfoByStatus) //查询指定病患的指定状态的就诊记录
		patientGroup.GET("/medical/cancel", controllers.CancelMedicalInfo)           //取消就诊
		patientGroup.GET("/process/getAll", controllers.GetProcessByVisitID)         //查询指定就诊记录的就诊过程记录
		patientGroup.POST("/medical/uploadChain", controllers.CreateMedicalRecord)   //上链-就诊过程文件
		patientGroup.GET("/get/pharmacy", controllers.GetAllPharmacy)                //获取所有药房信息
	}

}
