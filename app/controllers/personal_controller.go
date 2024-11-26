//go:build pkcs11
// +build pkcs11

package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"log"
	"lyods-fabric-demo/app/internal/crypto"
	"lyods-fabric-demo/app/internal/db"
	"lyods-fabric-demo/app/internal/fabric"
	"lyods-fabric-demo/app/internal/pkg/tool"
	"net/http"
)

type CreatePersonalRecordRequest struct {
	UserID           string   `json:"userID"`           // 用户ID
	BaseInfo         string   `json:"baseInfo"`         // 基本信息
	ContactInfo      string   `json:"contactInfo"`      // 联系方式
	BaseCryptoExp    string   `json:"cryptoExp"`        // 基本信息的规则加密策略
	BaseAuth         []string `json:"auth"`             // 基本信息的授权机构列表
	ContactCryptoExp string   `json:"contactCryptoExp"` // 联系方式的规则加密策略
	ContactAuth      []string `json:"contactAuth"`      // 联系方式的授权机构列表
}

// 根据userid从数据库中获取
func getCertAndMspID(userID string) ([]byte, string, error) {
	//连接数据库
	dbClient, err := db.InitDB()
	if err != nil {
		return nil, "", err
	}
	cert, mspID, err := db.GetUserCertAndMspID(dbClient, userID)
	if err != nil {
		return nil, "", err
	}
	return cert, mspID, nil
}
func CreatePersonalRecord(c *gin.Context) {
	var request CreatePersonalRecordRequest
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
	//创建Gateway
	getaway, err := fabric.SetupGateway(mspID, cert, hsmSign)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get contract: " + err.Error()})
		return
	}
	defer getaway.Close()
	encryptBase, encryptContact, sign, err := processPersonalRecord(request, hsmSign, getaway)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process personal record: " + err.Error()})
		return
	}
	err = fabric.CreatePersonalRecord(getaway, request.UserID, encryptBase, encryptContact, sign)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create personal record: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Created personal record successfully"})
}

func UpdatePersonalRecord(c *gin.Context) {
	var request CreatePersonalRecordRequest
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
	//创建Gateway
	getaway, err := fabric.SetupGateway(mspID, cert, hsmSign)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get contract: " + err.Error()})
		return
	}
	defer getaway.Close()
	encryptBase, encryptContact, sign, err := processPersonalRecord(request, hsmSign, getaway)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process personal record: " + err.Error()})
		return
	}
	err = fabric.UpdatePersonalRecord(getaway, request.UserID, encryptBase, encryptContact, sign)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create personal record: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Created personal record successfully"})
}

// 处理个人记录的创建逻辑
func processPersonalRecord(request CreatePersonalRecordRequest, hsmSign identity.Sign, getaway *client.Gateway) (string, string, string, error) {
	hashData := tool.CalculateSHA256Hash(request.BaseInfo)       // 计算基本信息的哈希值
	hashContact := tool.CalculateSHA256Hash(request.ContactInfo) // 计算联系方式的哈希值
	// 签名
	signature, err := hsmSign([]byte(hashData + hashContact))
	if err != nil {
		return "", "", "", fmt.Errorf("failed to sign: %w", err)
	}
	sign := tool.EncodeToString(signature) // 签名

	// 加密基本信息
	encryptBase, err := crypto.Encrypt(request.UserID, request.BaseInfo, request.BaseCryptoExp, request.BaseAuth)
	if err != nil {
		return "", "", "", err
	}

	// 加密联系方式
	encryptContact, err := crypto.Encrypt(request.UserID, request.ContactInfo, request.ContactCryptoExp, request.ContactAuth)
	if err != nil {
		return "", "", "", err
	}
	return tool.EncodeToString(encryptBase), tool.EncodeToString(encryptContact), sign, nil
}

func QueryPersonalRecord(c *gin.Context) {
	//获取userID和requesterID
	userID := c.Query("userId")
	requesterID := c.Query("requesterId")
	if userID == "" || requesterID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId and requesterId are required"})
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
	record, err := fabric.QueryPersonalRecord(getaway, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query personal record: " + err.Error()})
		return
	}
	err = respondWithRecord(record, requesterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to respond with record: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Query personal record successfully", "data": record})
}

// 响应查询到的个人记录，并进行必要的解密
func respondWithRecord(record *fabric.PersonalRecord, requesterID string) error {
	//解密数据
	baseBytes, err := tool.DecodeToString(record.BaseInfo)
	if err != nil {
		return fmt.Errorf("failed to decode base info: %w", err)
	}
	contactBytes, err := tool.DecodeToString(record.ContactInfo)
	if err != nil {
		return fmt.Errorf("failed to decode contact info: %w", err)
	}
	baseDec, err := crypto.Decrypt(baseBytes, requesterID)
	if err != nil {
		log.Println("failed to decrypt base info: ", err)
	} else {
		record.BaseInfo = baseDec
	}
	contractDec, err := crypto.Decrypt(contactBytes, requesterID)
	if err != nil {
		log.Println("failed to decrypt contact info: ", err)
	} else {
		record.ContactInfo = contractDec
	}

	return nil
}
