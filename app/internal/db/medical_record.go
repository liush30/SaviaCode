package db

import (
	"fmt"
	"gorm.io/gorm"
)

type MedicalRecord struct {
	TmrID     string `gorm:"primaryKey;column:tmr_id"`
	PatientID string `gorm:"column:patient_id"`
	//DepartmentID string `gorm:"column:department_id"`
	DoctorID  string `gorm:"column:doctor_id"`
	VisitDate string `gorm:"column:visit_date"`
	Status    string `gorm:"column:status"`
	CreateAt  string `gorm:"column:create_at"`
	UpdateAt  string `gorm:"column:update_at"`
	Version   int    `gorm:"column:version"`
}

func (MedicalRecord) TableName() string {
	return "t_medical_record"
}

// CreateMedicalRecord 新增record
func CreateMedicalRecord(db *gorm.DB, record *MedicalRecord) error {
	result := db.Create(record)
	if result.Error != nil {
		return fmt.Errorf("failed to create record: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected, record creation might have failed")
	}
	return nil
}

// GetMedicalRecord 查询record
func GetMedicalRecord(db *gorm.DB, tdID string) (*MedicalRecord, error) {
	var record MedicalRecord
	result := db.Where("tmr_id = ?", tdID).First(&record)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get record: %v", result.Error)
	}
	return &record, nil
}

// QueryMedicalRecordByConditions 查询指定用户的所有就诊记录
func QueryMedicalRecordByConditions(db *gorm.DB, conditions map[string]interface{}, page, size int) ([]MedicalRecord, error) {
	var records []MedicalRecord
	result := db.Where(conditions).Offset((page - 1) * size).Limit(size).Find(&records)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to query records: %v", result.Error)
	}
	return records, nil
}

// QueryMedicalRecordByConditionsNoOffset 查询指定用户的所有就诊记录
func QueryMedicalRecordByConditionsNoOffset(db *gorm.DB, conditions map[string]interface{}) ([]MedicalRecord, error) {
	var records []MedicalRecord
	result := db.Where(conditions).Find(&records)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to query records: %v", result.Error)
	}
	return records, nil
}

// QueryMedicalRecordByDoctorID 查询指定医生的所有就诊记录
func QueryMedicalRecordByDoctorID(db *gorm.DB, doctorID, status string, page, size int) ([]MedicalRecord, error) {
	var records []MedicalRecord
	result := db.Where("doctor_id = ? AND status = ?", doctorID, status).Offset((page - 1) * size).Limit(size).Find(&records)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to query records: %v", result.Error)
	}
	return records, nil
}

func UpdateMedicalRecordStatus(db *gorm.DB, id, status string) error {
	result := db.Model(&MedicalRecord{}).Where("tmr_id = ?", id).Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("failed to update record: %v", result.Error)
	}
	return nil
}
