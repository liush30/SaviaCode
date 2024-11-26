package tool

import (
	"testing"
)

func TestGenerateUUIDWithoutDashes(t *testing.T) {

	uuid := GenerateUUIDWithoutDashes()

	t.Log(uuid)
}
