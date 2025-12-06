# Shred Wallet - User Guide

Welcome to **Shred**, the terminal-native Hedera wallet. This guide will help you download, install, and run Shred on your computer.

## 📥 Download

Download the latest release from the [GitHub Releases](https://github.com/divin3circle/shred/releases) page.

Choose the correct file for your system:

| Operating System                   | File Name                             |
| :--------------------------------- | :------------------------------------ |
| **macOS (Apple Silicon M1/M2/M3)** | `shred_{version}_darwin_arm64.tar.gz` |
| **macOS (Intel)**                  | `shred_{version}_darwin_amd64.tar.gz` |
| **Windows**                        | `shred_{version}_windows_amd64.zip`   |
| **Linux (AMD64)**                  | `shred_{version}_linux_amd64.tar.gz`  |
| **Linux (ARM64)**                  | `shred_{version}_linux_arm64.tar.gz`  |

**Note:** After extracting the archive, you'll find a binary named `shred` (or `shred.exe` on Windows).

---

## 🍎 macOS Installation

1.  **Download** the appropriate file for your Mac:

    - Apple Silicon (M1/M2/M3): `shred_{version}_darwin_arm64.tar.gz`
    - Intel: `shred_{version}_darwin_amd64.tar.gz`

2.  **Extract the archive**: Double-click the downloaded `.tar.gz` file, or in Terminal:

    ```bash
    cd ~/Downloads
    tar -xzf shred_*_darwin_*.tar.gz
    ```

3.  **Open Terminal**: Press `Cmd + Space`, type `Terminal`, and press Enter.

4.  **Navigate to the extracted folder**:

    ```bash
    cd ~/Downloads
    ```

5.  **Make it executable**:

    ```bash
    chmod +x shred
    ```

6.  **Run the app**:

    ```bash
    ./shred
    ```

    **Optional:** Move to a location in your PATH for easier access:

```bash
sudo mv shred /usr/local/bin/
```

Then you can run `shred` from anywhere.

> **Security Warning?**
> If you see "Unidentified Developer", go to **System Settings > Privacy & Security** and click **Open Anyway**. Alternatively, right-click the file in Finder and select **Open**.

---

## 🪟 Windows Installation

1.  **Download** `shred_{version}_windows_amd64.zip`

2.  **Extract the file**: Right-click the `.zip` file and select **Extract All**.

3.  **Open PowerShell**: Press `Win + X` and select **Terminal** or **PowerShell**.

4.  **Navigate to the extracted folder**:

    ```powershell
    cd Downloads
    cd shred_*_windows_amd64
    ```

5.  **Run the app**:
    ```powershell
    .\shred.exe
    ```

> **Windows Protected Your PC?**
> If SmartScreen blocks the app, click **More Info** -> **Run Anyway**.

---

## 🐧 Linux Installation

1.  **Download** the appropriate file:

    - AMD64: `shred_{version}_linux_amd64.tar.gz`
    - ARM64: `shred_{version}_linux_arm64.tar.gz`

2.  **Extract the archive**:

    ```bash
    tar -xzf shred_*_linux_*.tar.gz
    ```

3.  **Make it executable**:

    ```bash
    chmod +x shred
    ```

4.  **Run the app**:

    ```bash
    ./shred
    ```

    **Optional:** Install system-wide:

```bash
sudo mv shred /usr/local/bin/
```

Then you can run `shred` from anywhere.

---

## 🚀 Getting Started

Once the application is running:

1.  **Create New Wallet**: Press `N` to generate a new 24-word recovery phrase.
2.  **Save Your Phrase**: Write down the 24 words on paper. **Do not lose them.**
3.  **Verify**: Re-enter specific words to verify you saved them correctly.
4.  **Set Password**: Create a strong password to encrypt your wallet file.
5.  **Activate Account**: Hedera accounts after creation need to be funded by a small amount of HBAR < 1HBAR for the account to be activated. After creation submit an issue [here](https://github.com/divin3circle/shred/issues/new/choose) to receive some HBAR or get them from [https://portal.hedera.com/faucet](this faucet)

**Controls:**

- `↑/↓` or `j/k`: Navigate menus
- `Enter`: Select / Confirm
- `Esc`: Go Back
- `q`: Quit
