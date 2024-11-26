package db

import (
	"fmt"
	"gorm.io/gorm"
)

type Department struct {
	TdID           string `gorm:"primaryKey;column:td_id"`
	HospitalID     string `gorm:"column:hospital_id"`
	Name           string `gorm:"column:name"`
	Code           string `gorm:"column:code"`
	Category       string `gorm:"column:category"`
	Description    string `gorm:"column:description"`
	DepartmentHead string `gorm:"column:department_head"`
	ContactNumber  string `gorm:"column:contact_number"`
	Location       string `gorm:"column:location"`
	Specialties    string `gorm:"column:specialties"`
	WorkingHours   string `gorm:"column:working_hours"`
	Status         string `gorm:"column:status"`
	CreateAt       string `gorm:"column:create_at"`
	UpdateAt       string `gorm:"column:update_at"`
	Version        int    `gorm:"column:version"`
}

func (Department) TableName() string {
	return "t_departments"
}

func CreateDepartment(db *gorm.DB, department Department) error {
	// 执行创建操作
	result := db.Create(&department)

	// 检查是否有错误发生
	if result.Error != nil {
		// 返回原始的错误信息
		return fmt.Errorf("failed to create department: %v", result.Error)
	}

	// 检查是否新增了记录
	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected, department creation might have failed")
	}

	// 成功时返回 nil
	return nil
}

// GetDepartmentsByHospitalID 查询医院的所有科室
func GetDepartmentsByHospitalID(db *gorm.DB, hospitalID string) ([]Department, error) {
	var departments []Department
	if err := db.Where("hospital_id = ?", hospitalID).Find(&departments).Error; err != nil {
		return nil, err
	}
	return departments, nil
}

// UpdateDepartment 修改科室信息
func UpdateDepartment(db *gorm.DB, department Department) error {
	result := db.Save(&department)
	if result.Error != nil {
		return fmt.Errorf("failed to update department: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected, department update might have failed")
	}
	return nil
}
