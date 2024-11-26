package db

import "testing"

func TestGetUserCertAndMspID(t *testing.T) {
	dbClient, err := InitDB()
	if err != nil {
		t.Fatalf("failed to init internal: %v", err)
	}

	cert, msp, err := GetUserCertAndMspID(dbClient, "b6a7473b06f64e4c946967b861e9ffac")
	if err != nil {
		t.Fatalf("failed to get user attributes: %v", err)
	}

	t.Logf("cert: %v", string(cert))
	t.Logf("msp: %v", msp)

}

func TestGetUserAttributesByUserID(t *testing.T) {

	dbClient, err := InitDB()
	if err != nil {
		t.Fatalf("failed to init internal: %v", err)
	}

	userAttributes, err := GetUserAttributesByUserID(dbClient, "b6a7473b06f64e4c946967b861e9ffac")
	if err != nil {
		t.Fatalf("failed to get user attributes: %v", err)
	}

	t.Logf("user attributes: %v", userAttributes)
}
