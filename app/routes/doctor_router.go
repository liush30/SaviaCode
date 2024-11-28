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
		doctorGroup.GET("/register/pending", controllers.GetDoctorRegisterInfo) //获取医生待就诊记录
		doctorGroup.GET("/register/active", controllers.GetActiveRegisterInfo)  //获取医生就诊中记录
		doctorGroup.GET("/register/accept", controllers.AcceptRegisterInfo)     //获取医生已就诊记录
		doctorGroup.GET("/register/end", controllers.EndRegisterInfo)           //获取医生结束就诊
		doctorGroup.POST("/process/add", controllers.AddProcess)                //添加过程文件
		doctorGroup.GET("/process/get", controllers.GetProcess)                 //获取指定过程内容
		doctorGroup.GET("/process/getAll", controllers.GetProcessByVisitID)     //获取指定过程内容
		doctorGroup.POST("/process/update", controllers.UpdateProcess)          //更新过程文件
		doctorGroup.GET("/process/delete", controllers.DeleteProcess)
	}

}
