# {📻} RepoRad.io CLI

**Turn your Git repository into a podcast.**
RepoRadio is an open source CLI tool that generates narrated audio episodes directly from your codebase.

It’s free forever, powered by your own OpenAI API key, and works entirely from the command line.

GitHub: [https://github.com/RepoRadio/reporadio-cli](https://github.com/RepoRadio/reporadio-cli)

## ✨ What It Does

## 🤖 Agent-Assisted Setup

You can use the GitHub Copilot agent or another agentic coder to set up RepoRadio CLI automatically.

**Quickstart:**

1. Download the setup instructions:
   ```bash
   curl -O https://raw.githubusercontent.com/RepoRadio/reporadio-cli/refs/heads/main/LLMs.txt
   ```
2. Give the `LLMs.txt` file to your agent (e.g., Copilot agent) and prompt:
   > Follow these instructions to setup reporadio-cli

The agent will guide you through the setup process interactively.

---

- 🔍 Analyzes your repo (README, structure, metadata, commits)
- 🎙️ Generates narrated audio content:
  - Contributor onboarding episodes
  - Consumer-facing getting-started guides
  - Change log summaries

- 🧑‍💻 Built for developers who prefer audio over reading long docs

---

## 💾 Installation

### ✅ Prerequisite

You’ll need a [free OpenAI API key](https://platform.openai.com/account/api-keys).
Set it as an environment variable:

```bash
export OPENAI_API_KEY=sk-...
```

### 📦 Install via Go

```bash
go install github.com/reporadio/reporadio-cli@main
```

Make sure `$GOPATH/bin` is in your `$PATH`.

---

## 🚀 Usage Examples

### Initialize template customization:

```bash
reporadio-cli init
```

This creates a `.reporadio/prompts/` directory with default templates that you can customize.

### Create a new podcast:

```bash
reporadio-cli create my-podcast
```

### Generate a podcast:

```bash
reporadio-cli generate my-podcast --audio
```

---

## 🎨 Template Customization

RepoRadio allows you to customize the AI prompts used for content generation:

### 1. Initialize custom templates:

```bash
reporadio-cli init
```

This creates a `.reporadio/prompts/` directory containing the default template files.

### 2. Customize templates:

Edit any `.tmpl` file in `.reporadio/prompts/` to customize how content is generated:

- `system_prompt.tmpl` - Main system instructions for the AI
- `episode_transcript.tmpl` - Template for episode content
- `extract_episodes.tmpl` - Template for extracting episodes from content
- And more...

### 3. Automatic usage:

When you run `create` or `generate` commands, RepoRadio will:
- Use your custom templates from `.reporadio/prompts/` when they exist
- Fall back to built-in templates for any missing custom files
- Ensure full compatibility with existing workflows

This allows you to fine-tune the AI's behavior while maintaining backward compatibility.

---

## 🧪 What It's Great For

- Solo developers switching between repos
- Open source maintainers offering better onboarding
- Consultants moving across unfamiliar codebases
- Anyone who learns better by listening than reading

---

## 🐶 RepoRadio on RepoRadio

We use RepoRadio to explain itself:

- 🎧 **Welcome podcast**: [`.reporadio/welcome/episodes/`](./.reporadio/welcome/episodes/) - Getting started guide
- 🎧 **Contributor podcast**: [`.reporadio/contributor/episodes/`](./.reporadio/contributor/episodes/) - Development onboarding

Each directory contains MP3 files and transcripts.

---

## 🧠 Philosophy

> Docs should adapt to how people learn—not the other way around.

RepoRadio helps you create and consume spoken-word documentation, powered by the tools you already use.

No hosted accounts. No lock-in. Just a CLI and your API key.

---

## 🤝 Contribute

We welcome contributions, feedback, and new use cases!

- Star this repo ⭐
- Open issues 🐛
- Submit pull requests 🛠️

---

## 📮 Contact

- GitHub Issues: [https://github.com/RepoRadio/reporadio-cli/issues](https://github.com/RepoRadio/reporadio-cli/issues)
- Email: `hello@reporad.io`

---

> Built with Go. Always free. Powered by your own OpenAI key.
