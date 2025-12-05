package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/argon2"
)

const (
	argon2Time    = 3
	argon2Memory  = 64 * 1024
	argon2Threads = 4
	keyLen        = 32
	saltLen       = 16
)

type EncryptedWallet struct {
	Salt       []byte `json:"salt"`
	Nonce      []byte `json:"nonce"`
	Ciphertext []byte `json:"ciphertext"`
}

func SaveWallet(mnemonic []byte, passphrase string, path string) error {
	salt := make([]byte, saltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return err
	}

	key := argon2.IDKey([]byte(passphrase), salt, argon2Time, argon2Memory, argon2Threads, keyLen)

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	ciphertext := gcm.Seal(nil, nonce, mnemonic, nil)

	walletData := EncryptedWallet{
		Salt:       salt,
		Nonce:      nonce,
		Ciphertext: ciphertext,
	}

	data, err := json.Marshal(walletData)
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func LoadWallet(passphrase string, path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var walletData EncryptedWallet
	if err := json.Unmarshal(data, &walletData); err != nil {
		return nil, err
	}

	key := argon2.IDKey([]byte(passphrase), walletData.Salt, argon2Time, argon2Memory, argon2Threads, keyLen)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, walletData.Nonce, walletData.Ciphertext, nil)
	if err != nil {
		return nil, errors.New("invalid passphrase or corrupted wallet file")
	}

	return plaintext, nil
}
