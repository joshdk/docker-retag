[![License](https://img.shields.io/github/license/joshdk/docker-retag.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/joshdk/docker-retag)](https://goreportcard.com/report/github.com/joshdk/docker-retag)
[![CircleCI](https://circleci.com/gh/joshdk/docker-retag.svg?&style=shield)](https://circleci.com/gh/joshdk/docker-retag/tree/master)

# Docker Retag

🐳 Retag an existing Docker image without the overhead of pulling and pushing

## Motivation

There are certain situation where it is desirable to give an existing Docker image an additional tag. This is usually acomplished by a `docker pull`, followed by a `docker tag` and a `docker push`.

That approach has the downside of downloading the contents of every layer from Docker Hub, which has bandwidth and performance implications, especially in a CI environment.

This tool uses the [Docker Hub API](https://docs.docker.com/registry/spec/api/) to pull and push only a tiny [manifest](https://docs.docker.com/registry/spec/manifest-v2-2/) of the layers, bypassing the download overhead. Using this approach, an image of any size can be retagged in approximately 2 seconds.

## Installing

### From source

You can use `go get` to install this tool by running:

```bash
$ go get -u github.com/joshdk/docker-retag
```

### Precompiled binary

Alternatively, you can download a static Linux [release](https://github.com/joshdk/docker-retag/releases) binary by running:

```bash
$ wget -q https://github.com/joshdk/docker-retag/releases/download/0.0.1/docker-retag
$ sudo install docker-retag /usr/bin
```

## Usage

### Setup

Since `docker-retag` communicates with the [Docker Hub](https://hub.docker.com/) API, you must first export your account credentials into the working environment. These are the same credentials that you would use during `docker login`.

```bash
$ export DOCKER_USER='joshdk'
$ export DOCKER_PASS='hunter2'
```

### Examples

There are three argument forms for this tool. The first separately specifies the image name, current tag, and desired new tag.

```bash
$ docker-retag joshdk/hello-world 1.0.0 1.0.1
  Retagged joshdk/hello-world:1.0.0 as joshdk/hello-world:1.0.1
```

The second specifies the image name and current tag joined with a colon.

```bash
$ docker-retag joshdk/hello-world:1.0.0 1.0.1
  Retagged joshdk/hello-world:1.0.0 as joshdk/hello-world:1.0.1
```

The third defaults to `latest` if no tag is specified, similar to `docker pull`, etc.

```bash
$ docker-retag joshdk/hello-world 1.0.1
  Retagged joshdk/hello-world:latest as joshdk/hello-world:1.0.1
```

In all cases, the image and current tag **must** already exist in Docker Hub.

## License

This library is distributed under the [MIT License](https://opensource.org/licenses/MIT), see [LICENSE.txt](https://github.com/joshdk/docker-retag/blob/master/LICENSE.txt) for more information.