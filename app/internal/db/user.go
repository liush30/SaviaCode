package db

import (
	"fmt"
	"gorm.io/gorm"
)

type User struct {
	UserID         string `gorm:"column:USER_ID;primaryKey" json:"userId"`
	Username       string `gorm:"column:USERNAME" json:"username"`
	PasswordHash   string `gorm:"column:PASSWORD_HASH" json:"passwordHash"`
	MspID          string `gorm:"column:MSP_ID" json:"mspId"`
	EnrollmentCert []byte `gorm:"column:ENROLLMENT_CERT" json:"enrollmentCert"`
	IdentityID     string `gorm:"column:IDENTITY_ID" json:"identityId"`
	CreatedAt      string `gorm:"column:CREATED_AT" json:"createdAt"`
	UpdatedAt      string `gorm:"column:UPDATED_AT" json:"updatedAt"`
	Version        int    `gorm:"column:VERSION;default:1" json:"version"`
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
