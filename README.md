# {📻} RepoRad.io CLI

**Turn your Git repository into a podcast.**
RepoRadio is an open source CLI tool that generates narrated audio episodes directly from your codebase.

It’s free forever, powered by your own OpenAI API key, and works entirely from the command line.

GitHub: [https://github.com/RepoRadio/reporadio-cli](https://github.com/RepoRadio/reporadio-cli)

---

## ✨ What It Does

* 🔍 Analyzes your repo (README, structure, metadata, commits)
* 🎙️ Generates narrated audio content:

  * Contributor onboarding episodes
  * Consumer-facing getting-started guides
  * Change log summaries
* 🧑‍💻 Built for developers who prefer audio over reading long docs

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
go install github.com/reporad-io/reporadio@latest
```

Make sure `$GOPATH/bin` is in your `$PATH`.

---

## 🚀 Usage Examples

### Create a new podcast:

```bash
reporadio create my-podcast
```

### Generate a podcast:

```bash
reporadio generate my-podcast
```
---

## 🧪 What It's Great For

* Solo developers switching between repos
* Open source maintainers offering better onboarding
* Consultants moving across unfamiliar codebases
* Anyone who learns better by listening than reading

---

## 🐶 RepoRadio on RepoRadio

We use RepoRadio to explain itself.
🎧 [Listen to the onboarding podcast for this repo](https://reporad.io/onboarding)

---

## 🧠 Philosophy

> Docs should adapt to how people learn—not the other way around.

RepoRadio helps you create and consume spoken-word documentation, powered by the tools you already use.

No hosted accounts. No lock-in. Just a CLI and your API key.

---

## 🤝 Contribute

We welcome contributions, feedback, and new use cases!

* Star this repo ⭐
* Open issues 🐛
* Submit pull requests 🛠️

---

## 📮 Contact

* GitHub Issues: [https://github.com/RepoRadio/reporadio-cli/issues](https://github.com/RepoRadio/reporadio-cli/issues)
* Email: `hello@reporad.io`

---

> Built with Go. Always free. Powered by your own OpenAI key.
