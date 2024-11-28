package db

import (
	"fmt"
	"gorm.io/gorm"
)

const (
	activeStatus = "接诊中"
)

type Doctor struct {
	ThID        string `gorm:"primaryKey;column:td_id"`
	HospitalID  string `gorm:"column:hospital_id"`
	Name        string `gorm:"column:name"`
	Gender      string `gorm:"column:gender"`
	Specialty   string `gorm:"column:specialty"`
	Status      string `gorm:"column:status"`
	Level       string `gorm:"column:level"`
	Description string `gorm:"column:desc"`
	CreateAt    string `gorm:"column:create_at"`
	UpdateAt    string `gorm:"column:update_at"`
	Version     int    `gorm:"column:version"`
}

func (d *Doctor) TableName() string {
	return "t_doctor"
}

// CreateDoctor 创建一个doctor
func CreateDoctor(db *gorm.DB, doctor *Doctor) error {
	result := db.Create(doctor)

	if result.Error != nil {
		return fmt.Errorf("failed to create doctor: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected, doctor creation might have failed")
	}

	return nil
}

// UpdateDoctor UpdateDoctor
func UpdateDoctor(db *gorm.DB, doctor *Doctor) error {
	result := db.Save(doctor)

	if result.Error != nil {
		return fmt.Errorf("failed to update doctor: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected, doctor update might have failed")
	}

	return nil
}

// GetDoctor 获取doctor信息
func GetDoctor(db *gorm.DB, thID string) (Doctor, error) {
	var doctor Doctor
	if err := db.Where("td_id = ?", thID).Find(&doctor).Error; err != nil {
		return Doctor{}, err
	}
	return doctor, nil
}

// GetDoctorsByDepartmentID 根据科室ID查询就诊中的所有医生信息
func GetDoctorsByDepartmentID(db *gorm.DB, departmentID string) ([]Doctor, error) {
	var doctors []Doctor
	if err := db.Where("specialty = ? AND status = ?", departmentID, activeStatus).Find(&doctors).Error; err != nil {
		return nil, err
	}
	return doctors, nil
}
