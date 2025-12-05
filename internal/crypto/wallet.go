package crypto

import (
	"crypto/ed25519"

	sdk "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
	"github.com/tyler-smith/go-bip39"
)

type Wallet struct {
	Mnemonic    []byte
	PrivateKey  ed25519.PrivateKey
	PublicKey   ed25519.PublicKey
	AccountID   sdk.AccountID
	EVMAddress  string
}

func (w *Wallet) Wipe() {
	if w == nil {
		return
	}
	for i := range w.Mnemonic {
		w.Mnemonic[i] = 0
	}
	w.Mnemonic = nil

	for i := range w.PrivateKey {
		w.PrivateKey[i] = 0
	}
	w.PrivateKey = nil
}

func NewMnemonic() ([]byte, error) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return nil, err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, err
	}
	return []byte(mnemonic), nil
}

func ValidateMnemonic(mnemonic []byte) bool {
	return bip39.IsMnemonicValid(string(mnemonic))
}

func DeriveKey(mnemonic []byte) (ed25519.PrivateKey, error) {
	hMnemonic, err := sdk.MnemonicFromString(string(mnemonic))
	if err != nil {
		return nil, err
	}

	key, err := sdk.PrivateKeyFromMnemonic(hMnemonic, "")
	if err != nil {
		return nil, err
	}
	
	return key.Bytes(), nil
}

func DeriveECDSAKey(mnemonic []byte) (sdk.PrivateKey, error) {
	edKeyBytes, err := DeriveKey(mnemonic)
	if err != nil {
		return sdk.PrivateKey{}, err
	}
	
	return sdk.PrivateKeyFromBytesECDSA(edKeyBytes)
}

func CalculateEVMAddress(key sdk.PrivateKey) string {
    return key.PublicKey().ToEvmAddress()
}
