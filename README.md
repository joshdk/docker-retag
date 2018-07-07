[![CircleCI][circleci-badge]][circleci-link]
[![Go Report Card][go-report-card-badge]][go-report-card-link]
[![License][license-badge]][license-link]
[![Github downloads][github-downloads-badge]][github-release-link]
[![GitHub release][github-release-badge]][github-release-link]

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

Alternatively, you can download a static Linux [release][github-release-link] binary by running:

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

This library is distributed under the [MIT License][license-link], see [LICENSE.txt][license-file] for more information.

[circleci-badge]:         https://circleci.com/gh/joshdk/docker-retag.svg?&style=shield
[circleci-link]:          https://circleci.com/gh/joshdk/docker-retag/tree/master
[github-downloads-badge]: https://img.shields.io/github/downloads/joshdk/docker-retag/total.svg
[github-release-badge]:   https://img.shields.io/github/release/joshdk/docker-retag.svg
[github-release-link]:    https://github.com/joshdk/docker-retag/releases/latest
[go-report-card-badge]:   https://goreportcard.com/badge/github.com/joshdk/docker-retag
[go-report-card-link]:    https://goreportcard.com/report/github.com/joshdk/docker-retag
[license-badge]:          https://img.shields.io/github/license/joshdk/docker-retag.svg
[license-file]:           https://github.com/joshdk/docker-retag/blob/master/LICENSE.txt
[license-link]:           https://opensource.org/licenses/MIT
