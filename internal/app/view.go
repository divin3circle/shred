package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/divin3circle/shred/internal/crypto"
	"github.com/mdp/qrterminal/v3"
)

var (
	styleBox = lipgloss.NewStyle().
			Padding(1, 2).
			Height(30)

	styleTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500")).
			Bold(true).
			Align(lipgloss.Center).
			Width(80)

	styleSubTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Bold(true)
)

func (m Model) viewWelcome() string {
	content := fmt.Sprintf(`
%s

No wallet found. Create or import?

[N] Create new wallet
[I] Import seed phrase
[Q] Quit
`, GetStyledLogo())

	boxedContent := styleBox.Render(content)

	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, boxedContent)
}

func (m Model) viewWalletList() string {
	var walletList strings.Builder

	if len(m.AvailableWallets) == 0 {
		content := fmt.Sprintf(`
%s

No wallets found.

[N] Create new wallet
[Q] Quit
`, GetStyledLogo())
		boxedContent := styleBox.Render(content)
		return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, boxedContent)
	}

	walletList.WriteString(fmt.Sprintf("%s\n\n", GetStyledLogo()))
	walletList.WriteString("Available wallets:\n\n")

	for i, wallet := range m.AvailableWallets {
		cursor := "  "
		if i == m.SelectedWalletIndex {
			cursor = "→ "
		}

		evmDisplay := wallet.EVMAddress
		if len(evmDisplay) > 20 {
			evmDisplay = evmDisplay[:10] + "..." + evmDisplay[len(evmDisplay)-8:]
		}

		status := "Unverified"
		if wallet.AccountID != "" {
			status = wallet.AccountID
		}

		walletList.WriteString(fmt.Sprintf("%s[%d] %s\n", cursor, i+1, wallet.FileName))
		walletList.WriteString(fmt.Sprintf("    EVM: %s\n", evmDisplay))
		walletList.WriteString(fmt.Sprintf("    Status: %s\n\n", status))
	}

	content := fmt.Sprintf(`
%s

[↑↓] Navigate  [Enter] Select  [N] New Wallet  [Q] Quit
`, walletList.String())

	boxedContent := styleBox.Render(content)
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, boxedContent)
}

func (m Model) viewWalletUnlock() string {
	var walletInfo crypto.WalletInfo
	if m.SelectedWalletIndex < len(m.AvailableWallets) {
		walletInfo = m.AvailableWallets[m.SelectedWalletIndex]
	}

	evmDisplay := walletInfo.EVMAddress
	if len(evmDisplay) > 20 {
		evmDisplay = evmDisplay[:10] + "..." + evmDisplay[len(evmDisplay)-8:]
	}

	content := fmt.Sprintf(`
%s

Unlock wallet: %s
EVM Address: %s

Enter passphrase:

%s

[Enter] Unlock  [Esc] Back
`, styleTitle.Render("Unlock Wallet"), walletInfo.FileName, evmDisplay, m.Input.View())

	boxedContent := styleBox.Render(content)
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, boxedContent)
}

func (m Model) viewCreate() string {
	var wordsView strings.Builder
	for i, word := range m.MnemonicWords {
		wordsView.WriteString(fmt.Sprintf("%2d. %-10s ", i+1, word))
		if (i+1)%4 == 0 {
			wordsView.WriteString("\n")
		}
	}

	content := fmt.Sprintf(`
%s

Write down these 24 words. They are the ONLY way to recover your funds.

%s

[Enter] I have written them down
`, styleTitle.Render("Secret Recovery Phrase"), wordsView.String())

	boxedContent := styleBox.Render(content)
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, boxedContent)
}

func (m Model) viewVerify() string {
	targetIndex := m.VerifyIndices[m.CurrentVerifyIndex]
	content := fmt.Sprintf(`
%s

Verify your recovery phrase.
Enter word #%d:

%s
`, styleTitle.Render("Verify Phrase"), targetIndex+1, m.Input.View())

	boxedContent := styleBox.Render(content)
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, boxedContent)
}

func (m Model) viewPassword() string {
	errorMsg := ""
	if m.ErrorMessage != "" {
		errorMsg = fmt.Sprintf("\n⚠️  %s\n", m.ErrorMessage)
	}

	content := fmt.Sprintf(`
%s

Set a strong passphrase to encrypt your wallet on disk.%s

%s
`, styleTitle.Render("Set Passphrase"), errorMsg, m.Input.View())

	boxedContent := styleBox.Render(content)
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, boxedContent)
}

