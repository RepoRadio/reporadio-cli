# Installing RepoRadio: Step-by-Step Guide

Hello there, and welcome back to another episode of RepoRadio! I hope you’re all doing wonderfully today. If you're tuning in again after our last episode, you’re probably getting excited about the possibilities RepoRadio offers for transforming your documentation into an audio format. Today, we’re going to build on what we've covered by walking through the process of installing RepoRadio on your system. So, grab a cup of coffee, sit back, and let's get started!

## Getting Ready: Pre-Installation Checklist

First things first. Before we dive into the installation itself, there’s an important step to complete, which is obtaining your OpenAI API key. As some of you might remember, RepoRadio utilizes this key to power its audio generation capabilities. Don’t worry, it’s straightforward!

Head over to the OpenAI platform at platform.openai.com and sign up for a free account if you haven’t already. Once you’re in, navigate to the API keys section and generate a new key. Keep it safe, because we’ll be using it in just a moment.

Next, you'll want to set this API key as an environment variable on your machine. This step is crucial for RepoRadio to access the OpenAI API. Run the following command in your terminal:

```bash
export OPENAI_API_KEY=sk-...
```

Remember to replace `sk-...` with your actual API key. Now, we're ready to proceed to the installation.

## Installing RepoRadio via Go

RepoRadio works as a command-line interface tool, and installation is nice and easy with Go. If you haven’t installed Go yet, you’ll need to do so from the official Go website. Once you have Go set up, make sure your `$GOPATH/bin` is included in your system's `$PATH` so that RepoRadio can be run from anywhere in your terminal.

Now, you can install RepoRadio using the following command:

```bash
go install github.com/reporad-io/reporadio@latest
```

And voila! RepoRadio is now installed on your system. To verify, just type `reporadio` in your terminal, and if everything is set up correctly, you’ll see a list of commands you can use.

## Creating and Generating Podcasts

Okay, with RepoRadio installed, what’s next? Let’s test it out.

Start by creating a new podcast project. You can do that by running:

```bash
reporadio create my-podcast
```

This command sets up a new podcast setup with the name "my-podcast". Feel free to replace "my-podcast" with whatever name suits your series.

Once you have your podcast created, generating it is just as simple. Use this command:

```bash
reporadio generate my-podcast
```

And there you go, your very own audio documentation from your repository is ready for listening!

## Troubleshooting Tips

Now, let's address some common issues you might encounter during installation. One of the typical stumbling blocks is related to the PATH configuration. Double-check that your `$GOPATH/bin` is in your `$PATH`, as this path is crucial for Go to locate the RepoRadio binaries.

Another common hiccup involves the API key setup. If RepoRadio is having trouble accessing the OpenAI API, the environment variable might not be set correctly. Verify that you exported the key in the same terminal session, or consider adding the `export` command to your shell’s profile file for permanence.

## Wrapping Up

Well folks, there you have it! Installing RepoRadio and setting up your first podcast project can be a breeze if you follow these steps. As always, the goal of RepoRadio is to make documentation accessible in ways that suit different learning styles, whether that's reading or listening.

I hope this guide has been helpful and that you’re now ready to create your own audio documentation using RepoRadio. If you run into any trouble or have questions, don’t hesitate to reach out via our GitHub page or drop us an email. Your feedback is incredibly valuable as we continue to improve and expand RepoRadio.

Thanks for tuning in, and don't forget to check back for more insights and guides about making the most of your codebase documentation. Until next time, happy coding and happy listening!