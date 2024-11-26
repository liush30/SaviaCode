package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"time"
)

type DispenseContract struct {
	contractapi.Contract
}

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
func (d *DispenseContract) CreatePrescriptionDispensing(ctx contractapi.TransactionContextInterface, prescriptionID, pharmacyID, patientID, scheduledDispenseTime string) (string, error) {
	//判断记录是否已经存在
	exists, err := d.checkRecordExists(ctx, prescriptionID)
	if exists {
		return "", fmt.Errorf("prescription dispensing record already exists: %s", prescriptionID)
	}
	record := &DispenseRecord{
		PrescriptionID:        prescriptionID,
		PharmacyID:            pharmacyID,
		PatientID:             patientID,
		ScheduledDispenseTime: scheduledDispenseTime,
		Status:                0,
	}
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return "", fmt.Errorf("failed to marshal prescription dispensing record: %v", err)
	}

	if err := ctx.GetStub().PutState(prescriptionID, recordJSON); err != nil {
		return "", fmt.Errorf("failed to put state: %v", err)
	}
	return ctx.GetStub().GetTxID(), nil
}

// QueryDispensing 查询取药单据
func (d *DispenseContract) QueryDispensing(ctx contractapi.TransactionContextInterface, prescriptionID string) (DispenseRecord, error) {
	recordJSON, err := ctx.GetStub().GetState(prescriptionID)
	if err != nil {
		return DispenseRecord{}, fmt.Errorf("failed to get state: %v", err)
	}
	if recordJSON == nil {
		return DispenseRecord{}, fmt.Errorf("record not found: %s", prescriptionID)
	}
	var record DispenseRecord
	if err := json.Unmarshal(recordJSON, &record); err != nil {
		return DispenseRecord{}, fmt.Errorf("failed to unmarshal record: %v", err)
	}
	return record, nil
}

// ConfirmPharmacySignature 确认药房签名
func (d *DispenseContract) ConfirmPharmacySignature(ctx contractapi.TransactionContextInterface, prescriptionID, signature string) error {
	if signature == "" {
		return fmt.Errorf("empty signature")
	}
	//查询记录是否存在
	record, err := d.QueryDispensing(ctx, prescriptionID)
	if err != nil {
		return err
	}
	if record.PharmacySign != "" {
		return fmt.Errorf("pharmacy signature already exists: %s", prescriptionID)
	}
	//decodedSign, err := base64.StdEncoding.DecodeString(signature)
	//if err != nil {
	//	return fmt.Errorf("failed to decode doctorSign: %v", err)
	//}
	record.PharmacySign = signature
	//判断患者是否签名,若双方均已签名，表示完成取药
	if record.PatientSign != "" {
		// 获取交易的时间戳
		timestamp, err := ctx.GetStub().GetTxTimestamp()
		if err != nil {
			return fmt.Errorf("failed to get transaction timestamp: %v", err)
		}
		record.ActualDispenseTime = timestamp.AsTime().Format(time.RFC3339)
		record.Status = 1
	}
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal record: %v", err)
	}
	if err := ctx.GetStub().PutState(prescriptionID, recordJSON); err != nil {
		return fmt.Errorf("failed to put state: %v", err)
	}
	return nil
}

// ConfirmPatientSignature 确认取药者签名
func (d *DispenseContract) ConfirmPatientSignature(ctx contractapi.TransactionContextInterface, prescriptionID, signature string) error {
	if signature == "" {
		return fmt.Errorf("empty signature")
	}
	//查询记录是否存在
	record, err := d.QueryDispensing(ctx, prescriptionID)
	if err != nil {
		return err
	}
	if record.PatientSign != "" {
		return fmt.Errorf("patient signature already exists: %s", prescriptionID)
	}
	//decodedSign, err := base64.StdEncoding.DecodeString(signature)
	//if err != nil {
	//	return fmt.Errorf("failed to decode doctorSign: %v", err)
	//}
	record.PatientSign = signature
	//判断患者是否签名,若双方均已签名，表示完成取药
	if record.PharmacySign != "" {
		// 获取交易的时间戳
		timestamp, err := ctx.GetStub().GetTxTimestamp()
		if err != nil {
			return fmt.Errorf("failed to get transaction timestamp: %v", err)
		}
		record.ActualDispenseTime = timestamp.AsTime().Format(time.RFC3339)
		record.Status = 1
	}
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal record: %v", err)
	}
	if err := ctx.GetStub().PutState(prescriptionID, recordJSON); err != nil {
		return fmt.Errorf("failed to put state: %v", err)
	}
	return nil
}

// 检查记录是否存在
func (d *DispenseContract) checkRecordExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	recordJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to get state: %v", err)
	}
	return recordJSON != nil, nil
}
