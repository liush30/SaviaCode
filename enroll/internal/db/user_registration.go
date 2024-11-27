package db

import (
	"fmt"
	"gorm.io/gorm"
)

type UserRegistration struct {
	TurID          string `gorm:"primaryKey;column:tur_id"`
	Name           string `gorm:"column:name"`
	IDNumber       string `gorm:"column:id_number"`
	Gender         string `gorm:"column:gender"`
	DateOfBirth    string `gorm:"column:date_of_birth"`
	BloodType      string `gorm:"column:blood_type"`
	HeightCm       string `gorm:"column:height_cm"`
	WeightKg       string `gorm:"column:weight_kg"`
	Address        string `gorm:"column:address"`
	PhoneNumber    string `gorm:"column:phone_number"`
	AllergyHistory string `gorm:"column:allergy_history"`
	MedicalHistory string `gorm:"column:medical_history"`
	CreateAt       string `gorm:"column:create_at"`
	UpdateAt       string `gorm:"column:update_at"`
	Result         string `gorm:"column:result"`
	Version        int    `gorm:"column:version"`
	Password       string `gorm:"column:password"`
}

func (UserRegistration) TableName() string {
	return "t_user_registration"
}
func CreateUserRegistration(db *gorm.DB, registration *UserRegistration) error {
	result := db.Create(registration)
	if result.Error != nil {
		return fmt.Errorf("failed to create registration: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected, registration creation might have failed")
	}
	return nil

}

func GetAllRegistrationWithConditions(db *gorm.DB, conditions map[string]interface{}, page, pageSize int) ([]UserRegistration, error) {
	var req []UserRegistration

	// 计算 OFFSET
	offset := (page - 1) * pageSize

	// 使用 Limit 和 Offset 进行分页查询，并根据条件构建 WHERE 子句
	err := db.Where(conditions).Limit(pageSize).Offset(offset).Find(&req).Error
	return req, err
}

func UpdateRegistration(db *gorm.DB, id string, updatedFields interface{}) error {
	return db.Model(&UserRegistration{}).Where("tur_id = ?", id).Updates(updatedFields).Error
}

// GetRegistration 查询指定id的数据
func GetRegistration(db *gorm.DB, id string) (*UserRegistration, error) {
	var req UserRegistration
	err := db.Where("tur_id = ?", id).First(&req).Error
	return &req, err
}