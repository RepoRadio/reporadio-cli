# Setting Up Your Development Environment

Hello, everyone, and welcome back to another episode of RepoRadio! Today, we're diving into a crucial starting point for anyone looking to contribute to RepoRadio: setting up your development environment. Whether you're a seasoned coder or new to the world of open source, I've got you covered. So let's roll up our sleeves and get started!

First things first, before you can contribute, you'll need to get your system ready. This involves installing some necessary tools and gaining access to a key component—your OpenAI API key. Let’s break it down step by step.

## Why Set Up Your Environment?

RepoRadio is a unique open-source CLI tool that transforms your Git repository into podcast episodes. Not only does this mean you'll be working with inventive tech, but it also allows you to create engaging narrated audio content directly from your codebase. Imagine having contributor onboarding episodes or change log summaries narrated to you—that's the power of RepoRadio!

## Prerequisites

Before we jump into installations, you need to have an OpenAI API key. Don’t worry, it's free! Just head over to [OpenAI's website](https://platform.openai.com/account/api-keys) to grab your key. Once you have it, you’ll need to configure it as an environment variable in your terminal. Here’s the command you’ll use:

```bash
export OPENAI_API_KEY=sk-...
```

Just replace the "..." with your actual API key, and you’re good to go!

## Installing RepoRadio

Now that we've got our key, let’s proceed with the installation. RepoRadio is built with Go, so make sure you have Go installed on your system. Then, run the following command to install RepoRadio via Go:

```bash
go install github.com/reporadio/reporadio-cli@main
```

After installation, confirm that `$GOPATH/bin` is included in your `$PATH`. This ensures you can freely execute RepoRadio from anywhere in your terminal.

## How to Use RepoRadio

Great! If you've followed along, your environment is set up to contribute to RepoRadio. Let’s explore a couple of usage examples, so you know how to get started with generating podcasts:

- **Creating a new podcast**: Use the command:
  ```bash
  reporadio-cli create my-podcast
  ```
  Replace "my-podcast" with whatever name you give your podcast.

- **Generating a podcast**: Once you've created it, you can generate the audio content using:
  ```bash
  reporadio-cli generate my-podcast --audio
  ```

And there you have it! With these steps, you're ready to start creating engaging audio content from your repositories.

## Why You'll Love RepoRadio

RepoRadio isn’t just a tool; it’s a philosophy. It’s about adapting documentation to how people learn—by listening, not just reading. Whether you’re a developer switching repos, an open source maintainer enhancing onboarding, or a consultant getting to grips with unfamiliar codebases, RepoRadio has your needs covered. Plus, it’s totally free, built with Go, and powered by your OpenAI key. No hosted accounts, no lock-in—just the freedom to innovate.

## Get Involved

We’re always thrilled to welcome new contributors! Whether you want to star the repo, open issues, or submit pull requests, we’re eager to hear your feedback and see how you might improve RepoRadio. Every little bit helps keep this project thriving.

And with that, we've reached the end of today’s episode. Thanks for tuning in, and happy coding! Be sure to join us next time as we continue exploring the possibilities with RepoRadio. Until then, keep transforming those codebases into captivating audio experiences! 🎧

---

For any questions, thoughts, or if you just want to say hi, feel free to reach out on GitHub or email us at `hello@reporad.io`. 

Stay curious and keep learning!