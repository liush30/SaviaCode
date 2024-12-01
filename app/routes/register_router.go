package routes

import (
	"eldercare_health/app/controllers"
	"github.com/gin-gonic/gin"
)

// RegisterRegisterRoutes 挂号
func RegisterRegisterRoutes(r *gin.Engine) {
	registerGroup := r.Group("/register")
	{
		registerGroup.GET("/get/hospital", controllers.GetAllHospitals)                     //获取所有医院信息
		registerGroup.GET("/get/doctor", controllers.GetAllDoctors)                         //获取所有医生信息
		registerGroup.GET("/get/departmentCategory", controllers.GetAllDepartmentsCategory) //获取医院所有科室类别
		registerGroup.GET("/get/department", controllers.GetAllDepartments)                 //获取医院所有科室
		registerGroup.GET("/add", controllers.Registry)                                     //挂号
	}

}
