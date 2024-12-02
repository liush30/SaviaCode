//go:build pkcs11
// +build pkcs11

package controllers

import (
	"eldercare_health/app/internal/db"
	"eldercare_health/app/internal/pkg/tool"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	medicalActive   = "营业中"
	medicalInActive = "暂停营业"
	medicalClose    = "已关闭"
	medicalHospital = "医院"
	medicalPharmacy = "药房"
)

// GetAllHospitals 获取所有正常营业的医院
func GetAllHospitals(c *gin.Context) {
	initDB, err := db.InitDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not connect to database:" + err.Error()})
		return
	}
	facility, err := db.QueryMedicalFacility(initDB, medicalHospital, medicalActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get facility:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": facility})
}

func GetAllDepartmentsCategory(c *gin.Context) {
	initDB, err := db.InitDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not connect to database:" + err.Error()})
		return
	}
	//获取id
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	//获取指定id的数据
	facility, err := db.GetCategoriesByHospitalID(initDB, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get facility:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": facility})
}
func GetAllDepartments(c *gin.Context) {
	initDB, err := db.InitDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not connect to database:" + err.Error()})
		return
	}
	//获取id
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	category := c.Query("category")
	if category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category is required"})
		return
	}
	//获取指定id的数据
	facility, err := db.GetDepartNameAndIdByHospitalID(initDB, id, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get facility:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": facility})
}

func GetAllDoctors(c *gin.Context) {
	initDB, err := db.InitDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not connect to database:" + err.Error()})
		return
	}
	//获取id
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	//获取指定id的数据
	department, err := db.GetDoctorsByDepartmentID(initDB, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get facility:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": department})
}

// Registry 挂号
func Registry(c *gin.Context) {
	//获取userId
	userId := c.MustGet("userId").(string)
	//获取doctor_id
	doctorId := c.Query("doctor_id")
	if doctorId == "" || userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId or doctor_id is required"})
		return
	}
	dbClient, err := db.InitDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not connect to database:" + err.Error()})
		return
	}
	//查看指定doctor的信息
	doctor, err := db.GetDoctor(dbClient, doctorId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get doctor:" + err.Error()})
		return
	}
	maxNumber := doctor.MaxNumber
	//判断是否超过最大挂号数量
	number, err := db.QueryByDoctorAndDatePrefix(dbClient, doctorId, time.Now().Format(time.DateOnly))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get records:" + err.Error()})
		return
	}
	if number+1 > maxNumber {
		c.JSON(http.StatusBadRequest, gin.H{"error": "over max number"})
		return
	}
	record := db.MedicalRecord{
		TmrID:     tool.GenerateUUIDWithoutDashes(),
		PatientID: userId,
		DoctorID:  doctorId,
		Status:    recordStatusPend,
		CreateAt:  tool.GetNowTime(),
		UpdateAt:  tool.GetNowTime(),
		Version:   1,
	}
	err = db.CreateMedicalRecord(dbClient, &record)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "success"})
}
