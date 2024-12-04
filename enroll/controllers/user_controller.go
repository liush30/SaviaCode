package controllers

import (
	"eldercare_health/enroll/internal/db"
	"eldercare_health/enroll/internal/fabric"
	"eldercare_health/enroll/internal/pkg/tool"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// UserAttributeRequest 用户属性注册请求参数
type UserAttributeRequest struct {
	UserID     string           `json:"userId"` // 用户ID
	Attributes []UserAttributes `json:"attributes"`
}

// UserAttributes 用户属性
type UserAttributes struct {
	TaKey     string `json:"taKey"`     // 属性key值
	Attribute string `json:"attribute"` // 属性value值
}

type EnrollUserRequest struct {
	UserName string `json:"userName"` //用户名
	Password string `json:"password"` //密码
}

const InitVersion = 1

// EnrollUser 注册用户
func EnrollUser(c *gin.Context) {
	var enrollRequest EnrollUserRequest
	if err := c.ShouldBindJSON(&enrollRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}
	userID := tool.GenerateUUIDWithoutDashes()
	//调用enroll方法
	err := fabric.Enroll(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user:" + err.Error()})
		return
	}
	// 创建一个更新结构体
	user := db.User{
		Username:     enrollRequest.UserName,
		PasswordHash: enrollRequest.Password,
		CreatedAt:    tool.GetNowTime(),
		UpdatedAt:    tool.GetNowTime(),
		Version:      InitVersion,
	}
	dbClient, err := db.InitDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not connect to database:" + err.Error()})
		return
	}
	err = db.UpdateUser(dbClient, userID, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user:" + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully!", "data": userID})
}

// EnrollUserAttributes 注册用户属性
func EnrollUserAttributes(c *gin.Context) {
	var attributeRequest UserAttributeRequest

	// 绑定请求数据
	if err := c.ShouldBindJSON(&attributeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}
	fmt.Printf("User attributes to insert: %+v\n", attributeRequest)

	var userAttributes []db.UserAttribute
	for _, attr := range attributeRequest.Attributes {
		nowTime := tool.GetNowTime()
		userAttribute := db.UserAttribute{
			TuaKey:    tool.GenerateUUIDWithoutDashes(),
			UserID:    attributeRequest.UserID,
			TaKey:     attr.TaKey,
			Attribute: attr.Attribute,
			CreatedAt: nowTime,
			UpdatedAt: nowTime,
		}
		userAttributes = append(userAttributes, userAttribute)
	}
	dbClient, err := db.InitDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user attribute"})
		return
	}
	// 保存用户属性到数据库
	if err := db.BatchInsertUserAttributes(dbClient, userAttributes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user attribute"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User attribute created successfully"})
}

func DeleteUserAttributes(c *gin.Context) {
	//获取url中的key
	key := c.Query("key")
	dbClient, err := db.InitDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete user attribute"})
		return
	}
	if err := db.DeleteUserAttribute(dbClient, key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete user attribute"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User attribute deleted successfully"})

}

func GetUserAttributes(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get user attribute"})
		return
	}
	var conditions = make(map[string]interface{})
	err = c.ShouldBindJSON(&conditions)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	data, err := db.GetUserAttributeByCondition(dbClient, conditions, pageInt, sizeInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get user attribute"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data, "total": len(data)})

}
