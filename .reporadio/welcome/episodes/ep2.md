# Installing RepoRadio: Step-by-Step Guide

**[Intro Music Fades Out]**

**Host:**  
Hello, everyone, and welcome back to another episode of RepoRadio, the podcast where we make codebases sing—quite literally! I'm your host, Chris, and today we're diving into the nitty-gritty details of getting RepoRadio installed on your system. If you love working with code and prefer learning by listening rather than reading, you're in the right place. So, get comfortable, and let’s get RepoRadio up and running on your machine!

## First Things First: What Is RepoRadio?

RepoRadio is a fantastic open-source command line tool that turns your Git repository into a podcast. Yes, you heard that right! It takes the key elements of your repo—like readmes, metadata, and commit histories—and generates narrated audio episodes. Whether you’re a developer switching between repos or an open-source maintainer trying to improve onboarding, RepoRadio has got you covered, making it easier to catch up on complex codebases, earbud-style.

## Prerequisites: The OpenAI API Key

Now, before we jump into the installation, there’s a crucial step we need to handle first: obtaining your OpenAI API key. RepoRadio is powered by OpenAI, and you’ll need your own API key to unlock all its magical, audio-generating capabilities.

Here’s how to get it:
1. Head over to the [OpenAI API Keys page](https://platform.openai.com/account/api-keys).
2. Register for a free account if you haven’t already.
3. Once you’re logged in, generate your API key.

With your API key ready, let’s set it as an environment variable. Open your terminal and type:

```bash
export OPENAI_API_KEY=sk-...
```

Remember to replace `sk-...` with your actual key!

## Installing RepoRadio via Go

With your API key set, it's time to install RepoRadio. This part's easy! RepoRadio is written in Go, so you’ll install it using `go install`. Enter the following command in your terminal:

```bash
go install github.com/reporadio/reporadio-cli@main
```

Make sure that `$GOPATH/bin` is in your `$PATH` so you can run RepoRadio from anywhere on your system. If you're not sure how to do this, check out some documentation on setting up your Go environment.

## Ready to Go!

Congratulations! You’ve now set up RepoRadio. You’re just moments away from generating your very first podcast from a codebase. Here’s how you start:

To create a new podcast, just type:

```bash
reporadio-cli create my-podcast
```

Then, to generate the audio content, use:

```bash
reporadio-cli generate my-podcast --audio
```

And there you have it! Your repository is speaking your language—literally!

## Troubleshooting Tips

Before we wrap up, let’s cover common troubleshooting issues you might encounter:
- **Issue: Command not found.** Double-check that `$GOPATH/bin` is indeed in your `$PATH`.
- **Issue: API Key error.** Ensure your API key is correctly set and that it has the right permissions on OpenAI’s platform.
- **Issue: Network problems.** Sometimes issues can arise due to internet connectivity, so make sure you’re connected.

## Get Involved

RepoRadio is a living project with a thriving community. We welcome your contributions, feedback, and new uses for the tool. Head over to our GitHub page, star the repo, or submit issues if you encounter any bugs.

And remember, we built this with a simple philosophy: Docs should adapt to how people learn, not the other way around!

**[Outro Music Fades In]**

**Host:**  
That’s all for today’s episode! Feel free to reach out at `hello@reporad.io` if you have questions or just want to share your experiences with us. Until next time, happy coding and happy listening! 

**[Outro Music Fades Out]**