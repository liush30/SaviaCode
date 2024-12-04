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

			enrollGroup.POST("/identity", controllers.EnrollUser) // 注册用户身份注册
		}
		registerGroup := userGroup.Group("/register")
		{
			registerGroup.POST("/request", controllers.RegisterRequest) //发起注册请求
			registerGroup.POST("/get", controllers.GetAllUsersRegisterRequest)
			registerGroup.GET("/approve", controllers.ApproveRegistration)
		}
		attrGroup := userGroup.Group("/attributes")
		{
			attrGroup.POST("/register", controllers.EnrollUserAttributes) // 注册用户属性
			attrGroup.GET("/delete", controllers.DeleteUserAttributes)    // 删除用户属性
			attrGroup.POST("/get", controllers.GetUserAttributes)         // 查询用户属性
		}
	}
}
