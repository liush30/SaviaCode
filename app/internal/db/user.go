package db

import (
	"eldercare_health/app/internal/pkg/tool"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type User struct {
	UserID         string `gorm:"column:USER_ID;primaryKey" json:"userId"`
	Username       string `gorm:"column:USERNAME" json:"username"`
	IdNumber       string `gorm:"column:ID_NUMBER" json:"idNumber"`
	PasswordHash   string `gorm:"column:PASSWORD_HASH" json:"passwordHash"`
	MspID          string `gorm:"column:MSP_ID" json:"mspId"`
	EnrollmentCert []byte `gorm:"column:ENROLLMENT_CERT" json:"enrollmentCert"`
	IdentityID     string `gorm:"column:IDENTITY_ID" json:"identityId"`
	CreatedAt      string `gorm:"column:CREATED_AT" json:"createdAt"`
	UpdatedAt      string `gorm:"column:UPDATED_AT" json:"updatedAt"`
	Version        int    `gorm:"column:VERSION;default:1" json:"version"`
	Type           string `gorm:"column:TYPE" json:"type"`
}

func (User) TableName() string {
	return "t_users"
}

// GetUserCertAndMspID 获取用户的 ENROLLMENT_CERT 和 MSP_ID
func GetUserCertAndMspID(db *gorm.DB, userID string) (enrollmentCert []byte, mspID string, err error) {
	var result struct {
		EnrollmentCert []byte `gorm:"column:ENROLLMENT_CERT" json:"enrollmentCert"`
		MspID          string `gorm:"column:MSP_ID" json:"mspId"`
	}

	if err := db.Table("t_users").
		Select("ENROLLMENT_CERT, MSP_ID").
		Where("USER_ID = ?", userID).
		Scan(&result).Error; err != nil {
		return nil, "", err
	}
	// 判断用户是否存在
	if result.EnrollmentCert == nil && result.MspID == "" {
		return nil, "", fmt.Errorf("user ID %s does not exist", userID)
	}
	return result.EnrollmentCert, result.MspID, nil
}

// Login 根据身份证号码和密码进行用户验证，返回用户的 ID
func Login(db *gorm.DB, idNumber, password string) (string, string, error) {
	var user User

	// 根据身份证号码查找用户
	err := db.Where("id_number = ? AND password = ?", idNumber, tool.CalculateSHA256Hash(password)).First(&user).Error
	if err != nil {
		return "", "", errors.New("invalid id number or password")
	}

	// 返回ID
	return user.UserID, user.Type, nil
}