func (m Model) viewDashboard() string {
	statusLine := ""
	if m.IsRefreshing {
		statusLine = "\n🔄 Refreshing account information...\n"
	} else if m.RefreshError != "" {
		statusLine = fmt.Sprintf("\n⚠️  Error: %s\n", m.RefreshError)
	}

	content := fmt.Sprintf(`
%s

Account: %s
Balance: %s
EVM Alias: %s%s

[s] Send   [r] Receive   [t] Tokens   [f] Refresh   [h] History   [q] Quit
`, styleTitle.Render("shred Dashboard")+"\n"+styleSubTitle.Render("Network: Testnet"), m.AccountID, m.Balance, m.EVMAddress, statusLine)

	if len(m.TokenBalances) > 0 {
		content += "\nTokens:\n"
		for _, token := range m.TokenBalances {
			alias := ""
			if m.TokenAliases != nil {
				if a, ok := m.TokenAliases[token.TokenID]; ok {
					alias = fmt.Sprintf(" (%s)", a)
				}
			}
			content += fmt.Sprintf("- %s: %d%s\n", token.TokenID, token.Balance, alias)
		}
	}

	boxedContent := styleBox.Render(content)
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, boxedContent)
}

func (m Model) viewReceive() string {
	var content strings.Builder

	content.WriteString(styleTitle.Render("Receive Funds") + "\n\n")

	content.WriteString(styleSubTitle.Render("Hedera Account ID") + "\n")
	content.WriteString(m.AccountID + "\n\n")

	content.WriteString(styleSubTitle.Render("EVM Address") + "\n")
	content.WriteString(m.EVMAddress + "\n\n")

	qrContent := m.AccountID
	if m.AccountID == "Unverified" || m.AccountID == "Inactive" {
		qrContent = m.EVMAddress
	}

	config := qrterminal.Config{
		Level:     qrterminal.M,
		Writer:    &content,
		BlackChar: qrterminal.BLACK,
		WhiteChar: qrterminal.WHITE,
		QuietZone: 1,
	}
	qrterminal.GenerateWithConfig(qrContent, config)

	content.WriteString("\n\n[c] Copy Account ID  [e] Copy EVM Address  [Esc] Back\n")

	boxedContent := styleBox.Render(content.String())
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, boxedContent)
}

func (m Model) viewTokenMenu() string {
	var content strings.Builder

	content.WriteString(styleTitle.Render("Manage Tokens") + "\n\n")
	content.WriteString("Select a token to add an alias:\n\n")

	for i, token := range m.TokenBalances {
		cursor := "  "
		if i == m.SelectedTokenIndex {
			cursor = "→ "
		}

		alias := ""
		if m.TokenAliases != nil {
			if a, ok := m.TokenAliases[token.TokenID]; ok {
				alias = fmt.Sprintf(" (%s)", a)
			}
		}

		content.WriteString(fmt.Sprintf("%s%s: %d%s\n", cursor, token.TokenID, token.Balance, alias))
	}

	content.WriteString("\n" + m.Input.View() + "\n")
	content.WriteString("\n[↑↓] Navigate  [Enter] Save Alias  [Esc] Back\n")

	boxedContent := styleBox.Render(content.String())
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, boxedContent)
}

func (m Model) viewSendSelectToken() string {
	var content strings.Builder
	content.WriteString(styleTitle.Render("Send Funds - Select Asset") + "\n\n")

	cursor := "  "
	if m.SelectedTokenIndex == 0 {
		cursor = "→ "
	}
	content.WriteString(fmt.Sprintf("%sHBAR (Balance: %s)\n", cursor, m.Balance))

	for i, token := range m.TokenBalances {
		cursor = "  "
		if i+1 == m.SelectedTokenIndex {
			cursor = "→ "
		}

		alias := ""
		if m.TokenAliases != nil {
			if a, ok := m.TokenAliases[token.TokenID]; ok {
				alias = fmt.Sprintf(" (%s)", a)
			}
		}

		content.WriteString(fmt.Sprintf("%s%s: %d%s\n", cursor, token.TokenID, token.Balance, alias))
	}

	content.WriteString("\n[↑↓] Navigate  [Enter] Select  [Esc] Cancel\n")
	boxedContent := styleBox.Render(content.String())
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, boxedContent)
}

