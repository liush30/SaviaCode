package crypto

import "testing"

func TestGenerateUserKey(t *testing.T) {
	//key, err := generateUserKey("user1")
	//if err != nil {
	//	return
	//}
	//t.Logf("user key: %v", key)

	key, err := generateUserKey("b6a7473b06f64e4c946967b861e9ffac")
	if err != nil {
		t.Fatalf("failed to generate user key: %v", err)
		return
	}
	t.Logf("user key: %v", key)
}
