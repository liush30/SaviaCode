//go:build pkcs11
// +build pkcs11

package routes

import (
	"github.com/gin-gonic/gin"
	"lyods-fabric-demo/app/controllers"
)

func RegisterMedicalRoutes(r *gin.Engine) {
	medicalGroup := r.Group("/medicalRecords")
	{
		medicalGroup.POST("/create", controllers.CreateMedicalRecord)  // 创建就诊记录
		medicalGroup.GET("/queryAll", controllers.QueryMedicalRecord)  // 根据就诊ID查询所有相关就诊记录
		medicalGroup.POST("/queryType", controllers.QueryPrescription) // 根据就诊ID查询指定类型的就诊记录
		medicalGroup.POST("/update", controllers.UpdateMedicalRecord)  // 更新就诊记录
		infoGroup := medicalGroup.Group("/info")
		{
			infoGroup.POST("/create", controllers.CreateMedicalInfo) // 创建就诊信息
			infoGroup.POST("/update", controllers.UpdateMedicalInfo) // 更新就诊信息
		}
	}

}
