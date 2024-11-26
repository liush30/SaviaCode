package db

import (
	"fmt"
	"gorm.io/gorm"
)

type Log struct {
	TlID            string `gorm:"primaryKey;column:tl_id"`
	UserID          string `gorm:"column:user_id"`
	OperationType   string `gorm:"column:operation_type"`
	OperationDetail string `gorm:"column:operation_detail"`
	Status          string `gorm:"column:status"`
	ErrorMessage    string `gorm:"column:error_message"`
	CreateAt        string `gorm:"column:create_at"`
	UpdateAt        string `gorm:"column:update_at"`
	Version         int    `gorm:"column:version"`
}

func (Log) TableName() string {
	return "t_log"
}

func CreateLog(db *gorm.DB, log *Log) error {
	result := db.Create(log)

	if result.Error != nil {
		return fmt.Errorf("failed to create log: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected, log creation might have failed")
	}

	return nil

}

// QueryLogsByUserID 查询指定用户的日志
func QueryLogsByUserID(db *gorm.DB, userID string) ([]Log, error) {
	var logs []Log
	result := db.Where("user_id = ?", userID).Find(&logs)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to query logs: %v", result.Error)
	}
	return logs, nil
}
