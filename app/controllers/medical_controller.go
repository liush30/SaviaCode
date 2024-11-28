package controllers

import (
	"eldercare_health/app/internal/crypto"
	"eldercare_health/app/internal/db"
	"eldercare_health/app/internal/pkg/tool"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type MedicalRecordRequest struct {
	UserID string `json:"userId"` // 用户ID
	//PersonID   string   `json:"personId"`   // 患者ID
	RecordID   string   `json:"recordId"`   // 就诊记录
	RecordType string   `json:"recordType"` // 记录类型
	DoctorID   string   `json:"doctorId"`   // 医生ID
	DoctorSign string   `json:"doctorSign"` // 医生签名
	Data       string   `json:"data"`       // 加密数据
	CryptoExp  string   `json:"cryptoExp"`  // 规则加密策略
	Auth       []string `json:"auth"`       // 授权机构列表
}

const (
	recordStatusPend   = "待就诊"
	recordStatusDone   = "就诊结束"
	recordStatusActive = "就诊中"
)

func RegistryMedicalRecord(c *gin.Context) {
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

// CreateMedicalRecord 初始化就诊记录
func CreateMedicalRecord(c *gin.Context) {
	var request MedicalRecordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}
	cert, mspID, err := getCertAndMspID(request.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get cert and msp id: " + err.Error()})
		return
	}
	//创建HSMSignerFactory
	hsmSignerFactory, err := fabric.CreateHSMSignerFactory()
	defer hsmSignerFactory.Dispose()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create hsm signer factory: " + err.Error()})
		return
	}
	//创建HSMSign
	hsmSign, hsmSignClose, err := fabric.CreateHSMSign(hsmSignerFactory, cert)
	defer hsmSignClose()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create hsm sign: " + err.Error()})
		return
	}
	//计算数据hash
	//hashData := tool.CalculateSHA256Hash(request.Data)
	//加密数据
	encrypt, err := crypto.Encrypt(request.UserID, request.Data, request.CryptoExp, request.Auth)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encrypt: " + err.Error()})
		return
	}
	getaway, err := fabric.SetupGateway(mspID, cert, hsmSign)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get contract: " + err.Error()})
		return
	}
	defer getaway.Close()
	err = fabric.CreateMedicalRecord(getaway, request.UserID, request.RecordID, request.RecordType, request.DoctorID, request.DoctorSign, tool.EncodeToString(encrypt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create medical record: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Created medical record successfully"})
}

// UpdateMedicalRecord 更新就诊记录
func UpdateMedicalRecord(c *gin.Context) {
	var request MedicalRecordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}
	cert, mspID, err := getCertAndMspID(request.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get cert and msp id: " + err.Error()})
		return
	}
	//创建HSMSignerFactory
	hsmSignerFactory, err := fabric.CreateHSMSignerFactory()
	defer hsmSignerFactory.Dispose()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create hsm signer factory: " + err.Error()})
		return
	}
	//创建HSMSign
	hsmSign, hsmSignClose, err := fabric.CreateHSMSign(hsmSignerFactory, cert)
	defer hsmSignClose()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create hsm sign: " + err.Error()})
		return
	}
	//计算数据hash
	//hashData := tool.CalculateSHA256Hash(request.Data)
	//加密数据
	encrypt, err := crypto.Encrypt(request.UserID, request.Data, request.CryptoExp, request.Auth)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encrypt: " + err.Error()})
		return
	}
	getaway, err := fabric.SetupGateway(mspID, cert, hsmSign)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get contract: " + err.Error()})
		return
	}
	defer getaway.Close()
	err = fabric.UpdateMedicalRecord(getaway, request.UserID, request.RecordID, request.RecordType, request.DoctorID, request.DoctorSign, tool.EncodeToString(encrypt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create medical record: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Created medical record successfully"})
}

// QueryMedicalRecord 查询患者的所有就诊记录
func QueryMedicalRecord(c *gin.Context) {
	// 获取要查询的用户的ID
	userID := c.Query("userId")
	//申请查询的请求者ID
	requesterID := c.Query("requesterId")
	if userID == "" || requesterID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId and requesterID is required"})
		return
	}
	cert, mspID, err := getCertAndMspID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get cert and msp id: " + err.Error()})
		return
	}
	//创建HSMSignerFactory
	hsmSignerFactory, err := fabric.CreateHSMSignerFactory()
	defer hsmSignerFactory.Dispose()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create hsm signer factory: " + err.Error()})
		return
	}
	//创建HSMSign
	hsmSign, hsmSignClose, err := fabric.CreateHSMSign(hsmSignerFactory, cert)
	defer hsmSignClose()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create hsm sign: " + err.Error()})
		return
	}
	//创建Gateway
	getaway, err := fabric.SetupGateway(mspID, cert, hsmSign)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get contract: " + err.Error()})
		return
	}
	defer getaway.Close()
	//根据用户id查询就诊记录
	records, err := fabric.QueryMedicalRecord(getaway, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query medical record: " + err.Error()})
		return
	}
	//解密数据
	for key, value := range records {
		for i, record := range value {
			dataBytes, err := tool.DecodeToString(record.Data)
			if err != nil {
				log.Printf("failed to decode data for record %v: %s", record, err.Error())
				continue
			}
			decrypt, err := crypto.Decrypt(dataBytes, requesterID)
			if err != nil {
				log.Printf("failed to decrypt data for record %v: %s", record, err.Error())
				continue // 如果解密失败，则跳过此记录，不修改数据
			}

			records[key][i].Data = decrypt
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Query personal record successfully", "data": records})
}

type QueryPrescriptionRequest struct {
	RequesterID string `json:"requesterId"`
	UserId      string `json:"userId"`
	RecordID    string `json:"recordId"`
	RecordType  string `json:"recordType"`
}

// QueryPrescription 查询指定类型的用户就诊记录,并根据请求者的id对加密数据进行解密，若解密成功直接返回解密后的数据。
func QueryPrescription(c *gin.Context) {
	var request QueryPrescriptionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}
	//检查参数
	if request.RequesterID == "" || request.UserId == "" || request.RecordID == "" || request.RecordType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId or recordId or recordType is required"})
		return
	}

	cert, mspID, err := getCertAndMspID(request.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get cert and msp id: " + err.Error()})
		return
	}
	//创建HSMSignerFactory
	hsmSignerFactory, err := fabric.CreateHSMSignerFactory()
	defer hsmSignerFactory.Dispose()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create hsm signer factory: " + err.Error()})
		return
	}
	//创建HSMSign
	hsmSign, hsmSignClose, err := fabric.CreateHSMSign(hsmSignerFactory, cert)
	defer hsmSignClose()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create hsm sign: " + err.Error()})
		return
	}
	//创建Gateway
	getaway, err := fabric.SetupGateway(mspID, cert, hsmSign)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get contract: " + err.Error()})
		return
	}
	defer getaway.Close()
	// 获取链上数据
	record, err := fabric.QueryPrescription(getaway, request.UserId, request.RecordID, request.RecordType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query prescription: " + err.Error()})
		return
	}
	//fmt.Println(record.Data)
	//解密数据
	dataBytes, err := tool.DecodeToString(record.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode data: " + err.Error()})
	}
	decrypt, err := crypto.Decrypt(dataBytes, request.RequesterID)
	if err == nil {
		record.Data = decrypt
	} else {
		log.Println("failed to decrypt data: ", err)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Query prescription successfully", "data": record})
}
