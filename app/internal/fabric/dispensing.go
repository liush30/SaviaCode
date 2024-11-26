//go:build pkcs11
// +build pkcs11

package fabric

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"log"
)

type DispenseRecord struct {
	PrescriptionID        string `json:"prescriptionId"`        // 处方ID
	PharmacyID            string `json:"pharmacyId"`            // 药房ID
	PatientID             string `json:"patientId"`             // 患者ID
	PharmacySign          string `json:"pharmacySign"`          // 药房签名（数字签名）
	PatientSign           string `json:"patientSign"`           // 患者签名（数字签名）
	ScheduledDispenseTime string `json:"scheduledDispenseTime"` // 预定取药时间
	ActualDispenseTime    string `json:"actualDispenseTime"`    // 实际取药时间
	Status                int64  `json:"status"`                // 取药状态，例如 0-"已取药"、1-"待取药"
}

// CreatePrescriptionDispensing 创建取药单据
func CreatePrescriptionDispensing(getaway *client.Gateway, prescriptionID, pharmacyID, patientID, scheduledDispenseTime string) error {
	contract := getContract(getaway, channelName, dispensingChaincodeName)
	_, err := contract.SubmitTransaction("CreatePrescriptionDispensing", prescriptionID, pharmacyID, patientID, scheduledDispenseTime)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %v", err)
	}
	log.Println("Created prescription dispensing successfully")
	return nil
}

func QueryDispensing(getaway *client.Gateway, prescriptionID string) (DispenseRecord, error) {
	contract := getContract(getaway, channelName, dispensingChaincodeName)
	data, err := contract.EvaluateTransaction("QueryDispensing", prescriptionID)
	if err != nil {
		return DispenseRecord{}, fmt.Errorf("failed to evaluate transaction: %v", err)
	}

	var record DispenseRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return DispenseRecord{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}
	return record, nil
}

func ConfirmPharmacySignature(getaway *client.Gateway, prescriptionID, pharmacySign string) error {
	contract := getContract(getaway, channelName, dispensingChaincodeName)
	// 编码数据
	_, err := contract.SubmitTransaction("ConfirmPharmacySignature", prescriptionID, pharmacySign)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %v", err)
	}
	log.Println("Confirmed pharmacy signature successfully")
	return nil
}

func ConfirmPatientSignature(getaway *client.Gateway, prescriptionID, patientSign string) error {
	contract := getContract(getaway, channelName, dispensingChaincodeName)
	// 编码数据
	_, err := contract.SubmitTransaction("ConfirmPatientSignature", prescriptionID, patientSign)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %v", err)
	}
	return nil
}
