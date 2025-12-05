package app

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/divin3circle/shred/internal/crypto"
	hedera_client "github.com/divin3circle/shred/internal/hedera"
	sdk "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

type refreshAccountMsg struct {
	AccountID    string
	Balance      string
	Tokens       []hedera_client.TokenBalance
	Error        error
}

type walletsFoundMsg struct {
	Wallets []crypto.WalletInfo
	Error   error
}

func checkForWallets() tea.Msg {
	wallets, err := crypto.ListWallets()
	return walletsFoundMsg{
		Wallets: wallets,
		Error:   err,
	}
}

func (m Model) updateWalletList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.SelectedWalletIndex > 0 {
				m.SelectedWalletIndex--
			}
			return m, nil
		case "down", "j":
			if m.SelectedWalletIndex < len(m.AvailableWallets)-1 {
				m.SelectedWalletIndex++
			}
			return m, nil
		case "enter":
			if len(m.AvailableWallets) > 0 && m.SelectedWalletIndex < len(m.AvailableWallets) {
				m.SelectedWalletPath = m.AvailableWallets[m.SelectedWalletIndex].FilePath
				m.State = StateWalletUnlock
				m.Input.Reset()
				m.Input.Placeholder = "Enter Passphrase"
				m.Input.EchoMode = textinput.EchoPassword
			}
			return m, nil
		case "n":
			mnemonic, err := crypto.NewMnemonic()
			if err != nil {
				return m, tea.Quit
			}
			m.Mnemonic = mnemonic
			m.MnemonicWords = strings.Split(string(mnemonic), " ")
			m.State = StateCreate
			return m, nil
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) updateWalletUnlock(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			passphrase := m.Input.Value()
			if passphrase == "" {
				return m, cmd
			}

			mnemonic, err := crypto.LoadWallet(passphrase, m.SelectedWalletPath)
			if err != nil {
				m.Input.Reset()
				m.Input.Placeholder = "Invalid passphrase. Try again:"
				m.Input.EchoMode = textinput.EchoPassword
				return m, cmd
			}

			m.Mnemonic = mnemonic

			client, err := hedera_client.NewClient()
			if err == nil {
				m.HederaClient = client

				ecdsaKey, err := crypto.DeriveECDSAKey(m.Mnemonic)
				if err == nil {
					m.EVMAddress = crypto.CalculateEVMAddress(ecdsaKey)

					accountID, err := client.GetAccountIDFromEVMAddress(m.EVMAddress)
					if err == nil && accountID != "" {
						m.AccountID = accountID

						id, _ := sdk.AccountIDFromString(accountID)
						info, err := client.GetAccountBalance(id)
						if err == nil {
							m.Balance = info.Balance.String()
							m.TokenBalances = info.Tokens
						}
					} else {
						m.AccountID = "Unverified"
						m.Balance = "0.00 ℏ"
					}

					metadata, err := crypto.LoadWalletMetadata(m.SelectedWalletPath)
					if err != nil {
						metadata = crypto.WalletMetadata{
							EVMAddress: m.EVMAddress,
							CreatedAt:  time.Now(),
							Network:    "testnet",
						}
					}
					metadata.AccountID = m.AccountID
					metadata.EVMAddress = m.EVMAddress
					if metadata.TokenAliases != nil {
						m.TokenAliases = metadata.TokenAliases
					} else {
						m.TokenAliases = make(map[string]string)
					}
					crypto.SaveWalletMetadata(m.SelectedWalletPath, metadata)
				}
			}

			m.State = StateDashboard
		case "esc":
			m.State = StateWalletList
			m.Input.Reset()
			return m, nil
		}
	}
	return m, cmd
}

func (m Model) updateWelcome(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch strings.ToLower(msg.String()) {
		case "n":
			mnemonic, err := crypto.NewMnemonic()
			if err != nil {
				return m, tea.Quit
			}
			m.Mnemonic = mnemonic
			m.MnemonicWords = strings.Split(string(mnemonic), " ")
			m.State = StateCreate
			return m, nil
		case "i":
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) updateCreate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.VerifyIndices = []int{3, 10, 18}
			m.CurrentVerifyIndex = 0
			m.State = StateVerify
			m.Input.Reset()
			m.Input.Placeholder = fmt.Sprintf("Word #%d", m.VerifyIndices[0]+1)
			return m, nil
		}
	}
	return m, nil
}

