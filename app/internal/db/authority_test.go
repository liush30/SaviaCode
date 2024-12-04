package db

import "testing"

func TestGetAuthoritiesNameAndAttributes(t *testing.T) {
	dbClient, err := InitDB()
	if err != nil {
		t.Fatal(err)
	}

	authorities, err := GetAuthoritiesNameAndAttributes(dbClient)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(authorities)
}
