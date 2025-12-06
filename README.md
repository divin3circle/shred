# Shred Wallet

**Shred** is a terminal-native Hedera wallet built with Go. Manage your Hedera accounts, send and receive HBAR, and interact with the Hedera network directly from your terminal.

## Features

- 🔐 **Secure Wallet Management**: BIP-39 mnemonic generation with Argon2id encryption
- 💰 **Hedera Integration**: Full support for Hedera accounts and HBAR transactions
- 🎨 **Beautiful TUI**: Terminal user interface built with Bubble Tea
- 🔄 **Multi-Wallet Support**: Create and manage multiple wallets
- 🚀 **Cross-Platform**: Works on macOS, Windows, and Linux
- 🔒 **Auto-Lock**: Automatic wallet locking after inactivity

## 📥 Download

Download the latest release from the [Releases](https://github.com/divin3circle/shred/releases) page.

Choose the correct file for your system:

| Operating System                   | File Name                             |
| :--------------------------------- | :------------------------------------ |
| **macOS (Apple Silicon M1/M2/M3)** | `shred_{version}_darwin_arm64.tar.gz` |
| **macOS (Intel)**                  | `shred_{version}_darwin_amd64.tar.gz` |
| **Windows**                        | `shred_{version}_windows_amd64.zip`   |
| **Linux (AMD64)**                  | `shred_{version}_linux_amd64.tar.gz`  |
| **Linux (ARM64)**                  | `shred_{version}_linux_arm64.tar.gz`  |

## 🐳 Docker

You can also run Shred using Docker:

```bash
docker pull yourusername/shred:latest
docker run -it yourusername/shred
```

## 🍎 macOS Installation

1. **Download** the appropriate file for your Mac (ARM64 or AMD64)
2. **Extract** the archive:
   ```bash
   tar -xzf shred_*_darwin_*.tar.gz
   ```
3. **Make it executable**:
   ```bash
   chmod +x shred
   ```
4. **Run the app**:
   ```bash
   ./shred
   ```

> **Security Warning?**  
> If you see "Unidentified Developer", go to **System Settings > Privacy & Security** and click **Open Anyway**. Alternatively, right-click the file in Finder and select **Open**.

## 🪟 Windows Installation

1. **Download** `shred_{version}_windows_amd64.zip`
2. **Extract** the ZIP file
3. **Open PowerShell** and navigate to the extracted folder
4. **Run the app**:
   ```powershell
   .\shred.exe
   ```

> **Windows Protected Your PC?**  
> If SmartScreen blocks the app, click **More Info** -> **Run Anyway**.

## 🐧 Linux Installation

1. **Download** the appropriate Linux archive
2. **Extract**:
   ```bash
   tar -xzf shred_*_linux_*.tar.gz
   ```
3. **Make it executable**:
   ```bash
   chmod +x shred
   ```
4. **Run the app**:
   ```bash
   ./shred
   ```

## 🚀 Getting Started

### First Time Setup

1. **Create New Wallet**: Press `N` to generate a new 24-word recovery phrase
2. **Save Your Phrase**: Write down the 24 words on paper. **Do not lose them.**
3. **Verify**: Re-enter specific words to verify you saved them correctly
4. **Set Password**: Create a strong password to encrypt your wallet file

### Using Your Wallet

- **Select Wallet**: On startup, choose from your existing wallets
- **Unlock**: Enter your passphrase to unlock your wallet
- **Dashboard**: View your account balance, EVM address, and account status
- **Refresh**: Press `f` to refresh account information

### Controls

- `↑/↓` or `j/k`: Navigate menus
- `Enter`: Select / Confirm
- `Esc`: Go Back
- `f`: Refresh account information
- `q`: Quit

## 🔧 Development

### Prerequisites

- Go 1.24.1 or later
- Git

### Build from Source

```bash
git clone https://github.com/divin3circle/shred.git
cd shred
go build ./cmd/shred
```

### Run Tests

```bash
go test ./...
```

## 📖 Documentation

- [User Guide](USER_GUIDE.md) - Detailed usage instructions
- [Deployment Guide](DEPLOYMENT.md) - Deployment and CI/CD information

## 🔒 Security

- Wallets are encrypted using Argon2id key derivation and AES-GCM encryption
- Mnemonics are never stored in plain text
- Auto-lock feature protects your wallet after inactivity
- Memory wiping on lock (best effort in Go)

## 📝 License

This project is licensed under the AGPL-3.0 License with additional commercial restrictions. See [LICENSE](LICENSE) for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## ⚠️ Disclaimer

This software is provided "as is" without warranty. Use at your own risk. Always verify your recovery phrase and keep it secure.

## 🔗 Links

- [GitHub Repository](https://github.com/divin3circle/shred)
- [Docker Hub](https://hub.docker.com/r/yourusername/shred)
