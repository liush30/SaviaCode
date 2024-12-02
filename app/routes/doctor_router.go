//go:build pkcs11
// +build pkcs11

package routes

import (
	"eldercare_health/app/controllers"
	"eldercare_health/app/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterDoctorRoutes(r *gin.Engine) {
	doctorGroup := r.Group("/doctor")
	doctorGroup.Use(middleware.AuthMiddleware())
	{
		medicalGroup := doctorGroup.Group("/medical")
		{
			medicalGroup.GET("/accept", controllers.AcceptRegisterInfo)
			medicalGroup.GET("/info", controllers.GetDoctorRegisterInfoByStatus) //获取指定状态的就诊记录信息
			medicalGroup.GET("/history", controllers.QueryMedicalRecord)
		}
		processGroup := doctorGroup.Group("/process")
		{
			processGroup.POST("/add", controllers.AddProcess)                 //添加过程文件
			processGroup.GET("/get", controllers.GetProcess)                  //获取指定过程内容
			processGroup.GET("/get/medical", controllers.GetProcessByVisitID) //获取指定过程内容
			processGroup.POST("/update", controllers.UpdateProcess)           //更新过程文件
			processGroup.GET("/delete", controllers.DeleteProcess)            //删除过程内容
		}
		//doctorGroup.GET("/update/status", controllers.UpdateDoctorStatus)           //更新状态

		//doctorGroup.GET("/medical/accept", controllers.AcceptRegisterInfo)          //接诊

	}

}
