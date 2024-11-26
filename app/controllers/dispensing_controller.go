//go:build pkcs11
// +build pkcs11

package controllers

import (
	"github.com/gin-gonic/gin"
	"lyods-fabric-demo/app/internal/fabric"
	"lyods-fabric-demo/app/internal/pkg/tool"
	"net/http"
)

type CreateDispensingRequest struct {
	//UserID                string `json:"userId"`                // 用户ID
	PrescriptionID        string `json:"prescriptionId"`        // 处方ID
	PharmacyID            string `json:"pharmacyId"`            // 药房ID
	PatientID             string `json:"patientId"`             // 患者ID
	ScheduledDispenseTime string `json:"scheduledDispenseTime"` // 预定取药时间
}

// CreateDispenseRecord 初始化取药单据
func CreateDispenseRecord(c *gin.Context) {
	var request CreateDispensingRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}
	cert, mspID, err := getCertAndMspID(request.PatientID)
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
	getaway, err := fabric.SetupGateway(mspID, cert, hsmSign)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get contract: " + err.Error()})
		return
	}
	defer getaway.Close()
	err = fabric.CreatePrescriptionDispensing(getaway, request.PrescriptionID, request.PharmacyID, request.PatientID, request.ScheduledDispenseTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create prescription dispensing: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Created prescription dispensing successfully"})
}

// QueryDispensing 查询取药单据
func QueryDispensing(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	userId := c.Query("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}
	cert, mspID, err := getCertAndMspID(userId)
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
	defer getaway.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get contract: " + err.Error()})
		return
	}
	record, err := fabric.QueryDispensing(getaway, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query prescription dispensing: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, record)

}

const (
	signTypePharmacy = "pharmacy"
	signTypePatient  = "patient"
)

// ConfirmSignature 确认签名
func ConfirmSignature(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	userId := c.Query("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}
	signType := c.Query("signType")
	if signType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "signType is required"})
		return
	}
	cert, mspID, err := getCertAndMspID(userId)
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
	sign, err := hsmSign([]byte(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign: " + err.Error()})
		return
	}
	getaway, err := fabric.SetupGateway(mspID, cert, hsmSign)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get contract: " + err.Error()})
		return
	}
	defer getaway.Close()
	if signType == signTypePharmacy {
		err = fabric.ConfirmPharmacySignature(getaway, id, tool.EncodeToString(sign))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to confirm pharmacy signature: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Confirmed pharmacy signature successfully"})
		return
	} else if signType == signTypePatient {
		err = fabric.ConfirmPatientSignature(getaway, id, tool.EncodeToString(sign))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to confirm patient signature: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Confirmed patient signature successfully"})
		return

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "signType is invalid"})
		return
	}
}
