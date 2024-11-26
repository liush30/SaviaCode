package db

import "gorm.io/gorm"

type User struct {
	UserID         string `gorm:"column:USER_ID;primaryKey" json:"userId"`
	Username       string `gorm:"column:USERNAME" json:"username"`
	PasswordHash   string `gorm:"column:PASSWORD_HASH" json:"passwordHash"`
	MspID          string `gorm:"column:MSP_ID" json:"mspId"`
	EnrollmentCert []byte `gorm:"column:ENROLLMENT_CERT" json:"enrollmentCert"`
	CreatedAt      string `gorm:"column:CREATED_AT" json:"createdAt"`
	UpdatedAt      string `gorm:"column:UPDATED_AT" json:"updatedAt"`
	Version        int    `gorm:"column:VERSION;default:1" json:"version"`
}

func (User) TableName() string {
	return "t_users"
}

// UpdateUser 更新用户信息
func UpdateUser(db *gorm.DB, userID string, user User) error {
	// 执行更新操作，使用 Where 条件限制更新的行
	if err := db.Model(&User{}).Where("USER_ID = ?", userID).Updates(user).Error; err != nil {
		return err
	}
	return nil
}

// GetUser 获取用户
func GetUser(db *gorm.DB, userID string) (*User, error) {
	var user User
	if err := db.First(&user, "USER_ID = ?", userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
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

	return result.EnrollmentCert, result.MspID, nil
}
