package routes

import (
	"eldercare_health/enroll/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine) {
	userGroup := r.Group("/users")
	{
		//enroll 路由组
		enrollGroup := userGroup.Group("/enroll")
		{
			enrollGroup.POST("/attributes", controllers.EnrollUserAttributes) // 注册用户属性
			enrollGroup.POST("/identity", controllers.EnrollUser)             // 注册用户身份注册
		}
		registerGroup := userGroup.Group("/register")
		{
			registerGroup.POST("/request", controllers.RegisterRequest) //发起注册请求
			registerGroup.POST("/get", controllers.GetAllUsersRegisterRequest)
			registerGroup.POST("/approve", controllers.ApproveRegistration)
		}
	}
}
