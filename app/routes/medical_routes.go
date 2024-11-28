//go:build pkcs11
// +build pkcs11

package routes

import (
	"eldercare_health/app/controllers"
	"eldercare_health/app/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterMedicalRoutes(r *gin.Engine) {
	medicalGroup := r.Group("/medicalRecords")
	medicalGroup.Use(middleware.AuthMiddleware())
	{
		medicalGroup.GET("/register", controllers.RegistryMedicalRecord) // 根据就诊ID查询就诊记录
		medicalGroup.POST("/create", controllers.CreateMedicalRecord)    // 创建就诊记录
		medicalGroup.GET("/queryAll", controllers.QueryMedicalRecord)    // 根据就诊ID查询所有相关就诊记录
		medicalGroup.POST("/queryType", controllers.QueryPrescription)   // 根据就诊ID查询指定类型的就诊记录
		medicalGroup.POST("/update", controllers.UpdateMedicalRecord)    // 更新就诊记录
		infoGroup := medicalGroup.Group("/info")
		{
			infoGroup.GET("/get/hospital", controllers.GetAllHospitals)                     //获取所有医院信息
			infoGroup.GET("/get/doctor", controllers.GetAllDoctors)                         //获取所有医生信息
			infoGroup.GET("/get/departmentCategory", controllers.GetAllDepartmentsCategory) //获取医院所有科室类别
			infoGroup.GET("/get/department", controllers.GetAllDepartments)                 //获取医院所有科室
			infoGroup.POST("/create", controllers.CreateMedicalInfo)                        // 创建就诊信息
			infoGroup.POST("/update", controllers.UpdateMedicalInfo)                        // 更新就诊信息
		}
	}

}
