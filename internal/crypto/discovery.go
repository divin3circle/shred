package crypto

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type WalletMetadata struct {
	EVMAddress string    `json:"evm_address"`
	CreatedAt  time.Time `json:"created_at"`
	AccountID    string            `json:"account_id,omitempty"`
	Network      string            `json:"network"`
	TokenAliases map[string]string `json:"token_aliases,omitempty"`
}

type WalletInfo struct {
	FilePath   string
	FileName   string
	EVMAddress string
	CreatedAt  time.Time
	AccountID  string
	Network    string
}

func GetWalletDirectory() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get config directory: %w", err)
	}
	walletDir := filepath.Join(configDir, "shred")
	return walletDir, nil
}

func GetWalletPath(identifier string) (string, error) {
	walletDir, err := GetWalletDirectory()
	if err != nil {
		return "", err
	}
	return filepath.Join(walletDir, fmt.Sprintf("wallet-%s.dat", identifier)), nil
}

func GetMetadataPath(walletPath string) string {
	return strings.TrimSuffix(walletPath, ".dat") + ".meta"
}

func SaveWalletMetadata(walletPath string, metadata WalletMetadata) error {
	metaPath := GetMetadataPath(walletPath)
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	return os.WriteFile(metaPath, data, 0600)
}

func LoadWalletMetadata(walletPath string) (WalletMetadata, error) {
	metaPath := GetMetadataPath(walletPath)
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return WalletMetadata{}, fmt.Errorf("failed to read metadata: %w", err)
	}

	var metadata WalletMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return WalletMetadata{}, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return metadata, nil
}

func ListWallets() ([]WalletInfo, error) {
	walletDir, err := GetWalletDirectory()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(walletDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create wallet directory: %w", err)
	}

	entries, err := os.ReadDir(walletDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read wallet directory: %w", err)
	}

	var wallets []WalletInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if !strings.HasSuffix(entry.Name(), ".dat") {
			continue
		}
		
		if !strings.HasPrefix(entry.Name(), "wallet") {
			continue
		}

		walletPath := filepath.Join(walletDir, entry.Name())
		
		fileInfo, err := os.Stat(walletPath)
		if err != nil {
			continue
		}
		
		if !fileInfo.Mode().IsRegular() {
			continue
		}
		
		info := fileInfo

		metadata, err := LoadWalletMetadata(walletPath)
		if err != nil {
			wallets = append(wallets, WalletInfo{
				FilePath:   walletPath,
				FileName:   entry.Name(),
				EVMAddress: "Unknown (decrypt to view)",
				CreatedAt:  info.ModTime(),
				AccountID:  "",
				Network:    "testnet",
			})
			continue
		}

		wallets = append(wallets, WalletInfo{
			FilePath:   walletPath,
			FileName:   entry.Name(),
			EVMAddress: metadata.EVMAddress,
			CreatedAt:  metadata.CreatedAt,
			AccountID:  metadata.AccountID,
			Network:    metadata.Network,
		})
	}

	return wallets, nil
}
