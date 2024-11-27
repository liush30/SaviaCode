package controllers

import (
	"eldercare_health/enroll/internal/db"
	"eldercare_health/enroll/internal/fabric"
	"eldercare_health/enroll/internal/pkg/tool"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func RegisterRequest(c *gin.Context) {
	var req db.UserRegistration
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}
	req.TurID = tool.GenerateUUIDWithoutDashes()
	req.CreateAt = tool.GetNowTime()
	req.UpdateAt = tool.GetNowTime()
	req.Version = 1
	req.Password = tool.CalculateSHA256Hash(req.Password)

	dbClient, err := db.InitDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not connect to database:" + err.Error()})
		return
	}
	err = db.CreateUserRegistration(dbClient, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user:" + err.Error()})
		return
	}
}
func GetAllUsersRegisterRequest(c *gin.Context) {
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
	//获取请求数据
	reqParams := map[string]interface{}{}
	err = c.ShouldBindJSON(reqParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}
	dbClient, err := db.InitDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not connect to database:" + err.Error()})
		return
	}
	conditions, err := db.GetAllRegistrationWithConditions(dbClient, reqParams, pageInt, sizeInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get all registration:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": conditions, "total": len(conditions)})
}

func ApproveRegistration(c *gin.Context) {
	//获取结果
	result := c.Query("result")
	if result == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "result is required"})
		return
	}
	//获取turID
	turID := c.Query("turID")
	if turID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "turID is required"})
		return
	}
	dbClient, err := db.InitDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not connect to database:" + err.Error()})
		return
	}
	//获取指定id的数据
	registration, err := db.GetRegistration(dbClient, turID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get registration:" + err.Error()})
		return
	}
	if registration == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "registration not found"})
		return
	}
	updateResult := make(map[string]interface{})
	updateResult["result"] = result
	updateResult["updateAt"] = tool.GetNowTime()
	updateResult["version"] = registration.Version + 1
	//更新数据
	err = db.UpdateRegistration(dbClient, turID, result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update registration:" + err.Error()})
		return
	}
	if result == "通过" {
		//调用enroll方法
		err := fabric.Enroll(turID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user:" + err.Error()})
			return
		}
		// 创建一个更新结构体
		user := db.User{
			Username:     registration.Name,
			PasswordHash: registration.Password,
			CreatedAt:    tool.GetNowTime(),
			UpdatedAt:    tool.GetNowTime(),
			Version:      InitVersion,
		}
		err = db.UpdateUser(dbClient, turID, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user:" + err.Error()})
			return
		}
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully!"})
}
