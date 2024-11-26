//go:build pkcs11
// +build pkcs11

package routes

import (
	"github.com/gin-gonic/gin"
	"lyods-fabric-demo/app/controllers"
)

func RegisterDispensingRoutes(r *gin.Engine) {
	dispensingGroup := r.Group("/dispensing")
	{
		dispensingGroup.POST("/create", controllers.CreateDispenseRecord)
		dispensingGroup.GET("/query", controllers.QueryDispensing)
		dispensingGroup.GET("/confirmSignature", controllers.ConfirmSignature)
	}
}