func (m Model) updateVerify(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			inputWord := strings.TrimSpace(strings.ToLower(m.Input.Value()))
			targetIndex := m.VerifyIndices[m.CurrentVerifyIndex]
			targetWord := m.MnemonicWords[targetIndex]

			if inputWord == targetWord {
				m.CurrentVerifyIndex++
				m.Input.Reset()
				
				if m.CurrentVerifyIndex >= len(m.VerifyIndices) {
					m.State = StatePassword
					m.Input.Placeholder = "Enter Passphrase"
					m.Input.EchoMode = textinput.EchoPassword
				} else {
					m.Input.Placeholder = fmt.Sprintf("Word #%d", m.VerifyIndices[m.CurrentVerifyIndex]+1)
				}
			} else {
				m.Input.Reset()
				m.Input.Placeholder = fmt.Sprintf("Incorrect! Try Word #%d again", targetIndex+1)
			}
		}
	}
	return m, cmd
}

func (m Model) updatePassword(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			passphrase := m.Input.Value()
			
			ecdsaKey, err := crypto.DeriveECDSAKey(m.Mnemonic)
			if err != nil {
				m.Input.Reset()
				m.Input.Placeholder = "Error deriving keys. Try again:"
				return m, cmd
			}
			
			evmAddress := crypto.CalculateEVMAddress(ecdsaKey)
			m.EVMAddress = evmAddress
			
			identifier := evmAddress
			if len(identifier) >= 2 && identifier[0:2] == "0x" {
				identifier = identifier[2:]
			}
			
			walletPath, err := crypto.GetWalletPath(identifier)
			if err != nil {
				m.ErrorMessage = fmt.Sprintf("Error getting wallet path: %v", err)
				m.Input.Reset()
				m.Input.Placeholder = "Error getting wallet path. Try again:"
				return m, cmd
			}
			
			err = crypto.SaveWallet(m.Mnemonic, passphrase, walletPath)
			if err != nil {
				m.ErrorMessage = fmt.Sprintf("Failed to save wallet to %s: %v", walletPath, err)
				m.Input.Reset()
				m.Input.Placeholder = fmt.Sprintf("Error: %v", err)
				return m, cmd
			}
			
			if _, err := os.Stat(walletPath); err != nil {
				m.ErrorMessage = fmt.Sprintf("Wallet file not found after save: %s (error: %v)", walletPath, err)
				m.Input.Reset()
				m.Input.Placeholder = "File creation failed. Try again:"
				return m, cmd
			}
			
			m.ErrorMessage = ""
			
			metadata := crypto.WalletMetadata{
				EVMAddress: evmAddress,
				CreatedAt:  time.Now(),
				Network:    "testnet",
			}
			err = crypto.SaveWalletMetadata(walletPath, metadata)
			if err != nil {
			}
			
			client, err := hedera_client.NewClient()
			if err == nil {
				m.HederaClient = client
				
				accountID, err := client.GetAccountIDFromEVMAddress(evmAddress)
				if err == nil && accountID != "" {
					m.AccountID = accountID
					
					metadata.AccountID = accountID
					crypto.SaveWalletMetadata(walletPath, metadata)
					
					id, _ := sdk.AccountIDFromString(accountID)
					info, err := client.GetAccountBalance(id)
					if err == nil {
						m.Balance = info.Balance.String()
						m.TokenBalances = info.Tokens
					}
				} else {
					m.AccountID = "Unverified"
					m.Balance = "0.00 ℏ"
				}
			}
			
			m.State = StateDashboard
		}
	}
	return m, cmd
}

func refreshAccountCmd(evmAddress string, client *hedera_client.Client) tea.Cmd {
	return func() tea.Msg {
		if client == nil || evmAddress == "" {
			return refreshAccountMsg{Error: fmt.Errorf("client or EVM address not available")}
		}

		accountID, err := client.GetAccountIDFromEVMAddress(evmAddress)
		if err != nil {
			return refreshAccountMsg{Error: fmt.Errorf("failed to query account: %w", err)}
		}

		if accountID == "" {
			return refreshAccountMsg{
				AccountID: "Unverified",
				Balance:   "0.00 ℏ",
			}
		}

		id, err := sdk.AccountIDFromString(accountID)
		if err != nil {
			return refreshAccountMsg{Error: fmt.Errorf("invalid account ID: %w", err)}
		}

		info, err := client.GetAccountBalance(id)
		if err != nil {
			return refreshAccountMsg{
				AccountID: accountID,
				Balance:   "Error",
				Error:     fmt.Errorf("failed to fetch balance: %w", err),
			}
		}

		return refreshAccountMsg{
			AccountID: accountID,
			Balance:   info.Balance.String(),
			Tokens:    info.Tokens,
		}
	}
}

