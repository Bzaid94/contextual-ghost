# 👻 Contextual Ghost (CG)
> *The Survival Agent for Developers.*

![Ghost Demo](https://media3.giphy.com/media/v1.Y2lkPTc5MGI3NjExcWd1NWk1YXM4aXRlNDM2c3dheWVrbjl1cHBocnd5c2Rkc2FtdHVubyZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/pmGFexW39VZRCHbSS8/giphy.gif)

**Stop Copy-Pasting Errors.**
**Stop Losing Context.**
**Stop Breaking Flow.**

Contextual Ghost (CG) is not just a tool; it's your undead pair programmer. It watches your CLI commands from the shadows. When you succeed, it stays silent. When you fail, it manifests immediately with the solution.

Powered by **GitHub Copilot CLI**, Ghost doesn't just guess—it *knows*. It sees your Git state, your environment, and your logs to provide surgically accurate fixes before you even switch to the browser.

## 🚀 Why Ghost?
- **Zero Friction**: Acts as a transparent wrapper. `ghost npm run build` behaves exactly like `npm run build`... until it breaks.
- **Context Aware**: It knows you just changed `app.js` and that you're on `macos`. It feeds this to the AI for better answers.
- **Proactive**: No need to type "explain this error". Ghost is already finding the fix while the error log is still scrolling.

## 🛠 Tech Stack
- **Lang**: Go (Golang) for raw speed and native process handling.
- **Brain**: GitHub Copilot CLI (`gh copilot`).
- **UI**: CharmBracelet (`bubbletea` & `lipgloss`) for a terminal UI that feels like 2077.

## 📦 Installation

### Option 1: Binary (Recommended for non-Go users)
Download the latest binary for your OS (Windows, macOS, Linux) from the [Releases](https://github.com/Bzaid94/contextual-ghost/releases) page.

```bash
# Example for macOS/Linux:
curl -sL https://github.com/Bzaid94/contextual-ghost/releases/latest/download/contextual-ghost_Darwin_x86_64.tar.gz | tar xz
sudo mv contextual-ghost /usr/local/bin/
```

### Option 2: Go Install (For Go users)
```bash
go install github.com/Bzaid94/contextual-ghost@latest
```

## 🧠 Brains (Required)
Ghost requires the GitHub CLI and the Copilot extension to function:
```bash
# Install GitHub CLI
# macOS: brew install gh
# Windows: winget install Microsoft.GitHubCLI

# Install Copilot extension
gh extension install github/gh-copilot
```

## 🎮 Usage
Don't change your workflow. Just invite the Ghost.

```bash
# General usage:
ghost <your-command>

# Example:
ghost ls /nonexistent
```

If it works? **Silence.**
If it breaks? **Salvation.**

## 📜 License
This project is licensed under the **Apache License 2.0**. See the [LICENSE](LICENSE) file for details.

---
*Built for the GitHub Challenge. Crafted with ❤️ and 👻 by @Bzaid94.*
