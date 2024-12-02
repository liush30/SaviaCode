//go:build pkcs11
// +build pkcs11

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
		infoGroup := patientGroup.Group("/info")
		{
			infoGroup.GET("/hospital", controllers.GetAllHospitals)                      //获取所有医院信息
			infoGroup.GET("/doctor", controllers.GetAllDoctors)                          //获取所有医生信息
			infoGroup.GET("/department/category", controllers.GetAllDepartmentsCategory) //获取医院所有科室类别
			infoGroup.GET("/department", controllers.GetAllDepartments)                  //获取医院所有科室
			infoGroup.GET("/pharmacy", controllers.GetAllPharmacy)                       //获取所有药房信息
		}
		patientGroup.GET("/register", controllers.Registry) //查询个人档案
		medicalGroup := patientGroup.Group("/medical")
		{
			medicalGroup.GET("/end", controllers.EndMedicalInfo)                 //结束就诊
			medicalGroup.GET("/cancel", controllers.CancelMedicalInfo)           //取消就诊
			medicalGroup.GET("/info", controllers.GetPatientMedicalInfoByStatus) //查询指定病患的指定状态的就诊记录
		}
		processGroup := patientGroup.Group("/process")
		{
			processGroup.GET("/getAll", controllers.GetProcessByVisitID)       //查询指定就诊记录的就诊过程记录
			processGroup.POST("/uploadChain", controllers.CreateMedicalRecord) //上链-就诊过程文件
		}
		dispensingGroup := patientGroup.Group("/dispensing")
		{
			dispensingGroup.GET("/create", controllers.CreateDispenseRecord)
			dispensingGroup.GET("/confirm", controllers.ConfirmSignature) //确认签名
			dispensingGroup.GET("/query", controllers.QueryDispensing)    //查询取药单据
		}

	}

}
