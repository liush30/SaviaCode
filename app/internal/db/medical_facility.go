package db

import (
	"fmt"
	"gorm.io/gorm"
)

type MedicalFacility struct {
	TmfID          string `gorm:"primaryKey;column:tmf_id"`
	Name           string `gorm:"column:name"`
	FacilityType   string `gorm:"column:facility_type"`
	LicenseNumber  string `gorm:"column:license_number"`
	ContactPerson  string `gorm:"column:contact_person"`
	PhoneNumber    string `gorm:"column:phone_number"`
	Email          string `gorm:"column:email"`
	Address        string `gorm:"column:address"`
	Status         string `gorm:"column:status"`
	Description    string `gorm:"column:desc"`
	OperatingHours string `gorm:"column:operating_hours"`
	HospitalLevel  string `gorm:"column:hospital_level"`
	CreateAt       string `gorm:"column:create_at"`
	UpdateAt       string `gorm:"column:update_at"`
	Version        int    `gorm:"column:version"`
}

func (MedicalFacility) TableName() string {
	return "t_medical_facility"
}

func CreateMedicalFacility(db *gorm.DB, facility *MedicalFacility) error {
	result := db.Create(facility)

	if result.Error != nil {
		return fmt.Errorf("failed to create facility: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected, facility creation might have failed")
	}

	return nil
}

// UpdateMedicalFacility 更新信息
func UpdateMedicalFacility(db *gorm.DB, facility MedicalFacility) error {
	result := db.Save(&facility)
	if result.Error != nil {
		return fmt.Errorf("failed to update facility: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected, facility update might have failed")
	}
	return nil
}

// QueryMedicalFacility 查询指定类型与状态的医疗机构
func QueryMedicalFacility(db *gorm.DB, facilityType string, status string) ([]MedicalFacility, error) {
	var facilities []MedicalFacility
	result := db.Where("facility_type = ? AND status = ?", facilityType, status).Find(&facilities)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to query facilities: %v", result.Error)
	}
	return facilities, nil
}

type Hospital struct {
	Name  string `gorm:"column:name"`
	TmfID string `gorm:"primaryKey;column:tmf_id"`
}

// QueryMedicalFacilityNameAndID 根据指定类型与状态查询名字与id
func QueryMedicalFacilityNameAndID(db *gorm.DB, facilityType string, status string) ([]Hospital, error) {

	var hospitals []Hospital
	result := db.Select("tmf_id,name").Where("facility_type = ? AND status = ?", facilityType, status).Find(&hospitals)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to query facilities: %v", result.Error)
	}
	return hospitals, nil
}
