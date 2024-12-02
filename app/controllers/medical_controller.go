//go:build pkcs11
// +build pkcs11

package controllers

import (
	"eldercare_health/app/internal/crypto"
	"eldercare_health/app/internal/db"
	"eldercare_health/app/internal/pkg/tool"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

type MedicalRecordRequest struct {
	//UserID    string `json:"userId"` // 用户ID
	ProcessID string `json:"processId"`
	//PersonID   string   `json:"personId"`   // 患者ID
	//RecordID   string `json:"recordId"`   // 就诊记录id
	//RecordType string `json:"recordType"` // 记录类型
	//DoctorID string `json:"doctorId"` // 医生ID
	//DoctorSign string   `json:"doctorSign"` // 医生签名
	//Data      string   `json:"data"`      // 加密数据
	ExpID     string   `json:"expId"`     // 策略ID
	CryptoExp string   `json:"cryptoExp"` // 规则加密策略
	Auth      []string `json:"auth"`      // 授权机构列表
}

const (
	recordStatusPend   = "待就诊"
	recordStatusDone   = "就诊结束"
	recordStatusActive = "就诊中"
	recordStatusCancel = "就诊取消"
)

// CreateMedicalRecord 初始化就诊记录
func CreateMedicalRecord(c *gin.Context) {
	var request MedicalRecordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}
	//获取user_id
	userID := c.MustGet("userId").(string)
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	dbClient, err := db.InitDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not connect to database:" + err.Error()})
		return
	}
	//根据process id获取process
	process, err := db.GetMedicalProcessByID(dbClient, request.ProcessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get process:" + err.Error()})
		return
	}
	record, err := db.GetMedicalRecord(dbClient, process.VisitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get record:" + err.Error()})
		return
	}
	//判断expid是否为空
	if request.ExpID != "" {
		expInfo, err := db.GetCryptoExpByID(dbClient, request.ExpID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get exp info:" + err.Error()})
			return
		}
		request.CryptoExp = expInfo.Exp
		request.Auth = strings.Split(expInfo.Auth, ",")
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
	encrypt, err := crypto.Encrypt(userID, process.RecordValue, request.CryptoExp, request.Auth)
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
	err = fabric.CreateMedicalRecord(getaway, userID, process.VisitID, process.RecordType, record.DoctorID, process.Sign, tool.EncodeToString(encrypt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create medical record: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Created medical record successfully"})
}

// UpdateMedicalRecord 更新就诊记录
//func UpdateMedicalRecord(c *gin.Context) {
//	var request MedicalRecordRequest
//	if err := c.ShouldBindJSON(&request); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
//		return
//	}
//	cert, mspID, err := getCertAndMspID(request.UserID)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get cert and msp id: " + err.Error()})
//		return
//	}
//	//创建HSMSignerFactory
//	hsmSignerFactory, err := fabric.CreateHSMSignerFactory()
//	defer hsmSignerFactory.Dispose()
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create hsm signer factory: " + err.Error()})
//		return
//	}
//	//创建HSMSign
//	hsmSign, hsmSignClose, err := fabric.CreateHSMSign(hsmSignerFactory, cert)
//	defer hsmSignClose()
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create hsm sign: " + err.Error()})
//		return
//	}
//	//计算数据hash
//	//hashData := tool.CalculateSHA256Hash(request.Data)
//	//加密数据
//	encrypt, err := crypto.Encrypt(request.UserID, request.Data, request.CryptoExp, request.Auth)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encrypt: " + err.Error()})
//		return
//	}
//	getaway, err := fabric.SetupGateway(mspID, cert, hsmSign)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get contract: " + err.Error()})
//		return
//	}
//	defer getaway.Close()
//	err = fabric.UpdateMedicalRecord(getaway, request.UserID, request.RecordID, request.RecordType, request.DoctorID, tool.EncodeToString(encrypt))
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create medical record: " + err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{"message": "Created medical record successfully"})
//}

// QueryMedicalRecord 查询患者的所有就诊记录
func QueryMedicalRecord(c *gin.Context) {
	// 获取要查询的用户的ID
	patientId := c.Query("patientId")
	doctorId := c.MustGet("userId").(string)
	if patientId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "patientId is required"})
		return
	}
	cert, mspID, err := getCertAndMspID(doctorId)
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
	records, err := fabric.QueryMedicalRecord(getaway, patientId)
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
			decrypt, err := crypto.Decrypt(dataBytes, doctorId)
			if err != nil {
				log.Printf("failed to decrypt data for record %v: %s", record, err.Error())
				continue // 如果解密失败，则跳过此记录，不修改数据
			}

			records[key][i].Data = decrypt
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": records, "total": len(records)})
}

type QueryPrescriptionRequest struct {
	//RequesterID string `json:"requesterId"`
	UserId     string `json:"userId"`
	RecordID   string `json:"recordId"`
	RecordType string `json:"recordType"`
}

// QueryPrescription 查询用户就诊记录,并根据请求者的id对加密数据进行解密，若解密成功直接返回解密后的数据。
func QueryPrescription(c *gin.Context) {
	//var request QueryPrescriptionRequest
	//if err := c.ShouldBindJSON(&request); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
	//	return
	//}
	userID := c.MustGet("userId").(string)
	//获取record id
	recordId := c.Query("recordId")
	if recordId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "recordId is required"})
		return
	}
	dbClient, err := db.InitDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to init db: " + err.Error()})
		return
	}
	medicalRecord, err := db.GetMedicalRecord(dbClient, recordId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get medical record: " + err.Error()})
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
	// 获取链上数据
	record, err := fabric.QueryPrescription(getaway, medicalRecord.PatientID, recordId, fabric.RecordTypePrescription)
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
	decrypt, err := crypto.Decrypt(dataBytes, userID)
	if err == nil {
		record.Data = decrypt
	} else {
		log.Println("failed to decrypt data: ", err)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Query prescription successfully", "data": record})
}
