//go:build pkcs11
// +build pkcs11

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

func GetPatientMedicalInfoByStatus(c *gin.Context) {
	//获取user_id
	userId := c.MustGet("userId").(string)
	//获取状态
	status := c.Query("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status is required"})
		return
	}
	//获取page和size
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
	dbClient, err := db.InitDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not connect to database:" + err.Error()})
		return
	}
	conditions := make(map[string]interface{})
	conditions["patient_id"] = userId
	if status != "" {
		conditions["status"] = status
	}
	//获取就诊记录
	info, err := db.QueryMedicalRecordByConditions(dbClient, conditions, pageInt, sizeInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get records:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": info, "total": len(info)})

}
func EndMedicalInfo(c *gin.Context) {
	//获取record_id
	recordID := c.Query("recordId")
	if recordID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "recordId is required"})
		return
	}
	dbClient, err := db.InitDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not connect to database:" + err.Error()})
		return
	}
	//获取
	err = db.UpdateMedicalRecordStatus(dbClient, recordID, recordStatusDone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get records:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

func CancelMedicalInfo(c *gin.Context) {
	//获取record_id
	recordID := c.Query("recordId")
	if recordID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "recordId is required"})
		return
	}
	dbClient, err := db.InitDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not connect to database:" + err.Error()})
		return
	}
	//获取
	err = db.UpdateMedicalRecordStatus(dbClient, recordID, recordStatusCancel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get records:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

// GetAllPharmacy  获取所有正常营业的药房
func GetAllPharmacy(c *gin.Context) {
	initDB, err := db.InitDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not connect to database:" + err.Error()})
		return
	}
	//获取指定id的数据
	department, err := db.QueryMedicalFacility(initDB, medicalPharmacy, medicalActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get facility:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": department})
}

//上传个人档案信息
