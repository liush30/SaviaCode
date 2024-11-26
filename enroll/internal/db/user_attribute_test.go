package db

import (
	"lyods-fabric-demo/enroll/internal/pkg/tool"
	"testing"
)

func TestGetUserAttributesByUserID(t *testing.T) {
	dbClient, err := InitDB()
	if err != nil {
		t.Fatal(err)
	}

	attributes, err := GetUserAttributesByUserID(dbClient, "user1")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(attributes)
}

func TestBatchInsertUserAttributes(t *testing.T) {
	dbClient, err := InitDB()
	if err != nil {
		t.Fatal(err)
	}
	var userAttributes []UserAttribute

	userAttributes = append(userAttributes, UserAttribute{
		TuaKey:    tool.GenerateUUIDWithoutDashes(),
		UserID:    "user1",
		TaKey:     "taKey2",
		Attribute: "attribute2",
	})
	userAttributes = append(userAttributes, UserAttribute{
		TuaKey:    tool.GenerateUUIDWithoutDashes(),
		UserID:    "user1",
		TaKey:     "taKey2",
		Attribute: "attribute2",
	})
	userAttributes = append(userAttributes, UserAttribute{
		TuaKey:    tool.GenerateUUIDWithoutDashes(),
		UserID:    "user1",
		TaKey:     "taKey3",
		Attribute: "attribute3",
	})
	err = BatchInsertUserAttributes(dbClient, userAttributes)
	if err != nil {
		t.Fatal(err)
	}

}
