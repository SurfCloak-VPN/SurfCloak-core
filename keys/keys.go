package keys

import (
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
