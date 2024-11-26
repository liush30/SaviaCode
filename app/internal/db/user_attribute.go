package db

import (
	"gorm.io/gorm"
)

type UserAttribute struct {
	TuaKey    string `gorm:"column:TUA_KEY;primaryKey" json:"id"`
	UserID    string `gorm:"column:USER_ID" json:"userId"`
	TaKey     string `gorm:"column:TA_KEY" json:"taKey"`
	Attribute string `gorm:"column:ATTRIBUTE" json:"attribute"`
	CreatedAt string `gorm:"column:CREATED_AT" json:"createdAt"`
	UpdatedAt string `gorm:"column:UPDATED_AT" json:"updatedAt"`
	Version   int    `gorm:"column:VERSION;default:1" json:"version"`
}

func (UserAttribute) TableName() string {
	return "t_user_attributes"
}

// GetUserAttributesByUserID 获取用户属性
func GetUserAttributesByUserID(db *gorm.DB, userID string) ([]AuthorityAttributes, error) {
	var results []AuthorityAttributes

	// 执行查询
	if err := db.Table("T_USER_ATTRIBUTES U").
		Select("A.AUTH_NAME, GROUP_CONCAT(U.ATTRIBUTE) AS ATTRIBUTE").
		Joins("JOIN T_AUTHORITY A ON A.TA_KEY = U.TA_KEY").
		Where("U.USER_ID = ?", userID).
		Group("A.AUTH_NAME").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}
