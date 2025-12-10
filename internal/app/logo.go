package app

import (
	"github.com/charmbracelet/lipgloss"
)

const shredLogo = `
 ███████╗██╗  ██╗██████╗ ███████╗██████╗ 
 ██╔════╝██║  ██║██╔══██╗██╔════╝██╔══██╗
 ███████╗███████║██████╔╝█████╗  ██║  ██║
 ╚════██║██╔══██║██╔══██╗██╔══╝  ██║  ██║
 ███████║██║  ██║██║  ██║███████╗██████╔╝
 ╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝╚═════╝ 
`

const shredTagline = "Secure Hedera Terminal Wallet"

func GetStyledLogo() string {
	logoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFA500")).
		Bold(true)

	taglineStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Italic(true)

	return logoStyle.Render(shredLogo) + "\n" +
		lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render(taglineStyle.Render(shredTagline))
}

const shredLogoSimple = `
  ____  _   _ ____  _____ ____  
 / ___|| | | |  _ \| ____|  _ \ 
 \___ \| |_| | |_) |  _| | | | |
  ___) |  _  |  _ <| |___| |_| |
 |____/|_| |_|_| \_\_____|____/ 
`

func GetSimpleLogo() string {
	logoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00D4FF")).
		Bold(true)

	return logoStyle.Render(shredLogoSimple)
}

const shredLogoBlock = `
 ▄▄▄▄▄▄▄▄▄▄▄  ▄         ▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄   
▐░░░░░░░░░░░▌▐░▌       ▐░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░▌  
▐░█▀▀▀▀▀▀▀▀▀ ▐░▌       ▐░▌▐░█▀▀▀▀▀▀▀█░▌▐░█▀▀▀▀▀▀▀▀▀ ▐░█▀▀▀▀▀▀▀█░▌ 
▐░▌          ▐░▌       ▐░▌▐░▌       ▐░▌▐░▌          ▐░▌       ▐░▌ 
▐░█▄▄▄▄▄▄▄▄▄ ▐░█▄▄▄▄▄▄▄█░▌▐░█▄▄▄▄▄▄▄█░▌▐░█▄▄▄▄▄▄▄▄▄ ▐░▌       ▐░▌ 
▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░▌       ▐░▌ 
 ▀▀▀▀▀▀▀▀▀█░▌▐░█▀▀▀▀▀▀▀█░▌▐░█▀▀▀▀█░█▀▀ ▐░█▀▀▀▀▀▀▀▀▀ ▐░▌       ▐░▌ 
          ▐░▌▐░▌       ▐░▌▐░▌     ▐░▌  ▐░▌          ▐░▌       ▐░▌ 
 ▄▄▄▄▄▄▄▄▄█░▌▐░▌       ▐░▌▐░▌      ▐░▌ ▐░█▄▄▄▄▄▄▄▄▄ ▐░█▄▄▄▄▄▄▄█░▌ 
▐░░░░░░░░░░░▌▐░▌       ▐░▌▐░▌       ▐░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░▌  
 ▀▀▀▀▀▀▀▀▀▀▀  ▀         ▀  ▀         ▀  ▀▀▀▀▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀▀▀▀   
`

func GetBlockLogo() string {
	topStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00D4FF"))

	bottomStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#0080FF"))

	lines := []string{
		" ▄▄▄▄▄▄▄▄▄▄▄  ▄         ▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄   ",
		"▐░░░░░░░░░░░▌▐░▌       ▐░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░▌  ",
		"▐░█▀▀▀▀▀▀▀▀▀ ▐░▌       ▐░▌▐░█▀▀▀▀▀▀▀█░▌▐░█▀▀▀▀▀▀▀▀▀ ▐░█▀▀▀▀▀▀▀█░▌ ",
		"▐░▌          ▐░▌       ▐░▌▐░▌       ▐░▌▐░▌          ▐░▌       ▐░▌ ",
		"▐░█▄▄▄▄▄▄▄▄▄ ▐░█▄▄▄▄▄▄▄█░▌▐░█▄▄▄▄▄▄▄█░▌▐░█▄▄▄▄▄▄▄▄▄ ▐░▌       ▐░▌ ",
		"▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░▌       ▐░▌ ",
		" ▀▀▀▀▀▀▀▀▀█░▌▐░█▀▀▀▀▀▀▀█░▌▐░█▀▀▀▀█░█▀▀ ▐░█▀▀▀▀▀▀▀▀▀ ▐░▌       ▐░▌ ",
		"          ▐░▌▐░▌       ▐░▌▐░▌     ▐░▌  ▐░▌          ▐░▌       ▐░▌ ",
		" ▄▄▄▄▄▄▄▄▄█░▌▐░▌       ▐░▌▐░▌      ▐░▌ ▐░█▄▄▄▄▄▄▄▄▄ ▐░█▄▄▄▄▄▄▄█░▌ ",
		"▐░░░░░░░░░░░▌▐░▌       ▐░▌▐░▌       ▐░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░▌  ",
		" ▀▀▀▀▀▀▀▀▀▀▀  ▀         ▀  ▀         ▀  ▀▀▀▀▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀▀▀▀   ",
	}

	result := ""
	for i, line := range lines {
		if i < 5 {
			result += topStyle.Render(line) + "\n"
		} else {
			result += bottomStyle.Render(line) + "\n"
		}
	}

	taglineStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Italic(true).
		Align(lipgloss.Center)

	result += "\n" + lipgloss.NewStyle().Width(70).Align(lipgloss.Center).
		Render(taglineStyle.Render("Secure Hedera Wallet"))

	return result
}
