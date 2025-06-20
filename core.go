package core

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func GenerateKeyPair() (privateKey, publicKey string, err error) {
	private, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return "", "", err
	}
	return private.String(), private.PublicKey().String(), nil
}

func ValidateKey(key string) error {
	_, err := wgtypes.ParseKey(key)
	if err != nil {
		return err
	}
	return nil
}

type WGConfigGenerator struct {
	serverConfig     ServerConfig
	generationConfig GenerationConfig
}

func NewGenerator(serverConfig ServerConfig, generationConfig GenerationConfig) (*WGConfigGenerator, error) {
	err := ValidateKey(serverConfig.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("ivalid server public key: %v", err)
	}

	if err := os.MkdirAll(generationConfig.Dir, 0755); err != nil {
		return nil, fmt.Errorf("failed ro create output dir: %v", err)
	}

	return &WGConfigGenerator{
		serverConfig:     serverConfig,
		generationConfig: generationConfig,
	}, nil
}

func (g *WGConfigGenerator) Generate(uniqueNum int, uniqueCode string, clientConfig ClientConfig) (path string, err error) {
	var publicKey string
	privateKey := clientConfig.PrivateKey
	if privateKey == "" {
		privKey, pubkey, err := GenerateKeyPair()
		if err != nil {
			return "", fmt.Errorf("key generation failed: %v", err)
		}
		privateKey = privKey
		clientConfig.PrivateKey = privKey
		publicKey = pubkey
	} else {
		privkey, err := wgtypes.ParseKey(privateKey)
		if err != nil {
			return "", fmt.Errorf("ivalid private key: %v", err)
		}
		publicKey = privkey.PublicKey().String()
	}

	if clientConfig.ClientIP == "" {
		clientConfig.ClientIP = fmt.Sprintf("10.8.0.%d", (uniqueNum%253)+2)
	}
	if clientConfig.PresistentKeepAlive == 0 {
		clientConfig.PresistentKeepAlive = 25
	}
	wgConfig := fmt.Sprintf("[Interface]\nPrivateKey = %s\nAddress = %s/24\nDNS = %s\n[Peer]\nPublicKey = %s\nEndpoint = %s:%d\nAllowedIPs = %s\nPersistentKeepalive = %d",
		privateKey,
		clientConfig.ClientIP,
		g.serverConfig.DNS,
		g.serverConfig.PublicKey,
		g.serverConfig.PublicIP,
		g.serverConfig.Port,
		g.serverConfig.AllowedIPs,
		clientConfig.PresistentKeepAlive,
	)

	wgConfigName := fmt.Sprintf("wg_%s.conf", uniqueCode)
	wgConfigPath := filepath.Join(g.generationConfig.Dir, wgConfigName)

	if err := os.WriteFile(wgConfigPath, []byte(wgConfig), 0600); err != nil {
		return "", fmt.Errorf("WriteFile error: %v", err)
	}

	if g.serverConfig.WGInterface == "" {
		g.serverConfig.WGInterface = "wg0"
	}

	cmd := exec.Command("wg", "set", g.serverConfig.WGInterface,
		"peer", publicKey,
		"allowed-ips", clientConfig.ClientIP+"/32")

	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("error adding client to server configuration: %v | output: %s", err, output)
	}

	return wgConfigPath, nil
}

type ServerConfig struct {
	DNS         string
	PublicKey   string
	PublicIP    string
	Port        int
	AllowedIPs  string
	WGInterface string
}

type ClientConfig struct {
	ClientIP            string
	PrivateKey          string
	PresistentKeepAlive int
}

type GenerationConfig struct {
	Dir string
}
