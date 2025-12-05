package hedera

import (
	"encoding/json"
	"fmt"
	"net/http"

	sdk "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

type Client struct {
	Client *sdk.Client
}

func NewClient() (*Client, error) {
	client := sdk.ClientForTestnet()
	return &Client{Client: client}, nil
}

type AccountInfo struct {
	Balance sdk.Hbar
	Tokens  []TokenBalance
}

type TokenBalance struct {
	TokenID string
	Balance uint64
}

func (c *Client) GetAccountBalance(accountID sdk.AccountID) (AccountInfo, error) {
	balance, err := sdk.NewAccountBalanceQuery().
		SetAccountID(accountID).
		Execute(c.Client)
	if err != nil {
		return AccountInfo{}, err
	}

	var tokens []TokenBalance
	for tokenID, bal := range balance.Tokens.GetAll() {
		tokens = append(tokens, TokenBalance{
			TokenID: tokenID,
			Balance: bal,
		})
	}

	return AccountInfo{
		Balance: balance.Hbars,
		Tokens:  tokens,
	}, nil
}

type MirrorAccountResponse struct {
	Accounts []struct {
		Account string `json:"account"`
	} `json:"accounts"`
}

type MirrorAccountDetailResponse struct {
	Account string `json:"account"`
	Balance struct {
		Balance int64 `json:"balance"`
	} `json:"balance"`
	EVMAddress   string              `json:"evm_address"`
	Transactions []MirrorTransaction `json:"transactions"`
	Links        MirrorLinks         `json:"links"`
}

type MirrorTransaction struct {
	TransactionID      string           `json:"transaction_id"`
	ConsensusTimestamp string           `json:"consensus_timestamp"`
	Result             string           `json:"result"`
	Name               string           `json:"name"`
	Transfers          []MirrorTransfer `json:"transfers"`
	MemoBase64         string           `json:"memo_base64"`
}

type MirrorTransfer struct {
	Account string `json:"account"`
	Amount  int64  `json:"amount"`
}

type MirrorLinks struct {
	Next string `json:"next"`
}

func (c *Client) GetAccountIDFromPublicKey(publicKey string) (string, error) {
	url := fmt.Sprintf("https://testnet.mirrornode.hedera.com/api/v1/accounts?account.publickey=%s", publicKey)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("mirror node returned status: %d", resp.StatusCode)
	}

	var result MirrorAccountResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Accounts) > 0 {
		return result.Accounts[0].Account, nil
	}

	return "", nil
}

func (c *Client) GetAccountIDFromEVMAddress(evmAddress string) (string, error) {
	address := evmAddress
	if len(address) >= 2 && address[0:2] == "0x" {
		address = address[2:]
	}
	
	url := fmt.Sprintf("https://testnet.mirrornode.hedera.com/api/v1/accounts/%s", address)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("mirror node returned status: %d", resp.StatusCode)
	}

	var result MirrorAccountDetailResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Account != "" {
		return result.Account, nil
	}

	return "", nil
}

func (c *Client) GetAccountInfoWithTransactions(evmAddress string, nextURL string) (*MirrorAccountDetailResponse, error) {
	var url string
	if nextURL != "" {
		url = "https://testnet.mirrornode.hedera.com" + nextURL
	} else {
		address := evmAddress
		if len(address) >= 2 && address[0:2] == "0x" {
			address = address[2:]
		}
		url = fmt.Sprintf("https://testnet.mirrornode.hedera.com/api/v1/accounts/%s", address)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("mirror node returned status: %d", resp.StatusCode)
	}

	var result MirrorAccountDetailResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) TransferHbar(senderID, recipientID string, amount float64, privateKey string) (string, error) {
	sender, err := sdk.AccountIDFromString(senderID)
	if err != nil {
		return "", fmt.Errorf("invalid sender ID: %w", err)
	}

	recipient, err := sdk.AccountIDFromString(recipientID)
	if err != nil {
		return "", fmt.Errorf("invalid recipient ID: %w", err)
	}

	key, err := sdk.PrivateKeyFromString(privateKey)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %w", err)
	}

	txID := sdk.TransactionIDGenerate(sender)
	tx, err := sdk.NewTransferTransaction().
		SetTransactionID(txID).
		AddHbarTransfer(sender, sdk.NewHbar(-amount)).
		AddHbarTransfer(recipient, sdk.NewHbar(amount)).
		SetTransactionMemo("Sent via shred").
		FreezeWith(c.Client)
	if err != nil {
		return "", fmt.Errorf("failed to create transaction: %w", err)
	}

	tx.Sign(key)

	resp, err := tx.Execute(c.Client)
	if err != nil {
		return "", fmt.Errorf("failed to execute transaction: %w", err)
	}

	receipt, err := resp.GetReceipt(c.Client)
	if err != nil {
		return "", fmt.Errorf("failed to get receipt: %w", err)
	}

	if receipt.Status != sdk.StatusSuccess {
		return "", fmt.Errorf("transaction failed with status: %s", receipt.Status)
	}

	return resp.TransactionID.String(), nil
}

func (c *Client) TransferToken(senderID, recipientID, tokenID string, amount float64, decimals int, privateKey string) (string, error) {
	sender, err := sdk.AccountIDFromString(senderID)
	if err != nil {
		return "", fmt.Errorf("invalid sender ID: %w", err)
	}

	recipient, err := sdk.AccountIDFromString(recipientID)
	if err != nil {
		return "", fmt.Errorf("invalid recipient ID: %w", err)
	}

	token, err := sdk.TokenIDFromString(tokenID)
	if err != nil {
		return "", fmt.Errorf("invalid token ID: %w", err)
	}

	key, err := sdk.PrivateKeyFromString(privateKey)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %w", err)
	}

	rawAmount := int64(amount * float64(1))
	multiplier := int64(1)
	for i := 0; i < decimals; i++ {
		multiplier *= 10
	}
	rawAmount = int64(amount * float64(multiplier))

	txID := sdk.TransactionIDGenerate(sender)
	tx, err := sdk.NewTransferTransaction().
		SetTransactionID(txID).
		AddTokenTransfer(token, sender, -rawAmount).
		AddTokenTransfer(token, recipient, rawAmount).
		SetTransactionMemo("Sent via shred").
		FreezeWith(c.Client)
	if err != nil {
		return "", fmt.Errorf("failed to create transaction: %w", err)
	}

	tx.Sign(key)

	resp, err := tx.Execute(c.Client)
	if err != nil {
		return "", fmt.Errorf("failed to execute transaction: %w", err)
	}

	receipt, err := resp.GetReceipt(c.Client)
	if err != nil {
		return "", fmt.Errorf("failed to get receipt: %w", err)
	}

	if receipt.Status != sdk.StatusSuccess {
		return "", fmt.Errorf("transaction failed with status: %s", receipt.Status)
	}

	return resp.TransactionID.String(), nil
}

func (c *Client) Close() error {
	return c.Client.Close()
}
