package controllers

import (
	"eldercare_health/app/internal/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	medicalActive   = "营业中"
	medicalInActive = "暂停营业"
	medicalClose    = "已关闭"
	medicalHospital = "医院"
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
