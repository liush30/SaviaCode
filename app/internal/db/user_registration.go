package db

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
	Version        int    `gorm:"column:version"`
}
