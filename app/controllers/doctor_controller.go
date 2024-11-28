package controllers

import (
	"eldercare_health/app/internal/db"
	"eldercare_health/app/internal/pkg/tool"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetDoctorRegisterInfo(c *gin.Context) {
	//获取user_id
	userID := c.Query("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
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
	//获取
	records, err := db.QueryMedicalRecordByDoctorID(dbClient, userID, recordStatusPend, pageInt, sizeInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get records:" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": records})

}

func GetActiveRegisterInfo(c *gin.Context) {
	//获取user_id
	userID := c.Query("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
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
	//获取
	records, err := db.QueryMedicalRecordByDoctorID(dbClient, userID, recordStatusActive, pageInt, sizeInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get records:" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": records, "total": len(records)})
}

func AcceptRegisterInfo(c *gin.Context) {
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
	err = db.UpdateMedicalRecordStatus(dbClient, recordID, recordStatusActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get records:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

func EndRegisterInfo(c *gin.Context) {
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

type ProcessReq struct {
	VisitID     string `json:"visitId"`
	RecordType  string `json:"recordType"`
	RecordValue string `json:"recordValue"`
}

func AddProcess(c *gin.Context) {
	var req ProcessReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}
	dbClient, err := db.InitDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not connect to database:" + err.Error()})
		return
	}
	process := db.MedicalProcess{
		ProcessID:   tool.GenerateUUIDWithoutDashes(),
		VisitID:     req.VisitID,
		RecordType:  req.RecordType,
		RecordValue: req.RecordValue,
		CreateAt:    tool.GetNowTime(),
	}
	err = db.CreateMedicalProcess(dbClient, &process)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

// GetProcessByVisitID 获取指定就诊记录的就诊过程信息
func GetProcessByVisitID(c *gin.Context) {
	visitID := c.Query("recordId")
	if visitID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "visitId is required"})
		return
	}
	dbClient, err := db.InitDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not connect to database:" + err.Error()})
		return
	}
	process, err := db.GetMedicalProcess(dbClient, visitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get records:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": process})
}

func GetProcess(c *gin.Context) {
	processId := c.Query("processId")
	if processId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "processId is required"})
		return
	}
	dbClient, err := db.InitDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not connect to database:" + err.Error()})
		return
	}
	process, err := db.GetMedicalProcessByID(dbClient, processId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get records:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": process})
}

func DeleteProcess(c *gin.Context) {
	processId := c.Query("processId")
	if processId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "processId is required"})
		return
	}
	dbClient, err := db.InitDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not connect to database:" + err.Error()})
		return
	}
	err = db.DeleteMedicalProcess(dbClient, processId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get records:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "success"})

}
