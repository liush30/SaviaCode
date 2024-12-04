//go:build pkcs11
// +build pkcs11

package fabric

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// MedicalRecord  病历信息
type MedicalRecord struct {
	RecordID   string `json:"record_id"`  // 就诊记录
	DoctorID   string `json:"doctorId"`   // 医生ID
	Data       string `json:"data"`       // 加密数据
	Timestamp  int64  `json:"timestamp"`  // 记录时间 (时间戳)
	DoctorSign string `json:"doctorSign"` // 医生签名 (数字签名)
	RecordType string `json:"recordType"` // 记录类型
}

const (
	RecordTypePrescription = "处方"        // 处方
	RecordTypeDiagnosis    = "diagnosis" // 诊断
	RecordTypeTreatment    = "treatment" // 诊疗
	RecordTypeProcedure    = "procedure" // 治疗方案
)

// CreateMedicalRecord 创建就诊记录
func CreateMedicalRecord(getaway *client.Gateway, personID, recordID, recordType, doctorID, doctorSign, data string) error {
	contract := getContract(getaway, channelName, medicalChaincodeName)
	_, err := contract.SubmitTransaction("CreateMedicalRecord", personID, recordID, recordType, doctorID, doctorSign, data)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %v", err)
	}
	return nil
}

// UpdateMedicalRecord 更新就诊记录
func UpdateMedicalRecord(getaway *client.Gateway, personID, recordID, recordType, doctorID, doctorSign, data string) error {
	contract := getContract(getaway, channelName, medicalChaincodeName)
	_, err := contract.SubmitTransaction("UpdateMedicalRecord", personID, recordID, recordType, doctorID, doctorSign, data)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %v", err)
	}
	return nil
}

// QueryMedicalRecord 查询患者的所有就诊记录
func QueryMedicalRecord(getaway *client.Gateway, personID string) (map[string][]MedicalRecord, error) {
	contract := getContract(getaway, channelName, medicalChaincodeName)
	evaluateResult, err := contract.EvaluateTransaction("QueryMedicalRecord", personID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %v", err)
	}
	if evaluateResult == nil {
		return nil, nil
	}
	var records []MedicalRecord
	if err := json.Unmarshal(evaluateResult, &records); err != nil {
		return nil, err
	}

	classifiedRecords := make(map[string][]MedicalRecord)
	// 遍历记录并按 Record 分类
	for _, record := range records {
		classifiedRecords[record.RecordID] = append(classifiedRecords[record.RecordID], record)
	}
	return classifiedRecords, nil
}

// QueryPrescription 查询处方
func QueryPrescription(getaway *client.Gateway, personID, recordID, recordType string) (MedicalRecord, error) {
	contract := getContract(getaway, channelName, medicalChaincodeName)
	evaluateResult, err := contract.EvaluateTransaction("QueryPrescription", personID, recordID, recordType)
	if err != nil {
		return MedicalRecord{}, err
	}
	if evaluateResult == nil {
		return MedicalRecord{}, nil
	}
	var record MedicalRecord
	if err := json.Unmarshal(evaluateResult, &record); err != nil {
		return MedicalRecord{}, fmt.Errorf("failed to unmarshal JSON: %v,json:%s", err, evaluateResult)
	}

	return record, nil
}
