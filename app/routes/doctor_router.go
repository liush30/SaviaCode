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
		doctorGroup.GET("/update/status", controllers.UpdateDoctorStatus)           //更新状态
		doctorGroup.GET("/medical/info", controllers.GetDoctorRegisterInfoByStatus) //获取指定状态的就诊记录信息
		doctorGroup.GET("/medical/accept", controllers.AcceptRegisterInfo)          //接诊
		doctorGroup.POST("/process/add", controllers.AddProcess)                    //添加过程文件
		doctorGroup.GET("/process/get", controllers.GetProcess)                     //获取指定过程内容
		doctorGroup.GET("/process/getAll", controllers.GetProcessByVisitID)         //获取指定过程内容
		doctorGroup.POST("/process/update", controllers.UpdateProcess)              //更新过程文件
		doctorGroup.GET("/process/delete", controllers.DeleteProcess)               //删除过程内容
	}

}
