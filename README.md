<div align="center">
  <a href="https://emby.media/"><img src="https://img.shields.io/badge/Emby-52b54b?logo=emby&logoColor=white"/></a>
  <a href="https://github.com/computer-geek64/emboxd/releases/latest"><img src="https://img.shields.io/github/v/release/computer-geek64/emboxd"/></a>
  <a href="https://github.com/search?q=repo%3Acomputer-geek64%2Femboxd++language%3AGo&type=code"><img src="https://img.shields.io/github/languages/top/computer-geek64/emboxd"/></a>
  <a href="https://github.com/computer-geek64/emboxd/issues?q=is%3Aissue%20state%3Aopen%20label%3Abug"><img src="https://img.shields.io/github/issues/computer-geek64/emboxd/bug"/></a>
  <a href="LICENSE"><img src="https://img.shields.io/github/license/computer-geek64/emboxd"/></a>
  <a href="https://github.com/computer-geek64/emboxd/forks"><img src="https://img.shields.io/github/forks/computer-geek64/emboxd"/></a>
  <a href="https://github.com/computer-geek64/emboxd/stargazers"><img src="https://img.shields.io/github/stars/computer-geek64/emboxd"/></a>

  <h1>EmBoxd</h1>

  <h4>Live sync server for Letterboxd users with self-hosted media platforms</h4>
</div>


## Table of Contents

- [About](#about)
- [Installation](#installation)
  - [Binary](#installation)
  - [Docker](#docker)
- [Usage](#usage)
  - [Configuration](#configuration)
  - [Running](#running)
- [Contributors](#contributors)
- [License](#license)


## About

EmBoxd provides live integration with Letterboxd for users of self-hosted media servers.
It tracks watch activity on the media server and synchronizes Letterboxd user data to match.
Changes to a movie's played status are reflected in the user's watched films, and movies that are fully played are logged in the user's diary.

The following media servers are currently supported or have planned support:

- [X] Emby
- [ ] Jellyfin [#4](https://github.com/computer-geek64/emboxd/issues/4)
- [ ] Plex [#6](https://github.com/computer-geek64/emboxd/issues/6)


## Installation

EmBoxd can either be setup and used as a binary or Docker image

### Binary

Building a binary from source requires the Go runtime

1. Clone repository:

```sh
git clone https://github.com/computer-geek64/emboxd.git --depth=1
cd emboxd/
```

2. Install Playwright browsers and OS dependencies:

```sh
go install github.com/playwright-community/playwright-go/cmd/playwright
playwright install --with-deps
```

3. Build and install binary (to GOPATH)

```sh
go install .
```

### Docker

Pull from GitHub container registry:

```sh
docker pull ghcr.io/computer-geek64/emboxd:latest
```

Or build image from source:

```sh
git clone https://github.com/computer-geek64/emboxd.git --depth=1
docker build -t emboxd:latest emboxd/
```

## Usage

### Configuration

The YAML configuration file describes how to link Letterboxd accounts with media server users.
The format should follow the example [`config.yaml`](config.yaml) in the repository root.

Supported media servers need to send webhook notifications for all (relevant) users to the EmBoxd server API.

Emby should send the following notifications to `/emby/webhook`:

- [X] Playback
  - [X] Start
  - [X] Pause
  - [X] Unpause
  - [X] Stop
- [X] Users
  - [X] Mark Played
  - [X] Mark Unplayed

### Running

Running EmBoxd starts the server and binds with port 80.
The `-c`/`--config` option specifies the config file to use with the server.

When running with Docker, the image expects the configuration file at `/config/config.yaml`.
It can be bind-mounted to the container or stored in a volume.

```sh
docker run --name=emboxd --restart=unless-stopped -v config.yaml:/config/config.yaml:ro -p 80:80 ghcr.io/computer-geek64/emboxd:latest
```


## Contributors

<a href="https://github.com/computer-geek64/emboxd/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=computer-geek64/emboxd"/>
</a>


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
