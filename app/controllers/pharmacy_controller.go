package controllers

import (
	"eldercare_health/app/internal/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetDispenseRecord(c *gin.Context) {
	//获取userId
	userId := c.MustGet("userId").(string)
	//获取page和size
	page := c.Query("page")
	size := c.Query("size")
	//转成int
	iPage, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to convert page to int: " + err.Error()})
		return
	}
	iSize, err := strconv.Atoi(size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to convert size to int: " + err.Error()})
		return
	}
	dbClient, err := db.InitDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to init db: " + err.Error()})
		return
	}
	//调用数据库
	info, err := db.QueryDispensingByPharmacyID(dbClient, userId, iPage, iSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get dispense record: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": info, "total": len(info)})
}