func (m Model) updateDashboard(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch strings.ToLower(msg.String()) {
		case "f":
			if m.HederaClient != nil && m.EVMAddress != "" && !m.IsRefreshing {
				m.IsRefreshing = true
				m.RefreshError = ""
				return m, refreshAccountCmd(m.EVMAddress, m.HederaClient)
			}
		case "s":
			m.State = StateSendSelectToken
			m.SelectedTokenIndex = 0
			return m, nil
		case "h":
			m.State = StateHistory
			m.HistoryIsLoading = true
			m.HistoryError = ""
			m.HistoryPrevURLs = []string{}
			return m, fetchHistoryCmd(m.EVMAddress, "", m.HederaClient)
		case "r":
			m.State = StateReceive
			return m, nil
		case "t":
			if len(m.TokenBalances) > 0 {
				m.State = StateTokenMenu
				m.SelectedTokenIndex = 0
				m.Input.Reset()
				m.Input.Placeholder = "Enter alias for selected token..."
			}
			return m, nil
		case "q":
			return m, tea.Quit
		}
	case refreshAccountMsg:
		m.IsRefreshing = false
		if msg.Error != nil {
			m.RefreshError = msg.Error.Error()
		} else {
			m.AccountID = msg.AccountID
			m.Balance = msg.Balance
			m.TokenBalances = msg.Tokens
			m.RefreshError = ""
		}
		return m, nil
	}
	return m, nil
}

func (m Model) updateReceive(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch strings.ToLower(msg.String()) {
		case "esc":
			m.State = StateDashboard
			return m, nil
		case "c":
			if m.AccountID != "" && m.AccountID != "Unverified" && m.AccountID != "Inactive" {
				clipboard.WriteAll(m.AccountID)
			}
			return m, nil
		case "e":
			if m.EVMAddress != "" {
				clipboard.WriteAll(m.EVMAddress)
			}
			return m, nil
		}
	}
	return m, nil
}

func (m Model) updateTokenMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.State = StateDashboard
			return m, nil
		case "up", "k":
			if m.SelectedTokenIndex > 0 {
				m.SelectedTokenIndex--
			}
			return m, nil
		case "down", "j":
			if m.SelectedTokenIndex < len(m.TokenBalances)-1 {
				m.SelectedTokenIndex++
			}
			return m, nil
		case "enter":
			alias := m.Input.Value()
			if alias != "" {
				if m.TokenAliases == nil {
					m.TokenAliases = make(map[string]string)
				}
				tokenID := m.TokenBalances[m.SelectedTokenIndex].TokenID
				m.TokenAliases[tokenID] = alias
				m.Input.Reset()
				
				if m.SelectedWalletPath != "" {
					metadata, err := crypto.LoadWalletMetadata(m.SelectedWalletPath)
					if err == nil {
						metadata.TokenAliases = m.TokenAliases
						crypto.SaveWalletMetadata(m.SelectedWalletPath, metadata)
					}
				}
			}
			return m, nil
		}
	}
	return m, cmd
}

func (m Model) updateSendSelectToken(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.State = StateDashboard
			return m, nil
		case "up", "k":
			if m.SelectedTokenIndex > 0 {
				m.SelectedTokenIndex--
			}
			return m, nil
		case "down", "j":
			if m.SelectedTokenIndex < len(m.TokenBalances) {
				m.SelectedTokenIndex++
			}
			return m, nil
		case "enter":
			if m.SelectedTokenIndex == 0 {
				m.SendSelectedToken = hedera_client.TokenBalance{TokenID: ""}
			} else {
				m.SendSelectedToken = m.TokenBalances[m.SelectedTokenIndex-1]
			}
			
			m.State = StateSendRecipient
			m.Input.Reset()
			m.Input.Placeholder = "Enter Recipient (Account ID or EVM Address)"
			m.Input.EchoMode = textinput.EchoNormal
			m.SendError = ""
			return m, nil
		}
	}
	return m, nil
}

func (m Model) updateSendRecipient(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.State = StateSendSelectToken
			return m, nil
		case "enter":
			recipient := strings.TrimSpace(m.Input.Value())
			if recipient == "" {
				return m, nil
			}
			
			if strings.HasPrefix(recipient, "0x") || len(recipient) == 40 || len(recipient) == 42 {
				m.SendRecipient = recipient
				return m, resolveRecipientCmd(recipient, m.HederaClient)
			}
			
			m.SendRecipient = recipient
			m.State = StateSendAmount
			m.Input.Reset()
			m.Input.Placeholder = "Enter Amount"
			return m, nil
		}
	case recipientResolvedMsg:
		if msg.Error != nil {
			m.SendError = msg.Error.Error()
			m.Input.Placeholder = "Invalid address. Try again:"
			m.Input.Reset()
		} else {
			m.SendRecipient = msg.AccountID
			m.State = StateSendAmount
			m.Input.Reset()
			m.Input.Placeholder = "Enter Amount"
			m.SendError = ""
		}
		return m, nil
	}
	return m, cmd
}

