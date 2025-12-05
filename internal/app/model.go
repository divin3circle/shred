package app

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/divin3circle/shred/internal/crypto"
	hedera_client "github.com/divin3circle/shred/internal/hedera"
)

type SessionState int

const (
	StateWelcome SessionState = iota
	StateWalletList
	StateWalletUnlock
	StateCreate
	StateVerify
	StatePassword
	StateDashboard
	StateReceive
	StateTokenMenu
	StateLocked
	StateSendSelectToken
	StateSendRecipient
	StateSendAmount
	StateSendConfirm
	StateSendSigning
	StateHistory
)

type Model struct {
	State        SessionState
	Width        int
	Height       int
	
	Mnemonic     []byte
	Wallet       *crypto.Wallet
	
	Input        textinput.Model
	
	MnemonicWords []string
	VerifyIndices []int
	CurrentVerifyIndex int
	
	LastActivity time.Time
	
	HederaClient *hedera_client.Client
	Balance      string
	AccountID    string
	EVMAddress   string
	
	TokenBalances []hedera_client.TokenBalance
	TokenAliases  map[string]string
	
	SelectedTokenIndex int
	TokenListCursor    int
	
	SendSelectedToken hedera_client.TokenBalance
	SendRecipient     string
	SendAmount        string
	SendMemo          string
	SendError         string
	SendSuccess       string
	
	HistoryTransactions []hedera_client.MirrorTransaction
	HistoryNextURL      string
	HistoryPrevURLs     []string
	HistoryIsLoading    bool
	HistoryError        string
	
	IsRefreshing bool
	RefreshError string
	
	AvailableWallets []crypto.WalletInfo
	SelectedWalletIndex int
	SelectedWalletPath string
	
	ErrorMessage string
}

func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Type here..."
	ti.Focus()

	return Model{
		State:        StateWelcome,
		Input:        ti,
		LastActivity: time.Now(),
	}
}

func (m Model) Init() tea.Cmd {
	return checkForWallets
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg, tea.MouseMsg:
		m.LastActivity = time.Now()
	}

	if time.Since(m.LastActivity) > 10*time.Minute {
		m.State = StateLocked
		if m.Wallet != nil {
			m.Wallet.Wipe()
			m.Wallet = nil
		}
		m.Mnemonic = nil
		m.MnemonicWords = nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	}

	switch msg := msg.(type) {
	case walletsFoundMsg:
		if msg.Error != nil {
			m.State = StateWelcome
			return m, nil
		}
		if len(msg.Wallets) > 0 {
			m.AvailableWallets = msg.Wallets
			m.SelectedWalletIndex = 0
			m.State = StateWalletList
			return m, nil
		}
		m.State = StateWelcome
		return m, nil
	}

	switch m.State {
	case StateWelcome:
		return m.updateWelcome(msg)
	case StateWalletList:
		return m.updateWalletList(msg)
	case StateWalletUnlock:
		return m.updateWalletUnlock(msg)
	case StateCreate:
		return m.updateCreate(msg)
	case StateVerify:
		return m.updateVerify(msg)
	case StatePassword:
		return m.updatePassword(msg)
	case StateDashboard:
		return m.updateDashboard(msg)
	case StateReceive:
		return m.updateReceive(msg)
	case StateTokenMenu:
		return m.updateTokenMenu(msg)
	case StateSendSelectToken:
		return m.updateSendSelectToken(msg)
	case StateSendRecipient:
		return m.updateSendRecipient(msg)
	case StateSendAmount:
		return m.updateSendAmount(msg)
	case StateSendConfirm:
		return m.updateSendConfirm(msg)
	case StateSendSigning:
		return m.updateSendSigning(msg)
	case StateHistory:
		return m.updateHistory(msg)
	}

	return m, nil
}

func (m Model) View() string {
	switch m.State {
	case StateWelcome:
		return m.viewWelcome()
	case StateWalletList:
		return m.viewWalletList()
	case StateWalletUnlock:
		return m.viewWalletUnlock()
	case StateCreate:
		return m.viewCreate()
	case StateVerify:
		return m.viewVerify()
	case StatePassword:
		return m.viewPassword()
	case StateDashboard:
		return m.viewDashboard()
	case StateReceive:
		return m.viewReceive()
	case StateTokenMenu:
		return m.viewTokenMenu()
	case StateSendSelectToken:
		return m.viewSendSelectToken()
	case StateSendRecipient:
		return m.viewSendRecipient()
	case StateSendAmount:
		return m.viewSendAmount()
	case StateSendConfirm:
		return m.viewSendConfirm()
	case StateSendSigning:
		return m.viewSendSigning()
	case StateHistory:
		return m.viewHistory()
	}
	return "Unknown State"
}
