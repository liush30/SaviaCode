package db

import "gorm.io/gorm"

type EncryptedData struct {
	TEDKey         string `gorm:"column:TED_KEY;primaryKey"`
	UserID         string `gorm:"column:USER_ID"`
	C0             []byte `gorm:"column:C0"`
	C1X            []byte `gorm:"column:C1X"`
	C2X            []byte `gorm:"column:C2X"`
	C3X            []byte `gorm:"column:C3X"`
	MSP            []byte `gorm:"column:MSP"`
	CreateDate     string `gorm:"column:CREATE_DATE"`
	LastModifyDate string `gorm:"column:LAST_MODIFY_DATE"`
	Version        int    `gorm:"column:VERSION;default:1"`
}

func (EncryptedData) TableName() string {
	return "t_encrypted_data"
}

// CreateEncryptedData 创建加密数据
func CreateEncryptedData(db *gorm.DB, data EncryptedData) error {
	return db.Create(&data).Error
}

// GetEncryptedData 获取加密数据
func GetEncryptedData(db *gorm.DB, tedKey string) (*EncryptedData, error) {
	var data EncryptedData
	if err := db.First(&data, "TED_KEY = ?", tedKey).Error; err != nil {
		return nil, err
	}
	return &data, nil
}
