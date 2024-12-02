package controllers

import (
	"eldercare_health/app/internal/db"
	"eldercare_health/app/internal/pkg/tool"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginInfo struct {
	IDNumber string `json:"id_number"` //身份证号码
	Password string `json:"password"`  //密码
}

func Login(c *gin.Context) {
	var info LoginInfo
	err := c.ShouldBind(&info)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind:" + err.Error()})
	}
	dbClient, err := db.InitDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "connection failed:" + err.Error()})
		return
	}
	id, userType, err := db.Login(dbClient, info.IDNumber, info.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login:" + err.Error()})
		return
	}
	if id == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user does not exist"})
		return
	}
	token, err := tool.GenerateJWT(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "type": userType, "data": "success"})
}
