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
	userID := "b6a7473b06f64e4c946967b861e9ffac"
	encryptUserID := "a337ceb77f604c13a8307d0af7fa2663"
	msg := "message"
	exp := "(role:pharmacy AND dept:dept1) OR (level:level1)"
	auth := []string{"role", "dept", "level"}
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
