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

// CreateUserAttribute 创建用户属性
func CreateUserAttribute(db *gorm.DB, data UserAttribute) error {
	return db.Create(&data).Error
}

// UpdateUserAttribute 更新用户属性
func UpdateUserAttribute(db *gorm.DB, data UserAttribute) error {
	return db.Save(&data).Error
}

// GetUserAttribute 获取用户属性
func GetUserAttribute(db *gorm.DB, key string) (*UserAttribute, error) {
	var data UserAttribute
	if err := db.First(&data, "TUA_KEY = ?", key).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetUserAttributeByCondition 根据条件查询用户属性
func GetUserAttributeByCondition(db *gorm.DB, condition map[string]interface{}) (*UserAttribute, error) {
	var data UserAttribute
	if err := db.Where(condition).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// BatchInsertUserAttributes 批量插入用户属性
func BatchInsertUserAttributes(db *gorm.DB, attributes []UserAttribute) error {
	return db.Create(&attributes).Error
}

// DeleteUserAttribute 删除用户属性
func DeleteUserAttribute(db *gorm.DB, key string) error {
	return db.Delete(&UserAttribute{}, "TUA_KEY = ?", key).Error
}

// GetUserAttributes 获取用户属性
func GetUserAttributes(db *gorm.DB, userID string) ([]UserAttribute, error) {
	var data []UserAttribute
	if err := db.Find(&data, "USER_ID = ?", userID).Error; err != nil {
		return nil, err
	}
	return data, nil
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
