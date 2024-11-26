package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type PersonalRecord struct {
	BaseInfo    string `json:"baseInfo"`    // 加密后的基本信息 (如姓名、性别、出生日期、血型、身高、体重、过敏史、既往病史等信息)
	ContactInfo string `json:"contactInfo"` // 加密后的联系方式
	Owner       string `json:"owner"`       // 个人档案拥有者
	ApprovalSig string `json:"approvalSig"` // 授权人签名
}

type PersonalRecordContract struct {
	contractapi.Contract
}

// CreatePersonalRecord 创建个人档案
func (p *PersonalRecordContract) CreatePersonalRecord(ctx contractapi.TransactionContextInterface, id, baseInfo, contactInfo, approvalSig string) (string, error) {
	if err := validateInput(id, baseInfo, contactInfo, approvalSig); err != nil {
		return "", err
	}

	// 检查ID是否已存在
	existingRecordJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return "", fmt.Errorf("failed to get state: %v", err)
	}
	if existingRecordJSON != nil {
		return "", fmt.Errorf("record with ID %s already exists", id)
	}
	owner, err := p.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get submitting client identity: %v", err)
	}
	record := &PersonalRecord{
		BaseInfo:    baseInfo,
		ContactInfo: contactInfo,
		ApprovalSig: approvalSig,
		Owner:       owner,
	}

	recordJSON, err := json.Marshal(record)
	if err != nil {
		return "", fmt.Errorf("failed to marshal record: %v", err)
	}

	if err := ctx.GetStub().PutState(id, recordJSON); err != nil {
		return "", fmt.Errorf("failed to put state: %v", err)
	}

	return ctx.GetStub().GetTxID(), nil
}

// QueryPersonalRecord 查询个人档案
func (p *PersonalRecordContract) QueryPersonalRecord(ctx contractapi.TransactionContextInterface, id string) (*PersonalRecord, error) {
	recordJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get state: %v", err)
	}
	if recordJSON == nil {
		return nil, fmt.Errorf("record not found: %s", id)
	}

	record := &PersonalRecord{}
	if err := json.Unmarshal(recordJSON, record); err != nil {
		return nil, fmt.Errorf("failed to unmarshal record: %v", err)
	}

	return record, nil
}

// UpdatePersonalRecord 更新个人档案
func (p *PersonalRecordContract) UpdatePersonalRecord(ctx contractapi.TransactionContextInterface, id string, baseInfo, contactInfo, approvalSig string) (string, error) {
	if err := validateInput(id, baseInfo, contactInfo, approvalSig); err != nil {
		return "", err
	}
	//判断调用者是否为个人档案拥有者
	owner, err := p.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get submitting client identity: %v", err)
	}

	// 检查ID是否存在
	record, err := p.QueryPersonalRecord(ctx, id)
	if err != nil {
		return "", fmt.Errorf("failed to query record: %v", err)
	}

	if record.Owner != owner {
		return "", fmt.Errorf("only owner can update record")
	}

	record.BaseInfo = baseInfo
	record.ContactInfo = contactInfo
	record.ApprovalSig = approvalSig

	recordJSON, err := json.Marshal(record)
	if err != nil {
		return "", fmt.Errorf("failed to marshal record: %v", err)
	}

	if err := ctx.GetStub().PutState(id, recordJSON); err != nil {
		return "", fmt.Errorf("failed to put state: %v", err)
	}

	return ctx.GetStub().GetTxID(), nil
}

// GetSubmittingClientIdentity  获取提交客户端身份
func (p *PersonalRecordContract) GetSubmittingClientIdentity(ctx contractapi.TransactionContextInterface) (string, error) {

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

// validateInput 验证输入参数
func validateInput(id string, baseInfo, contactInfo, approvalSig string) error {
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}
	if baseInfo == "" {
		return fmt.Errorf("baseInfo cannot be empty")
	}
	if contactInfo == "" {
		return fmt.Errorf("contactInfo cannot be empty")
	}
	if approvalSig == "" {
		return fmt.Errorf("approvalSig cannot be empty")
	}
	return nil
}
