# Setting Up Your Development Environment

**[Intro Music Fades Out]**

**Host:** Hey there, listeners! Welcome back to RepoRadio. If you joined us in our previous episode, you'll remember how we set the stage with an overview of RepoRadio's mission and features. Today, we're rolling up our sleeves and diving into the essential steps to get your development environment ready, so you can start transforming your Git repositories into captivating podcasts. I’m thrilled to guide you through this process, so let's get started!

Setting up your environment might sound daunting, but I promise you'll find it straightforward and rewarding. This episode is packed with actionable steps, so you might want to keep your terminal open and follow along. Ready? Let’s dive in!

**[Sound Effect: Keyboard Typing]**

## Step 1: Get Your OpenAI API Key

First up, you need to secure your OpenAI API key. Head over to [OpenAI's platform](https://platform.openai.com/account/api-keys) and grab a free API key if you haven’t already. This key is crucial—it acts as your portal to the AI magic that powers RepoRadio’s audio generation.

Once you've got your key, you'll want to set it in your environment. Let’s make this a breeze with a quick command:

```bash
export OPENAI_API_KEY=sk-your-api-key
```

This is your ticket to unleashing the capabilities of OpenAI within RepoRadio. Imagine it as the smart engine behind generating those delightful audio narratives!

## Step 2: Install RepoRadio via Go

With your API key in place, the next step is installing RepoRadio itself. RepoRadio is developed using Go, a language loved for its simplicity and efficiency. If you haven't installed Go yet, you might want to pause and take care of that now. Once you're set, install the latest version of RepoRadio with this command:

```bash
go install github.com/reporad-io/reporadio@latest
```

To ensure seamless operation, make sure that your Go binary path, usually found at `$GOPATH/bin`, is in your system’s `$PATH`. This ensures your terminal knows exactly how to execute `reporadio` commands whenever you call them.

**[Sound Effect: Successful Ding]**

## Step 3: Create and Generate Your Podcast

Now we get to the fun and creative part! With RepoRadio installed, it’s time to create your first podcast:

```bash
reporadio create my-podcast
```

This command initiates everything you need. Once you’ve initialized your podcast directory, it’s time to pull the magic switch—generating the podcast itself:

```bash
reporadio generate my-podcast
```

Here is where RepoRadio works its charm, crafting those spoken-word narratives from your repository's readme, structure, and even commit history. Whether you're looking to create an engaging change log summary, a contributor guide, or a user manual, RepoRadio is your go-to tool!

## Why You’ll Love RepoRadio

If you're juggling multiple projects or diving into new codebases like a consultant, RepoRadio brings audio documentation to your workspace, making the experience intuitive and engaging. It’s especially great for open-source maintainers keen on enhancing community interactions or for anyone who finds listening more enriching than reading dense pages of text.

Not to mention, RepoRadio is free, open-source, and gives you complete control without the shackles of subscriptions—just your expertise and an API key!

**[Sound Effect: Applause]**

And that’s a wrap on today’s episode! You’re now equipped to embark on podcasting your codebase using RepoRadio. Remember, you’re not just creating documentation; you’re crafting a narrative that tells your project’s story.

If you’re hungry for more guidance, why not delve back into some of our previous episodes, explore our README file, or even join the community on GitHub to contribute and learn?

**[Outro Music Fades In]**

**Host:** Thanks for joining me today on RepoRadio. We hope you enjoyed this episode and look forward to the incredible podcasts you’ll create. If this episode struck a chord, consider subscribing and leaving a review. Connect with us through GitHub or shoot us an email with any thoughts. Until next time, keep coding and happy podcasting!

**[Outro Music Fades Out]**