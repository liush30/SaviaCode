package db

import (
	"gorm.io/gorm"
)

type AuthorityAttributes struct {
	AuthorityName string `gorm:"column:AUTH_NAME"`
	Attributes    string `gorm:"column:ATTRIBUTE"`
}

// Authority 定义授权
type Authority struct {
	TaKey                string `gorm:"column:TA_KEY;primaryKey" json:"taKey"`
	AuthorityName        string `gorm:"column:AUTH_NAME" json:"authorityName"`
	AuthorityDescription string `gorm:"column:AUTH_DESCRIPTION" json:"authorityDescription"`
	Attributes           string `gorm:"column:ATTRIBUTES" json:"attributes"`
	CreatedAt            string `gorm:"column:CREATED_AT" json:"createdAt"`
	UpdatedAt            string `gorm:"column:UPDATED_AT" json:"updatedAt"`
	Version              int    `gorm:"column:VERSION;default:1" json:"version"`
}

func (Authority) TableName() string {
	return "t_authority"
}

// CreateAttribute 创建属性
func CreateAttribute(db *gorm.DB, data Authority) error {
	return db.Create(&data).Error
}

// UpdateAttribute 更新属性
func UpdateAttribute(db *gorm.DB, data Authority) error {
	return db.Save(&data).Error
}

// GetAttribute 获取属性
func GetAttribute(db *gorm.DB, attributeKey string) (*Authority, error) {
	var data Authority
	if err := db.First(&data, "TA_KEY = ?", attributeKey).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAuthoritiesNameAndAttributes 获取所有的授权
func GetAuthoritiesNameAndAttributes(db *gorm.DB) ([]AuthorityAttributes, error) {
	var authorities []AuthorityAttributes

	// 执行查询
	if err := db.Table("t_authority").
		Select("AUTH_NAME, ATTRIBUTE").
		Find(&authorities).Error; err != nil {
		return nil, err
	}

	return authorities, nil
}
