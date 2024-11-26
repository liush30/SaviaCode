package fabric

import (
	mspapi "github.com/hyperledger/fabric-sdk-go/pkg/msp/api"
	"log"
)

// Enroll 注册用户
func Enroll(id string) error {
	f := textFixture{}
	f.setup()
	secret, err := f.caClient.Register(&mspapi.RegistrationRequest{Name: id})
	if err != nil {
		log.Fatalf("Registration failed: %v", err)
	}
	log.Println("User successfully registered. Secret:", secret)

	return f.caClient.Enroll(&mspapi.EnrollmentRequest{Name: id, Secret: secret})
}
