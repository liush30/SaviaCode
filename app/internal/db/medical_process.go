package db

import (
	"fmt"
	"gorm.io/gorm"
)

type MedicalProcess struct {
	ProcessID   string `gorm:"primaryKey;column:tmp_id"` // 记录ID
	VisitID     string `gorm:"column:tmr_id"`            // 就诊记录ID
	RecordType  string `gorm:"column:record_type"`       // 就诊类型
	RecordValue string `gorm:"column:record_value"`      // 记录内容
	Status      string `gorm:"column:status"`            // 记录状态
	Sign        string `gorm:"column:sign"`              // 签名
	CreateAt    string `gorm:"column:create_at"`         // 创建时间
	UpdateAt    string `gorm:"column:update_at"`         // 更新时间
	//Version     int    `gorm:"column:version"`           // 版本
	//PharmacyID string `gorm:"column:pharmacy"` // 药房ID
}

const (
	ProcessStatusFinished = "已上链"
)

func (MedicalProcess) TableName() string {
	return "t_medical_process"
}

func CreateMedicalProcess(db *gorm.DB, record *MedicalProcess) error {
	result := db.Create(record)
	if result.Error != nil {
		return fmt.Errorf("failed to create record: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected, record creation might have failed")
	}
	return nil
}

// GetMedicalProcess 查看指定就诊记录的id
func GetMedicalProcess(db *gorm.DB, visitID string, page, size int) ([]MedicalProcess, error) {
	var records []MedicalProcess
	result := db.Limit(size).Offset((page-1)*size).Where("tmr_id = ?", visitID).Find(&records)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get records: %v", result.Error)
	}
	return records, nil
}

// UpdateMedicalProcess 更新就诊记录
func UpdateMedicalProcess(db *gorm.DB, record *MedicalProcess) error {
	result := db.Save(record)
	if result.Error != nil {
		return fmt.Errorf("failed to update record: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected, record update might have failed")
	}
	return nil
}

// GetMedicalProcessByID 根据id查询指定信息
func GetMedicalProcessByID(db *gorm.DB, id string) (*MedicalProcess, error) {
	var record MedicalProcess
	result := db.Where("tmp_id = ?", id).First(&record)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get record: %v", result.Error)
	}
	return &record, nil
}

// DeleteMedicalProcess 删除就诊过程记录
func DeleteMedicalProcess(db *gorm.DB, processId string) error {
	info, err := GetMedicalProcessByID(db, processId)
	if err != nil {
		return err
	}
	if (*info).Status == ProcessStatusFinished {
		return fmt.Errorf("record has been finished")
	}
	result := db.Delete(&MedicalProcess{}, "tmp_id = ?", processId)
	if result.Error != nil {
		return fmt.Errorf("failed to delete records: %v", result.Error)
	}
	return nil
}
