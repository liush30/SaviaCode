package routes

//func RegisterMedicalRoutes(r *gin.Engine) {
//	medicalGroup := r.Group("/medicalRecords")
//	medicalGroup.Use(middleware.AuthMiddleware())
//	{
//		//medicalGroup.GET("/register", controllers.RegistryMedicalRecord) // 根据就诊ID查询就诊记录
//		//medicalGroup.POST("/create", controllers.CreateMedicalRecord)    // 创建就诊记录
//		medicalGroup.GET("/queryAll", controllers.QueryMedicalRecord)  // 根据就诊ID查询所有相关就诊记录
//		medicalGroup.POST("/queryType", controllers.QueryPrescription) // 根据就诊ID查询指定类型的就诊记录
//		medicalGroup.POST("/update", controllers.UpdateMedicalRecord)  // 更新就诊记录
//	}
//
//}
