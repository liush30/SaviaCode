package fabric

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/cryptosuite"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/cryptosuite/bccsp/pkcs11"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/mocks"
	fabImpl "github.com/hyperledger/fabric-sdk-go/pkg/fab"
	msp2 "github.com/hyperledger/fabric-sdk-go/pkg/msp"
	mspapi "github.com/hyperledger/fabric-sdk-go/pkg/msp/api"
	"github.com/hyperledger/fabric-sdk-go/test/metadata"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type textFixture struct {
	endpointConfig fab.EndpointConfig //存储和管理网络配置（如对等节点、排序服务、CA 服务器等）的配置信息。提供对网络组件的配置，例如节点的地址和端口，网络的拓扑结构等。这些信息用于与 Fabric 网络进行交互。
	//存储和管理身份相关的配置信息。这些身份参数包括用户的凭证存储路径、证书、私钥和其他与组织成员身份认证相关的配置。
	//它的主要作用是在与网络交互时，提供用户身份信息的支持。
	identityConfig msp.IdentityConfig
	//配置加密套件的设置，包含有关加密操作的配置，如密钥存储路径、加密算法等，这些设置用于初始化和配置加密服务
	cryptSuiteConfig core.CryptoSuiteConfig
	//提供加密和解密服务，用于执行加密操作，如生成密钥对、签名、验证等，它是实际的加密操作的实现工具，依据 cryptSuiteConfig 进行配置。
	cryptoSuite core.CryptoSuite
	// 存储用户证书和相关信息。管理和存储用户的证书和私钥。它用于在系统中管理用户身份和凭证
	userStore msp.UserStore
	//与 CA 服务器进行交互，处理证书的注册和管理。提供与证书颁发机构（CA）进行交互的客户端功能，比如注册用户、颁发证书等。它负责在 Fabric 网络中处理证书的生命周期管理。
	caClient mspapi.CAClient
	//提供身份管理器的功能。根据组织名称返回相应的身份管理器。身份管理器用于处理组织的身份管理，包括用户身份的创建、验证和管理。
	identityManagerProvider msp.IdentityManagerProvider
}

func (f *textFixture) setup() {
	configPath := filepath.Join(metadata.GetProjectPath(), configTestFile)
	backend, err := getCustomBackend(configPath)
	if err != nil {
		panic(err)
	}
	//从给定的 core.ConfigBackend 后端中创建并返回 msp.IdentityConfig
	f.identityConfig, err = msp2.ConfigFromBackend(backend...)
	if err != nil {
		panic(err)
	}
	//从后端配置中提取Crypto Suite相关的配置信息
	f.cryptSuiteConfig = cryptosuite.ConfigFromBackend(backend...)

	f.endpointConfig, err = fabImpl.ConfigFromBackend(backend...)
	if err != nil {
		panic(fmt.Sprintf("Failed to read config : %s", err))
	}

	cleanup(f.identityConfig.CredentialStorePath())
	log.Println(f.cryptSuiteConfig.SecurityProvider())
	f.cryptoSuite, err = pkcs11.GetSuiteByConfig(f.cryptSuiteConfig)
	if err != nil {
		panic(err)
	}
	//检查从identityConfig中获取的凭证存储路径是否为空，如果路径不为空，说明配置文件定义了用于存储用户凭证的位置
	if f.identityConfig.CredentialStorePath() != "" {
		f.userStore, err = msp2.NewCertFileUserStore(f.identityConfig.CredentialStorePath())
		if err != nil {
			panic(fmt.Sprintf("creating a user store failed: %s", err))
		}
	}
	f.userStore = userStoreFromConfig()

	identityManagers := make(map[string]msp.IdentityManager)
	netConfig := f.endpointConfig.NetworkConfig()
	if netConfig == nil {
		panic("failed to get network config")
	}
	for orgName := range netConfig.Organizations {
		mgr, err1 := msp2.NewIdentityManager(orgName, f.userStore, f.cryptoSuite, f.endpointConfig)
		if err1 != nil {
			panic(fmt.Sprintf("failed to create identity manager for org %s: %s", orgName, err1))

		}
		identityManagers[orgName] = mgr
	}

	f.identityManagerProvider = &identityManagerProvider{identityManager: identityManagers}

	ctxProvider := context.NewProvider(context.WithIdentityManagerProvider(f.identityManagerProvider),
		context.WithUserStore(f.userStore), context.WithCryptoSuite(f.cryptoSuite),
		context.WithCryptoSuiteConfig(f.cryptSuiteConfig),
		context.WithEndpointConfig(f.endpointConfig),
		context.WithIdentityConfig(f.identityConfig))

	ctx := &context.Client{Providers: ctxProvider}

	if err != nil {
		panic(fmt.Sprintf("failed to created context for test setup: %s", err))
	}

	f.caClient, err = msp2.NewCAClient(org1, ctx)
	if err != nil {
		panic(fmt.Sprintf("NewCAClient returned error: %s", err))
	}

}

const (
	configTestFile = "config/config_e2e_pkcs11.yaml"
	org1           = "org1"
)

func userStoreFromConfig() msp.UserStore {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"lsh", "lsh666hh", "zhoupb.com", 33060, "eldercare_health")
	userStore, err := NewMySQLUserStore(dsn)
	if err != nil {
		panic(fmt.Sprintf("creating a user store failed: %s", err))
	}
	return userStore
}

func cleanup(storePath string) {
	err := os.RemoveAll(storePath)
	if err != nil {
		panic(fmt.Sprintf("Failed to remove dir %s: %s\n", storePath, err))
	}
}

func getCustomBackend(configPath string) ([]core.ConfigBackend, error) {
	// 从配置文件中加载配置
	configBackends, err := config.FromFile(configPath)()
	if err != nil {
		return nil, err
	}

	// 创建一个映射来存储 CA 客户端相关的配置
	backendMap := make(map[string]interface{})
	backendMap["client"], _ = configBackends[0].Lookup("client")
	backendMap["organizations"], _ = configBackends[0].Lookup("organizations")
	// 从配置中提取 CA 客户端相关的配置
	backendMap["certificateAuthorities"], _ = configBackends[0].Lookup("certificateAuthorities")
	// 创建并返回一个只包含 CA 客户端和 BCCSP 配置的 MockConfigBackend
	backends := append([]core.ConfigBackend{}, &mocks.MockConfigBackend{KeyValueMap: backendMap})
	backends = append(backends, configBackends...)
	return backends, nil
}

type identityManagerProvider struct {
	identityManager map[string]msp.IdentityManager
}

// IdentityManager returns the organization's identity manager
func (p *identityManagerProvider) IdentityManager(orgName string) (msp.IdentityManager, bool) {
	im, ok := p.identityManager[strings.ToLower(orgName)]
	if !ok {
		return nil, false
	}
	return im, true
}