func (m Model) viewSendRecipient() string {
	errorMsg := ""
	if m.SendError != "" {
		errorMsg = fmt.Sprintf("\n⚠️  %s\n", m.SendError)
	}

	assetName := "HBAR"
	if m.SendSelectedToken.TokenID != "" {
		assetName = m.SendSelectedToken.TokenID
		if m.TokenAliases != nil {
			if a, ok := m.TokenAliases[assetName]; ok {
				assetName = fmt.Sprintf("%s (%s)", assetName, a)
			}
		}
	}

	content := fmt.Sprintf(`
%s

Sending: %s

Enter Recipient (Account ID or EVM Address):
%s
%s
[Enter] Next  [Esc] Back
`, styleTitle.Render("Send Funds - Recipient"), assetName, m.Input.View(), errorMsg)

	boxedContent := styleBox.Render(content)
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, boxedContent)
}

func (m Model) viewSendAmount() string {
	assetName := "HBAR"
	maxAmount := m.Balance

	if m.SendSelectedToken.TokenID != "" {
		assetName = m.SendSelectedToken.TokenID
		if m.TokenAliases != nil {
			if a, ok := m.TokenAliases[assetName]; ok {
				assetName = fmt.Sprintf("%s (%s)", assetName, a)
			}
		}
		maxAmount = fmt.Sprintf("%d", m.SendSelectedToken.Balance)
	}

	content := fmt.Sprintf(`
%s

Sending: %s
To: %s
Max Available: %s

Enter Amount:
%s

[Enter] Next  [Esc] Back
`, styleTitle.Render("Send Funds - Amount"), assetName, m.SendRecipient, maxAmount, m.Input.View())

	boxedContent := styleBox.Render(content)
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, boxedContent)
}

func (m Model) viewSendConfirm() string {
	errorMsg := ""
	if m.SendError != "" {
		errorMsg = fmt.Sprintf("\n⚠️  %s\n", m.SendError)
	}

	assetName := "HBAR"
	if m.SendSelectedToken.TokenID != "" {
		assetName = m.SendSelectedToken.TokenID
		if m.TokenAliases != nil {
			if a, ok := m.TokenAliases[assetName]; ok {
				assetName = fmt.Sprintf("%s (%s)", assetName, a)
			}
		}
	}

	content := fmt.Sprintf(`
%s

Please review your transaction:

Asset:     %s
Amount:    %s
Recipient: %s
Memo:      Sent via shred

%s
[Y/Enter] Confirm & Sign  [Esc] Back
`, styleTitle.Render("Send Funds - Confirmation"), assetName, m.SendAmount, m.SendRecipient, errorMsg)

	boxedContent := styleBox.Render(content)
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, boxedContent)
}

func (m Model) viewSendSigning() string {
	errorMsg := ""
	if m.SendError != "" {
		errorMsg = fmt.Sprintf("\n⚠️  %s\n", m.SendError)
	}

	content := fmt.Sprintf(`
%s

Enter your wallet passphrase to sign the transaction:

%s
%s
[Enter] Sign & Send  [Esc] Cancel
`, styleTitle.Render("Sign Transaction"), m.Input.View(), errorMsg)

	boxedContent := styleBox.Render(content)
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, boxedContent)
}

func (m Model) viewHistory() string {
	var content strings.Builder
	content.WriteString(styleTitle.Render("Transaction History") + "\n\n")

	if m.HistoryIsLoading {
		content.WriteString("Loading transactions...\n")
	} else if m.HistoryError != "" {
		content.WriteString(fmt.Sprintf("⚠️  Error: %s\n", m.HistoryError))
	} else if len(m.HistoryTransactions) == 0 {
		content.WriteString("No transactions found.\n")
	} else {
		for _, tx := range m.HistoryTransactions {
			parts := strings.Split(tx.ConsensusTimestamp, ".")
			tsStr := parts[0]

			var amount int64
			for _, transfer := range tx.Transfers {
				if transfer.Account == m.AccountID {
					amount += transfer.Amount
				}
			}

			amountStr := fmt.Sprintf("%d tℏ", amount)
			if amount > 0 {
				amountStr = "+" + amountStr
			}

			memo := ""
			if tx.MemoBase64 != "" {
				memo = " (Memo)"
			}

			content.WriteString(fmt.Sprintf("%s  %s  %s%s  (%s)\n", tx.TransactionID, tx.Result, amountStr, memo, tsStr))
		}
	}

	content.WriteString("\n")
	if m.HistoryNextURL != "" {
		content.WriteString("[n] Next Page  ")
	}
	content.WriteString("[Esc] Back\n")

	boxedContent := styleBox.Render(content.String())
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, boxedContent)
}
