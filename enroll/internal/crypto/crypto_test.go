package crypto

import (
	"encoding/json"
	"testing"
)

func TestGeneratePublicKey(t *testing.T) {
	auth := []string{"role", "dept"}
	pk := generatePublicKey(auth)

	t.Logf("public keys: %v", pk)
}
func TestEncrypt(t *testing.T) {
	userID := "user3"
	encryptUserID := "user2"
	msg := "message"
	exp := "(role:physician AND dept:dept1) OR (role:pharmacy AND dept:dept2)"
	auth := []string{"role", "dept"}
	encrypt, err := Encrypt(userID, msg, exp, auth)
	if err != nil {
		t.Fatalf("failed to encrypt: %v", err)
	}

	t.Logf("encrypted message: %s", encrypt)
	var chainData OnChainData
	if err := json.Unmarshal(encrypt, &chainData); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	decrypt, err := Decrypt(encrypt, encryptUserID)
	if err != nil {
		t.Fatalf("failed to decrypt: %v", err)
	}

	t.Logf("decrypted message: %s", decrypt)
}
