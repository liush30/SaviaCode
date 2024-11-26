//go:build pkcs11
// +build pkcs11

package fabric

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"log"
)

type PersonalRecord struct {
	BaseInfo    string `json:"baseInfo"`    // 加密后的基本信息 (如姓名、性别、出生日期、血型、身高、体重、过敏史、既往病史等信息)
	ContactInfo string `json:"contactInfo"` // 加密后的联系方式
	ApprovalSig string `json:"approvalSig"` // 授权人签名
}

func CreatePersonalRecord(getaway *client.Gateway, userID, baseInfo, contactInfo, approvalSig string) error {
	contract := getContract(getaway, channelName, personalChaincodeName)
	_, err := contract.SubmitTransaction("CreatePersonalRecord", userID, baseInfo, contactInfo, approvalSig)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %v", err)
	}
	log.Println("Created personal record successfully")
	return nil
}

func QueryPersonalRecord(getaway *client.Gateway, userID string) (*PersonalRecord, error) {
	contract := getContract(getaway, channelName, personalChaincodeName)
	evaluateResult, err := contract.EvaluateTransaction("QueryPersonalRecord", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}
	var record PersonalRecord
	if err := json.Unmarshal(evaluateResult, &record); err != nil {
		return nil, err
	}
	return &record, nil

}

func UpdatePersonalRecord(getaway *client.Gateway, userID, baseInfo, contactInfo, approvalSig string) error {
	contract := getContract(getaway, channelName, personalChaincodeName)
	_, err := contract.SubmitTransaction("UpdatePersonalRecord", userID, baseInfo, contactInfo, approvalSig)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %v", err)
	}
	log.Println("Updated personal record successfully")
	return nil
}
