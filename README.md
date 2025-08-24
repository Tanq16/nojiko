<div align="center">

<img src=".github/assets/logo-mid.png" alt="Nojiko Logo" width="250"/>

<h1>Nojiko Self-Hosted Dashboard</h1>

[![Release Build](https://github.com/tanq16/nojiko/actions/workflows/release.yml/badge.svg)](https://github.com/tanq16/nojiko/actions/workflows/release.yml)
[![GitHub Release](https://img.shields.io/github/v/release/tanq16/nojiko)](https://github.com/Tanq16/nojiko/releases/latest)

A simple, beautiful, and customizable dashboard for your self-hosted domain. Nojiko focuses on simple status and feed trackers, with a comprehensive bookmarks system to be used as the primary webpage across all devices.

<a href="#features">Features</a> &bull; <a href="#installation">Installation</a> &bull; <a href="#usage">Usage</a> &bull; <a href="ai-development">AI-Development</a>

</div>

---

`Quickstart` &rarr;

```bash
docker run --rm -v "$HOME/nojiko/config.yaml:/app/config.yaml" -p 80:8080 --name "nojiko" -d tanq16/nojiko:main
# $HOME/nojiko should contain the config.yaml file
```

## Features

Nojiko focuses on simplicity in a dashboard application. It offers 4 main sections on the screen:
- Header - this shows a logo (optional) and title of your choice, along with optional weather given your location
- Bookmarks - this is a sidebar that shows your bookmarks with your choice of names and icons across your choice of categories and folders within categories (both foldable)
- Status Cards - these are the first set of cards on the main screen, currently 2 types of sections are supported - github repo status (shows PRs, issues, and stars), and service status (stats from other homelab services)
- Feed Cards - this is the next set of cards that shows youtube videos with thumbnails for your configured creators

> [!WARNING]
> Currently, only Adguard Home and Jellyfin are supported for services. Additionally, only Youtube feeds are supported. RSS feeds and additional services will be added in the future.

> [!IMPORTANT]
> All configuration of the dashboard is done via a `config.yaml` file. You can choose your icons, but they need to be part of [lucide](https://lucide.dev/icons/). The config file can also be edited directly from the dashboard page.

> [!TIP]
> When attempting to test out configurations, it's important to note that things like GitHub status and YouTube feeds can often lead to throttling. Therefore, it's best to try out configs with such sections commented out.

Here's a screenshot of what the app looks like:

![](.github/assets/ss.png)

## Installation

Create a directory to store the configuration file:

```bash
mkdir $HOME/nojiko
```

Create a config.yaml file similar to the `config.yaml.example` provided in the repo. Add the necessary details (read the Usage section below for more). With the config ready, use the following compose yaml to deploy the app:

```yaml
services:
  nojiko:
    image: tanq16/nojiko:main
    container_name: nojiko
    volumes:
      - /home/USER/nojiko/config.yaml:/app/config.yaml # this is the dir that has your config
    ports:
      - 8080:8080
```

The application with run at your `machineIP:8080`. Allow a few seconds for the first data fetch to complete before the server is enabled.

## Usage

Nojiko is configured entirely through the `config.yaml` file. The dashboard's content, from the header to the bookmarks, is determined by this file. You can edit this file directly or through the settings icon on the dashboard itself.

Below is a breakdown of the main configuration sections.

### Header

This section controls the main header at the top of the page.

```yaml
header:
  title: "Nojiko Dashboard"
  showLogo: false
  logoURL: "" # URL to your logo image
  showWeather: true
  latitude: 43.65
  longitude: -79.38
```

- **title**: The main title displayed in the header.
- **showLogo**: Set to `true` to display a logo next to the title.
- **logoURL**: The direct URL to your logo image.
- **showWeather**: Set to `true` to display the current weather.
- **latitude** & **longitude**: Your location coordinates for accurate weather data.

### Bookmarks

This section configures the bookmarks displayed in the sidebar. You can organize links into categories and folders.

```yaml
bookmarks:
  - category: "Social Media"
    color: "ctp-green" # Uses Catppuccin color names
    links:
      - name: "Twitter / X"
        url: "#"
        icon: "twitter" # Icon name from lucide.dev
  - category: "Development"
    color: "ctp-yellow"
    folded: true # Category is collapsed by default
    folders:
      - name: "Frameworks"
        icon: "folder"
        folded: false # Folder is expanded by default
        links:
          - name: "React Docs"
            url: "#"
            icon: "dot"
```

- **category**: The name of the bookmark group.
- **color**: Sets the color for the category title. Use Catppuccin color names (e.g., `ctp-green`, `ctp-mauve`).
- **links**: A list of direct bookmarks.
- **folders**: A list of folders within a category to further organize links.
- **folded**: Set to `true` to make a category or folder collapsed by default.
- **icon**: The name of an icon from [lucide.dev](https://lucide.dev/icons/) to display next to the link or folder.

### Status Cards

These cards display status information from various services, like GitHub repository statistics or the health of your self-hosted applications.

```yaml
statusCards:
  - title: "Project Status"
    icon: "github"
    type: "github"
    cards:
      - repo: "tanq16/nojiko"
  - title: "Service Status"
    icon: "bar-chart-3"
    type: "service"
    cards:
      - name: "Adguard Home"
        serviceType: "adguard"
        url: "http://adguard.local"
        username: "your_user"
        password: "your_password"
      - name: "Jellyfin"
        serviceType: "jellyfin"
        url: "http://jellyfin.local"
        apiKey: "your_api_key"
```

- **title**: The title for the status card section.
- **icon**: An icon for the section title from [lucide.dev](https://lucide.dev/icons/).
- **type**: The type of cards in this section. Supported types are `github` and `service`.
- **cards**: A list of items to display.
  - For `github`, specify the `repo` as `owner/repository`.
  - For `service`, specify the `name`, `serviceType` (`adguard` or `jellyfin`), `url`, and necessary credentials (`username`/`password` for Adguard, `apiKey` for Jellyfin).

### Thumbnail Feeds

This section displays a feed of cards with thumbnails, currently supporting YouTube channels.

```yaml
thumbFeeds:
  - title: "YouTube Feed"
    icon: "youtube"
    feedType: "youtube"
    limit: 10
    channels:
      - "mkbhd"
      - "LinusTechTips"
```

- **title**: The title for the feed section.
- **icon**: An icon for the section title from [lucide.dev](https://lucide.dev/icons/).
- **feedType**: The type of feed. Currently, only `youtube` is supported.
- **limit**: The maximum number of videos to display.
- **channels**: A list of YouTube channel handles (the `@` is not needed).

## AI-Development

This is my first project that I completely vibe-coded (and am continuing to do so). Vibe coded also means various things - the description I mean is the fact that I control implementation and tell AI what to do, not a simplistic variant of "hey AI, write me a dashboard application that has bookmarks". I described the entirety of the vibe code process I took for this project below, so feel free to expand and read.

<details>
<summary>Previous experience with vibe-coding vs. this one</summary>
Usually, I use AI autocomplete for quick writing and also to quicly build and iterate on UI configurations. I've been pretty successful at increasing productivity quite a bit this way. However, I've never had good luck with fully vibe coding (i.e., minimal intervention into code) the self-hosted apps I was writing - AI used to either overcomplicate, misunderstand, or competely botch the implementation.

So what's different with this one? After a bunch of experimentation, I finally landed on proper instructions to prompt the model to restrict its creativity. By that, I mean, I explicitly stated each and every thing I wanted to be present in the code including formats and structs, thereby leaving little room for the model to deviate from what I want. Additionally, I explicitly stated a large number of rules for it to follow for the implementation, which I had previously seen it go off on its own. Example, I've seen the model use arbitrary third party libraries for various things, add an extraneous amounts of comments unnecessarily, and implement it as if it were the backend implementation of Google, not a simple self-hosted app I wanted to write and use.

The result was actually really good! And it became a very quick way to iterate on the implementation. Design was always fast to iterate on anyway (hint: AI models are actually really good with UI stuff). It took a LOT of prompting. To quantify, around a total of around 6-7 hours of just writing English text describing everything I wanted in code and design. BUT, by restricting the model to do exactly what I wanted and leave little wiggle room, I was actually able to get the result I wanted, which I would have gotten to myself if I was actively writing code.
</details>

<details>
<summary>Vibe Code Process for this Project</summary>
Broadly speaking, I can divide the process I followed into 3 phases in order:

- Frontend Design: Why frontend first? It's hard for me to do frontend work, writing backend, specially for simple apps is much easier significantly. so even if I have no aid on the backend, it's totally manageable and fast as long as I quickly get a UI rolling.
- Backend Design: This is where I can completely imagine and write whatever is needed. So describing stuff is really easy in plain english. The only thing was to get the model to do it "my way".
- Iterations: Once the app is in place and the model does things the way I want, iterations to add features or fix things become trivial.

One quick cheatcode I discovered about the UI. Using a pre-existing color palette along with Tailwind and a prompt asking the AI to not use any CSS that needs to be managed, actually yields incredible results. After that it's just a question of tweaking.

My overall process step-wise was this:
- I had a general idea of the design I wanted. So I created a mockup on my iPad. By mockup, I mean hand-drawn boxes that I ideated on. I just gave that straight to the model. This took some time to get right, but eventually got there.
- Next thing I did (very important), looked at the frontend pieces to remove ambiguity and unnecessary complications. This required reading code and giving specifics to the model.
- With a plain single HTML file ready, I was free to ideate on the backend. I wrote down a large amount of prompt text indicating basic structs and interfaces I wanted to use, including the libraries I wanted to use and the overall idea I wanted to implement. As an example, this included the basic design and idea behind the state.go file, which I wanted to be the logic for maintaining a state that the frontend would receive and display.
- With the mock implementations ready, I was able to go deeper down into specifics and prompt my way through individual implementations like getting statuses from GitHub repos.
- Certain things became extremely easy due to AI, the primary one being fetching information for YouTube videos, which would have taken a while for me to do myself.
- At this point, iterations became simple with shorter prompts to add things. I, however, still followed the initial method of laying down the exact implementation of things I wanted to ensure the model doesn't do random things.

Overall, I also maintained a strategy to do two things:
- Every single primary prompt would be a large and structured piece of text aiming to do anywhere between 3-4 different things; and each prompt would represent one single chat only.
- Before switching to the new chat for the next big prompt, I also did a quick once over across the code base to ensure everything follows my vision and style.

With all this, honestly, I was able to get to a working sample with <24 hours of monitored vibe-coding. With only frontend being assisted by AI, I would expect this time to jump to over 1 week. And if I was to do it all myself, around 1 month LOL. Super happy with this fully vibe-coded project.
</details>
