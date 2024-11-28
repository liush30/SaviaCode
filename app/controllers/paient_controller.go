package controllers

import (
	"eldercare_health/app/internal/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetPatientAllRegisterInfo 获取患者所有就诊记录
func GetPatientAllRegisterInfo(c *gin.Context) {
	dbClient, err := db.InitDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not connect to database:" + err.Error()})
		return
	}
	//获取page 和 size
	page := c.Query("page")
	size := c.Query("size")
	if size == "" || page == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get page and size"})
		return
	}

	//将page和size转成int
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to convert page to int"})
		return
	}
	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to convert size to int"})
		return
	}
	//获取user_id
	userID := c.MustGet("userId").(string)
	condition := make(map[string]interface{})
	condition["patient_id"] = userID
	//获取就诊记录
	info, err := db.QueryMedicalRecordByConditions(dbClient, condition, pageInt, sizeInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get records:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": info, "total": len(info)})
}
func GetPatientActiveRegisterInfo(c *gin.Context) {
	dbClient, err := db.InitDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not connect to database:" + err.Error()})
		return
	}
	//获取page 和 size
	page := c.Query("page")
	size := c.Query("size")
	if size == "" || page == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get page and size"})
		return
	}
	//获取user_id
	userID := c.MustGet("userId").(string)
	condition := make(map[string]interface{})
	condition["patient_id"] = userID
	condition["status"] = recordStatusActive
	//获取就诊记录
	info, err := db.QueryMedicalRecordByConditionsNoOffset(dbClient, condition)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get records:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": info})
}
