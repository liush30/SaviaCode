package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type MedicalRecordContract struct {
	contractapi.Contract
}

// MedicalRecord  就诊记录信息
type MedicalRecord struct {
	RecordID   string `json:"record_id"`  // 就诊记录
	DoctorID   string `json:"doctorId"`   // 医生ID
	Data       string `json:"data"`       // 加密数据
	Timestamp  int64  `json:"timestamp"`  // 记录时间 (时间戳)
	DoctorSign string `json:"doctorSign"` // 医生签名 (数字签名)
	RecordType string `json:"recordType"` // 记录类型
	Owner      string `json:"owner"`      // 数据拥有者
}

// CreateMedicalRecord 创建就诊记录
func (m *MedicalRecordContract) CreateMedicalRecord(ctx contractapi.TransactionContextInterface, personID, recordID, recordType, doctorID, doctorSign, data string) error {
	err := checkInput(personID, recordID, recordType, doctorSign, doctorID, data)
	if err != nil {
		return err
	}
	//生成复合键
	key, err := ctx.GetStub().CreateCompositeKey(personID, []string{recordID, recordType, doctorID})
	if err != nil {
		return fmt.Errorf("failed to create composite key: %v", err)
	}
	exist, err := checkRecordExist(ctx, key)
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("record with key %s already exists", key)
	}
	// 获取交易的时间戳
	timestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get transaction timestamp: %v", err)
	}
	owner, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get client identity: %v", err)
	}
	// 将时间戳转换为更易读的格式
	txTime := timestamp.Seconds
	record := MedicalRecord{
		DoctorID:   doctorID,
		Data:       data,
		Timestamp:  txTime,
		RecordID:   recordID,
		DoctorSign: doctorSign,
		RecordType: recordType,
		Owner:      owner,
	}
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal record: %v", err)
	}
	// 将记录存储在链上
	err = ctx.GetStub().PutState(key, recordJSON)
	if err != nil {
		return fmt.Errorf("failed to put state: %v", err)
	}
	return nil
}

func (m *MedicalRecordContract) UpdateMedicalRecord(ctx contractapi.TransactionContextInterface, personID, recordID, recordType, doctorID, doctorSign, data string) error {
	err := checkInput(personID, recordID, recordType, doctorSign, doctorID, data)
	if err != nil {
		return err
	}
	//生成复合键
	key, err := ctx.GetStub().CreateCompositeKey(personID, []string{recordID, recordType, doctorID})
	if err != nil {
		return fmt.Errorf("failed to create composite key: %v", err)
	}
	//查询记录是否存在
	state, err := ctx.GetStub().GetState(key)
	if err != nil {
		return err
	}
	if state == nil {
		return fmt.Errorf("record with key %s does not exist", key)
	}
	var record MedicalRecord
	if err = json.Unmarshal(state, &record); err != nil {
		return fmt.Errorf("failed to unmarshal record: %v", err)
	}
	//获取调用者身份
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get client identity: %v", err)
	}
	if clientID != record.Owner {
		return fmt.Errorf("client identity does not match")
	}
	// 获取交易的时间戳
	timestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get transaction timestamp: %v", err)
	}
	// 将时间戳转换为更易读的格式
	record.Data = data
	record.Timestamp = timestamp.Seconds
	record.DoctorSign = doctorSign
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal record: %v", err)
	}
	// 将记录存储在链上
	err = ctx.GetStub().PutState(key, recordJSON)
	if err != nil {
		return fmt.Errorf("failed to put state: %v", err)
	}
	return nil
}

// QueryMedicalRecord 查询患者的所有就诊记录
func (m *MedicalRecordContract) QueryMedicalRecord(ctx contractapi.TransactionContextInterface, personID string) ([]MedicalRecord, error) {
	if personID == "" {
		return nil, fmt.Errorf("personID cannot be empty")
	}
	// 读取链上记录
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(personID, []string{})
	if err != nil {
		return nil, fmt.Errorf("failed to get state by partial composite key: %v", err)
	}
	defer resultsIterator.Close()
	var records []MedicalRecord
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next result: %v", err)
		}
		var record MedicalRecord
		if err := json.Unmarshal(queryResponse.Value, &record); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
		}
		records = append(records, record)
	}
	return records, nil
}

// QueryPrescription 查询指定类型的就诊记录
//
// personID: 患者ID
// recordID: 就诊记录ID
// recordType: 就诊记录类型
func (m *MedicalRecordContract) QueryPrescription(ctx contractapi.TransactionContextInterface, personID, recordID, recordType string) (*MedicalRecord, error) {
	if personID == "" {
		return nil, fmt.Errorf("personID cannot be empty")
	}

	if recordID == "" {
		return nil, fmt.Errorf("recordID cannot be empty")
	}

	if recordType == "" {
		return nil, fmt.Errorf("recordType cannot be empty")
	}
	// 读取链上记录
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(personID, []string{recordID, recordType})
	if err != nil {
		return nil, fmt.Errorf("failed to get state by partial composite key: %v", err)
	}
	defer resultsIterator.Close()
	var record *MedicalRecord
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next result: %v", err)
		}
		var r MedicalRecord
		if err := json.Unmarshal(queryResponse.Value, &r); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
		}
		record = &r
	}
	return record, nil
}

// GetSubmittingClientIdentity  获取提交客户端身份
func (m *MedicalRecordContract) GetSubmittingClientIdentity(ctx contractapi.TransactionContextInterface) (string, error) {

	b64ID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("failed to read clientID: %v", err)
	}
	decodeID, err := base64.StdEncoding.DecodeString(b64ID)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode clientID: %v", err)
	}
	return string(decodeID), nil
}

func checkInput(personID, recordID, recordType, doctorID, doctorSign, data string) error {
	if personID == "" {
		return fmt.Errorf("personID cannot be empty")
	}

	if recordID == "" {
		return fmt.Errorf("recordID cannot be empty")
	}

	if recordType == "" {
		return fmt.Errorf("recordType cannot be empty")
	}

	if doctorID == "" {
		return fmt.Errorf("doctorID cannot be empty")
	}

	if data == "" {
		return fmt.Errorf("data cannot be empty")
	}

	if doctorSign == "" {
		return fmt.Errorf("doctorSign cannot be empty")
	}

	return nil
}

// checkRecordExist 检查记录是否已经存在
// true 存在 false 不存在
func checkRecordExist(ctx contractapi.TransactionContextInterface, key string) (bool, error) {
	// 判断该记录是否已经存在
	existingRecord, err := ctx.GetStub().GetState(key)
	if err != nil {
		return false, fmt.Errorf("failed to get state: %v", err)
	}
	return existingRecord != nil, nil
}
