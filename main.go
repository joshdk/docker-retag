// Copyright 2018 Josh Komoroske. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file.

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joshdk/docker-retag/arguments"
	"strings"
)

const (
	dockerUsernameEnv = "DOCKER_USER"
	dockerPasswordEnv = "DOCKER_PASS"
)

func main() {
	if err := mainCmd(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "docker-retag: %s\n", err.Error())
		os.Exit(1)
	}
}

func mainCmd(args []string) error {
	var (
		repository, oldTag, newTag, err = arguments.Parse(args[1:])
	)

	if err != nil {
		return err
	}

	username, found := os.LookupEnv(dockerUsernameEnv)
	if !found {
		return errors.New(dockerUsernameEnv + " not found in environment")
	}

	password, found := os.LookupEnv(dockerPasswordEnv)
	if !found {
		return errors.New(dockerPasswordEnv + " not found in environment")
	}

	token, err := login(repository, username, password)
	if err != nil {
		return errors.New("failed to authenticate: " + err.Error())
	}

	manifest, err := pullManifest(token, repository, oldTag)
	if err != nil {
		return errors.New("failed to pull manifest: " + err.Error())
	}

	if err := pushManifest(token, repository, newTag, manifest); err != nil {
		return errors.New("failed to push manifest: " + err.Error())
	}

	separator := ":"
	if strings.HasPrefix(oldTag, "sha256:") {
		separator = "@"
	}

	fmt.Printf("Retagged %s%s%s as %s:%s\n", repository, separator, oldTag, repository, newTag)

	return nil
}

func login(repo string, username string, password string) (string, error) {
	var (
		client = http.DefaultClient
		url    = "https://auth.docker.io/token?service=registry.docker.io&scope=repository:" + repo + ":pull,push"
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data struct {
		Details string `json:"details"`
		Token   string `json:"token"`
	}

	if err := json.Unmarshal(bodyText, &data); err != nil {
		return "", err
	}

	if data.Token == "" {
		return "", errors.New("empty token")
	}

	return data.Token, nil
}

func pullManifest(token string, repository string, tag string) ([]byte, error) {
	var (
		client = http.DefaultClient
		url    = "https://index.docker.io/v2/" + repository + "/manifests/" + tag
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bodyText, nil
}

func pushManifest(token string, repository string, tag string, manifest []byte) error {
	var (
		client = http.DefaultClient
		url    = "https://index.docker.io/v2/" + repository + "/manifests/" + tag
	)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(manifest))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-type", "application/vnd.docker.distribution.manifest.v2+json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return errors.New(resp.Status)
	}

	return nil
}
