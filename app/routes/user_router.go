package routes

import "github.com/gin-gonic/gin"

func RegisterUserRoutes(r *gin.Engine) {
	userGroup := r.Group("/users")
}
