package db

import (
	"fmt"
	"gorm.io/gorm"
)

// Dispensing  定义 t_dispensing 结构体
type Dispensing struct {
	TdID           string `gorm:"primary_key;type:varchar(64);not null" json:"td_id"` // 记录ID
	PrescriptionID string `gorm:"type:varchar(10);not null" json:"prescription_id"`   // 处方ID
	PharmacyID     string `gorm:"type:varchar(64);not null" json:"pharmacy_id"`       // 药房ID
	Time           string `gorm:"type:varchar(64)" json:"time"`                       // 预约时间
	Status         string `gorm:"type:varchar(10);not null" json:"status"`            // 状态
}

const (
	statusDispensedPending = "待取药"
	statusDispensedDone    = "已取药"
)

func (d *Dispensing) TableName() string {
	return "t_dispensing"
}

func CreateDispensing(db *gorm.DB, dispensing *Dispensing) error {
	result := db.Create(dispensing)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected, record creation might have failed")
	}
	return nil
}

// QueryDispensingByPharmacyID 根据药房id查询记录
func QueryDispensingByPharmacyID(db *gorm.DB, pharmacyID string, page int, size int) ([]Dispensing, error) {
	var dispensings []Dispensing
	result := db.Where("pharmacy_id = ? and status = ?", pharmacyID, statusDispensedPending).Offset((page - 1) * size).Limit(size).Find(&dispensings)
	if result.Error != nil {
		return nil, result.Error
	}
	return dispensings, nil
}
