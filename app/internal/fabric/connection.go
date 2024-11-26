//go:build pkcs11
// +build pkcs11

package fabric

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/hash"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"os"
)

const (
	org1MSPID = "Org1MSP"
	org2MSPID = "Org2MSP"
	// org1 tls证书路径
	org1TlsCertPath = "/home/savia/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"
	// org1 peer endpoint
	org1PeerEndpoint = "dns:///localhost:7051"
	// org1 server name
	org1ServerName = "peer0.org1.example.com"
	// org2 tls证书路径
	org2TlsCertPath = "/home/savia/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"
	// org2 peer endpoint
	org2PeerEndpoint = "dns:///localhost:9051"
	// org2 server name
	org2ServerName = "peer0.org2.example.com"
	//softhsm key label
	label = "fabric-hsm"
	//softhsm pin
	pin = "987654"
)
const (
	channelName             = "mychannel"
	dispensingChaincodeName = "dispensing"
	personalChaincodeName   = "personal"
	medicalChaincodeName    = "medical"
)

var Org1GrpcConnection *grpc.ClientConn
var Org2GrpcConnection *grpc.ClientConn

func init() {
	Org1GrpcConnection = newGrpcConnection(org1TlsCertPath, org1PeerEndpoint, org1ServerName)
	Org2GrpcConnection = newGrpcConnection(org2TlsCertPath, org2PeerEndpoint, org2ServerName)
}

// newGrpcConnection 建立gRPC连接
func newGrpcConnection(tlsCertPath, peerEndpoint, org1ServerName string) *grpc.ClientConn {
	certificate, err := loadCertificate(tlsCertPath)
	if err != nil {
		panic(fmt.Errorf("failed to obtain commit status: %w", err))
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, org1ServerName)

	connection, err := grpc.NewClient(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}

	return connection
}

// CreateHSMSign 创建hsm签名
func CreateHSMSign(hsmSignerFactory *identity.HSMSignerFactory, cert []byte) (identity.Sign, identity.HSMSignClose, error) {
	return newHSMSign(hsmSignerFactory, getSKI(cert))
}

// CreateHSMSignerFactory 创建hsm签名工厂
func CreateHSMSignerFactory() (*identity.HSMSignerFactory, error) {
	hsmSignerFactory, err := identity.NewHSMSignerFactory(findSoftHSMLibrary())
	if err != nil {
		return nil, fmt.Errorf("failed to new hsm signer: %w", err)
	}
	return hsmSignerFactory, nil
}
func SetupGateway(mspID string, cert []byte, hsmSign identity.Sign) (*client.Gateway, error) {
	grpcConnection, err := getGrpcConnection(mspID)
	if err != nil {
		return nil, fmt.Errorf("failed to get grpc connection: %w", err)
	}
	id := newIdentity(cert, mspID)
	gateway, err := client.Connect(id, client.WithSign(hsmSign), client.WithHash(hash.SHA256), client.WithClientConnection(grpcConnection))
	if err != nil {
		return nil, fmt.Errorf("failed to gateway connect: %w", err)
	}
	return gateway, nil
}
func getGrpcConnection(mspID string) (*grpc.ClientConn, error) {
	if mspID == org1MSPID {
		return Org1GrpcConnection, nil
	} else if mspID == org2MSPID {
		return Org2GrpcConnection, nil
	} else {
		return nil, fmt.Errorf("mspID not found,mspID:%s", mspID)
	}
}
func getContract(gateway *client.Gateway, channelName, chaincodeName string) *client.Contract {
	network := gateway.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)
	return contract
}

// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func newIdentity(certificatePEM []byte, mspID string) *identity.X509Identity {
	cert, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}
	id, err := identity.NewX509Identity(mspID, cert)
	if err != nil {
		panic(err)
	}

	return id
}

func newHSMSign(h *identity.HSMSignerFactory, certPEM []byte) (identity.Sign, identity.HSMSignClose, error) {
	opt := identity.HSMSignerOptions{
		Label:      label,
		Pin:        pin,
		Identifier: string(certPEM),
	}

	sign, close, err := h.NewHSMSigner(opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to new hsm signer: %w", err)
	}

	return sign, close, nil
}

// loadCertificate 加载证书
func loadCertificate(filename string) (*x509.Certificate, error) {
	certificatePEM, err := os.ReadFile(filename) //#nosec G304
	if err != nil {
		return nil, err
	}

	return identity.CertificateFromPEM(certificatePEM)
}

// findSoftHSMLibrary 查找SoftHSM
func findSoftHSMLibrary() string {

	libraryLocations := []string{
		"/usr/lib/softhsm/libsofthsm2.so",
		"/usr/lib/x86_64-linux-gnu/softhsm/libsofthsm2.so",
		"/usr/local/lib/softhsm/libsofthsm2.so",
		"/usr/lib/libacsp-pkcs11.so",
		"/opt/homebrew/lib/softhsm/libsofthsm2.so",
	}
	pkcs11lib := os.Getenv("PKCS11_LIB")
	if pkcs11lib != "" {
		libraryLocations = append(libraryLocations, pkcs11lib)
	}
	for _, libraryLocation := range libraryLocations {
		if _, err := os.Stat(libraryLocation); !errors.Is(err, os.ErrNotExist) {
			return libraryLocation
		}
	}

	panic("No SoftHSM library can be found. The Sample requires SoftHSM to be installed")
}
func getSKI(certPEM []byte) []byte {
	block, _ := pem.Decode(certPEM)

	x590cert, _ := x509.ParseCertificate(block.Bytes)
	pk := x590cert.PublicKey

	return skiForKey(pk.(*ecdsa.PublicKey))
}

func skiForKey(pk *ecdsa.PublicKey) []byte {
	ski := sha256.Sum256(elliptic.Marshal(pk.Curve, pk.X, pk.Y))
	return ski[:]
}
