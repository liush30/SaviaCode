package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/fentec-project/bn256"
	"github.com/fentec-project/gofe/abe"
	"log"
	"lyods-fabric-demo/app/internal/db"
	"lyods-fabric-demo/app/internal/pkg/tool"
)

type OnChainData struct {
	TEDKey           string
	EncryptedMessage []byte // 加密后的消息
	IV               []byte // 初始化向量
	Hash             string // 消息哈希
}

// Encrypt 加密交易数据
//
// userID: 用户 ID
// msg: 消息
// exp: 属性访问规则表达式
// auth: 涉及授权机构
func Encrypt(userID, msg, exp string, auth []string) ([]byte, error) {
	// 创建新的 MAABE 结构，带有全局参数
	maabe := abe.NewMAABE()
	// 从布尔公式创建 MSP 结构
	msp, err := abe.BooleanToMSP(exp, false)
	if err != nil {
		return nil, fmt.Errorf("failed to create MSP: %v", err)
	}
	// 生成公钥
	pks := generatePublicKey(auth)
	if err != nil {
		return nil, fmt.Errorf("failed to generate public keys: %v", err)
	}
	// 使用 MSP 中的解密策略加密消息
	ct, err := maabe.Encrypt(msg, msp, pks)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt: %v", err)
	}
	c0Byte, err := json.Marshal(ct.C0)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal C0: %v", err)
	}
	c1xByte, err := json.Marshal(ct.C1x)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal C1x: %v", err)
	}
	c2xByte, err := json.Marshal(ct.C2x)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal C2x: %v", err)
	}
	c3xByte, err := json.Marshal(ct.C3x)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal C3x: %v", err)
	}
	mspByte, err := json.Marshal(ct.Msp)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal MSP: %v", err)
	}
	timeNow := tool.GetNowTime()
	id := tool.GenerateUUIDWithoutDashes()
	log.Println("generate encrypted id: ", id)
	encryptedData := db.EncryptedData{
		TEDKey:         id,
		UserID:         userID,
		C0:             c0Byte,
		C1X:            c1xByte,
		C2X:            c2xByte,
		C3X:            c3xByte,
		MSP:            mspByte,
		CreateDate:     timeNow,
		LastModifyDate: timeNow,
	}
	dbClient, err := db.InitDB()
	if err != nil {
		return nil, err
	}
	// 保存加密数据到数据库
	err = db.CreateEncryptedData(dbClient, encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to create encrypted data: %v", err)
	}
	// 计算加密数据的hash值
	hash := sha256.Sum256(ct.SymEnc)
	chainData := OnChainData{
		TEDKey:           id,
		EncryptedMessage: ct.SymEnc,
		IV:               ct.Iv,
		Hash:             hex.EncodeToString(hash[:]),
	}
	chainDataByte, err := json.Marshal(chainData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chain data: %v", err)
	}
	return chainDataByte, nil
}

// Decrypt 解密
func Decrypt(data []byte, userID string) (string, error) {
	// 创建新的 MAABE 结构，带有全局参数
	maabe := abe.NewMAABE()

	var chainDate OnChainData
	if err := json.Unmarshal(data, &chainDate); err != nil {
		return "", fmt.Errorf("failed to unmarshal chain data: %v", err)
	}
	dbClient, err := db.InitDB()
	if err != nil {
		return "", err
	}
	//根据id查询数据信息
	encryptedData, err := db.GetEncryptedData(dbClient, chainDate.TEDKey)
	if err != nil {
		return "", fmt.Errorf("failed to get encrypted data: %v", err)
	}
	var c0 *bn256.GT
	if err := json.Unmarshal(encryptedData.C0, &c0); err != nil {
		return "", fmt.Errorf("failed to unmarshal C0: %v", err)
	}

	var c1x map[string]*bn256.GT
	if err := json.Unmarshal(encryptedData.C1X, &c1x); err != nil {
		return "", fmt.Errorf("failed to unmarshal C1x: %v", err)
	}

	var c2x map[string]*bn256.G2
	if err := json.Unmarshal(encryptedData.C2X, &c2x); err != nil {
		return "", fmt.Errorf("failed to unmarshal C2x: %v", err)
	}

	var c3x map[string]*bn256.G2
	if err := json.Unmarshal(encryptedData.C3X, &c3x); err != nil {
		return "", fmt.Errorf("failed to unmarshal C3x: %v", err)
	}

	var msp *abe.MSP
	if err := json.Unmarshal(encryptedData.MSP, &msp); err != nil {
		return "", fmt.Errorf("failed to unmarshal MSP: %v", err)
	}
	// 使用 C0，C1x，C2x，C3x，MSP 中的解密策略解密消息
	ct := &abe.MAABECipher{
		SymEnc: chainDate.EncryptedMessage,
		Iv:     chainDate.IV,
		C0:     c0,
		C1x:    c1x,
		C2x:    c2x,
		C3x:    c3x,
		Msp:    msp,
	}

	userKeys, err := generateUserKey(userID)
	if err != nil {
		return "", fmt.Errorf("failed to generate user keys: %v", err)
	}
	msg, err := maabe.Decrypt(ct, userKeys)
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}
	return msg, nil
}

// 定义公钥集合
func generatePublicKey(auths []string) []*abe.MAABEPubKey {
	var pubKeys []*abe.MAABEPubKey
	for _, auth := range auths {
		pubKeys = append(pubKeys, mapAuth[auth].PubKeys())
	}
	return pubKeys
}
