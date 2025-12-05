# Shred Wallet - User Guide

Welcome to **Shred**, the terminal-native Hedera wallet. This guide will help you download, install, and run Shred on your computer.

## 📥 Download

Visit our official website to download the latest version for your operating system:
**[https://your-website.com/download](https://your-website.com/download)**

Choose the correct file for your system:

| Operating System | File Name |
| :--- | :--- |
| **macOS (Apple Silicon M1/M2/M3)** | `shred-mac-arm64.zip` |
| **macOS (Intel)** | `shred-mac-intel.zip` |
| **Windows** | `shred-windows.zip` |
| **Linux** | `shred-linux-amd64.zip` |

---

## 🍎 macOS Installation

1.  **Extract the file**: Double-click the downloaded `.zip` file.
2.  **Open Terminal**: Press `Cmd + Space`, type `Terminal`, and press Enter.
3.  **Navigate to the folder**:
    ```bash
    cd ~/Downloads
    ```
4.  **Make it executable** (only needed once):
    ```bash
    chmod +x shred-mac-arm64  # or shred-mac-intel
    ```
5.  **Run the app**:
    ```bash
    ./shred-mac-arm64
    ```

> **Security Warning?**
> If you see "Unidentified Developer", go to **System Settings > Privacy & Security** and click **Open Anyway**. Alternatively, right-click the file in Finder and select **Open**.

---

## 🪟 Windows Installation

1.  **Extract the file**: Right-click the `.zip` file and select **Extract All**.
2.  **Open PowerShell**: Press `Win + X` and select **Terminal** or **PowerShell**.
3.  **Navigate to the folder**:
    ```powershell
    cd Downloads\shred-windows
    ```
4.  **Run the app**:
    ```powershell
    .\shred-windows.exe
    ```

> **Windows Protected Your PC?**
> If SmartScreen blocks the app, click **More Info** -> **Run Anyway**.

---

## 🐧 Linux Installation

1.  **Extract the file**:
    ```bash
    unzip shred-linux-amd64.zip
    ```
2.  **Make it executable**:
    ```bash
    chmod +x shred-linux-amd64
    ```
3.  **Run the app**:
    ```bash
    ./shred-linux-amd64
    ```

---

## 🚀 Getting Started

Once the application is running:

1.  **Create New Wallet**: Press `N` to generate a new 24-word recovery phrase.
2.  **Save Your Phrase**: Write down the 24 words on paper. **Do not lose them.**
3.  **Verify**: Re-enter specific words to verify you saved them correctly.
4.  **Set Password**: Create a strong password to encrypt your wallet file.

**Controls:**
*   `↑/↓` or `j/k`: Navigate menus
*   `Enter`: Select / Confirm
*   `Esc`: Go Back
*   `q`: Quit
