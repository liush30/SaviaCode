package db

import "testing"

func TestQueryMedicalFacilityNameAndID(t *testing.T) {
	//新增数据
	//info1 := MedicalFacility{
	//	TmfID:          "1",
	//	Name:           "1",
	//	FacilityType:   "医院",
	//	LicenseNumber:  "1",
	//	ContactPerson:  "1",
	//	PhoneNumber:    "1",
	//	Email:          "1",
	//	Address:        "1",
	//	Status:         "1",
	//	Description:    "1",
	//	OperatingHours: "1",
	//	HospitalLevel:  "1",
	//	CreateAt:       "1",
	//	UpdateAt:       "1",
	//	Version:        1,
	//}
	//
	//info2 := MedicalFacility{
	//	TmfID:          "2",
	//	Name:           "2",
	//	FacilityType:   "医院",
	//	LicenseNumber:  "2",
	//	ContactPerson:  "2",
	//	PhoneNumber:    "2",
	//	Email:          "2",
	//	Address:        "2",
	//	Status:         "2",
	//	Description:    "2",
	//	OperatingHours: "2",
	//	HospitalLevel:  "2",
	//	CreateAt:       "2",
	//	UpdateAt:       "2",
	//	Version:        2,
	//}
	//inf3 := MedicalFacility{
	//	TmfID:          "3",
	//	Name:           "3",
	//	FacilityType:   "药房",
	//	LicenseNumber:  "3",
	//	ContactPerson:  "3",
	//	PhoneNumber:    "3",
	//	Email:          "3",
	//	Address:        "3",
	//	Status:         "3",
	//	Description:    "3",
	//	OperatingHours: "3",
	//	HospitalLevel:  "3",
	//	CreateAt:       "3",
	//	UpdateAt:       "3",
	//	Version:        3,
	//}
	//
	dbClient, err := InitDB()
	if err != nil {
		t.Fatal(err)
	}
	//err = CreateMedicalFacility(dbClient, &info1)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//err = CreateMedicalFacility(dbClient, &info2)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//err = CreateMedicalFacility(dbClient, &inf3)
	//if err != nil {
	//	t.Fatal(err)
	//}

	//查询数据
	medicalFacilities, err := QueryMedicalFacilityNameAndID(dbClient, "医院", "1")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(medicalFacilities)

}