type recipientResolvedMsg struct {
	AccountID string
	Error     error
}

func resolveRecipientCmd(address string, client *hedera_client.Client) tea.Cmd {
	return func() tea.Msg {
		accountID, err := client.GetAccountIDFromEVMAddress(address)
		if err != nil {
			return recipientResolvedMsg{Error: err}
		}
		if accountID == "" {
			return recipientResolvedMsg{Error: fmt.Errorf("account not found for address")}
		}
		return recipientResolvedMsg{AccountID: accountID}
	}
}

func (m Model) updateSendAmount(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.State = StateSendRecipient
			return m, nil
		case "enter":
			amount := strings.TrimSpace(m.Input.Value())
			if amount == "" {
				return m, nil
			}
			m.SendAmount = amount
			m.State = StateSendConfirm
			return m, nil
		}
	}
	return m, cmd
}

func (m Model) updateSendConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch strings.ToLower(msg.String()) {
		case "esc":
			m.State = StateSendAmount
			return m, nil
		case "y", "enter":
			m.State = StateSendSigning
			m.Input.Reset()
			m.Input.Placeholder = "Enter Wallet Passphrase to Sign"
			m.Input.EchoMode = textinput.EchoPassword
			return m, nil
		}
	}
	return m, nil
}

func (m Model) updateSendSigning(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.State = StateSendConfirm
			m.Input.Reset()
			m.Input.EchoMode = textinput.EchoNormal
			return m, nil
		case "enter":
			passphrase := m.Input.Value()
			if passphrase == "" {
				return m, nil
			}
			
			mnemonic, err := crypto.LoadWallet(passphrase, m.SelectedWalletPath)
			if err != nil {
				m.SendError = "Invalid passphrase"
				m.Input.Reset()
				return m, nil
			}
			
			ecdsaKey, err := crypto.DeriveECDSAKey(mnemonic)
			if err != nil {
				m.SendError = "Failed to derive key"
				return m, nil
			}
			privateKey := ecdsaKey.String()
			
			return m, sendTransactionCmd(m.HederaClient, m.AccountID, m.SendRecipient, m.SendSelectedToken, m.SendAmount, privateKey)
		}
	case transactionResultMsg:
		if msg.Error != nil {
			m.SendError = msg.Error.Error()
			m.State = StateSendConfirm
		} else {
			m.SendSuccess = fmt.Sprintf("Transaction Sent! ID: %s", msg.TransactionID)
			m.State = StateDashboard
			return m, refreshAccountCmd(m.EVMAddress, m.HederaClient)
		}
		return m, nil
	}
	return m, cmd
}

type transactionResultMsg struct {
	TransactionID string
	Error         error
}

func sendTransactionCmd(client *hedera_client.Client, senderID, recipientID string, token hedera_client.TokenBalance, amountStr string, privateKey string) tea.Cmd {
	return func() tea.Msg {
		var amount float64
		_, err := fmt.Sscanf(amountStr, "%f", &amount)
		if err != nil {
			return transactionResultMsg{Error: fmt.Errorf("invalid amount format")}
		}

		var txID string
		if token.TokenID == "" {
			txID, err = client.TransferHbar(senderID, recipientID, amount, privateKey)
		} else {
			txID, err = client.TransferToken(senderID, recipientID, token.TokenID, amount, 0, privateKey)
		}

		if err != nil {
			return transactionResultMsg{Error: err}
		}
		
		return transactionResultMsg{TransactionID: txID}
	}
}

func (m Model) updateHistory(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.State = StateDashboard
			return m, nil
		case "n":
			if m.HistoryNextURL != "" && !m.HistoryIsLoading {
				m.HistoryIsLoading = true
				return m, fetchHistoryCmd(m.EVMAddress, m.HistoryNextURL, m.HederaClient)
			}
		}
	case historyFetchedMsg:
		m.HistoryIsLoading = false
		if msg.Error != nil {
			m.HistoryError = msg.Error.Error()
		} else {
			m.HistoryTransactions = msg.Transactions
			m.HistoryNextURL = msg.NextURL
			m.HistoryError = ""
		}
		return m, nil
	}
	return m, nil
}

type historyFetchedMsg struct {
	Transactions []hedera_client.MirrorTransaction
	NextURL      string
	Error        error
}

func fetchHistoryCmd(evmAddress string, url string, client *hedera_client.Client) tea.Cmd {
	return func() tea.Msg {
		info, err := client.GetAccountInfoWithTransactions(evmAddress, url)
		if err != nil {
			return historyFetchedMsg{Error: err}
		}
		
		next := ""
		if info.Links.Next != "" {
			next = info.Links.Next
		}
		
		return historyFetchedMsg{
			Transactions: info.Transactions,
			NextURL:      next,
		}
	}
}
