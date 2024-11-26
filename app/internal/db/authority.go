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
