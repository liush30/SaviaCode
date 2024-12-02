package db

import (
	"fmt"
	"gorm.io/gorm"
)

type CryptoExp struct {
	TceID       string `gorm:"primaryKey;type:varchar(64);not null" json:"tce_id"`
	Name        string `gorm:"type:varchar(10);not null" json:"name"`
	Description string `gorm:"type:varchar(10);not null" json:"description"`
	Exp         string `gorm:"type:varchar(64);not null" json:"exp"`
	Auth        string `gorm:"type:varchar(64);nullable" json:"auth"`
	Status      string `gorm:"type:varchar(10);not null" json:"status"`
	CreateAt    string `gorm:"type:datetime;not null" json:"create_at"`
	UpdateAt    string `gorm:"type:datetime;not null" json:"update_at"`
	Version     int    `gorm:"type:int;not null" json:"version"`
}

type CryptoExpResponse struct {
	TceID       string `json:"tce_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (CryptoExp) TableName() string {
	return "t_crypto_exp"
}

func GetCryptoExpAll(db *gorm.DB) ([]CryptoExpResponse, error) {
	var records []CryptoExpResponse
	result := db.Model(&CryptoExp{}).Select("tce_id", "name", "description").Find(&records)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get records: %v", result.Error)
	}
	return records, nil
}

// GetCryptoExpByID 根据id查询指定信息
func GetCryptoExpByID(db *gorm.DB, id string) (*CryptoExp, error) {
	var record CryptoExp
	result := db.Where("tce_id = ?", id).First(&record)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get record: %v", result.Error)
	}
	return &record, nil
}
