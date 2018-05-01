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
		repository, oldTag, newTag, err = parseArgs(args)
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

	fmt.Printf("Retagged %s:%s as %s:%s\n", repository, oldTag, repository, newTag)

	return nil
}

func parseArgs(args []string) (string, string, string, error) {
	switch len(args) {
	case 4:
		// given:  "docker-retag", "repo/product", "1.2.3", "4.5.6"
		// return: "repo/product", "1.2.3", "4.5.6", nil
		return args[1], args[2], args[3], nil

	case 3:
		chunks := strings.SplitN(args[1], ":", 2)
		if len(chunks) == 2 {

			// given:  "docker-retag", "repo/product:1.2.3", "4.5.6"
			// return: "repo/product", "1.2.3", "4.5.6", nil
			return chunks[0], chunks[1], args[2], nil
		}

		// given:  "docker-retag", "repo/product", "4.5.6"
		// return: "repo/product", "latest", "4.5.6", nil
		return chunks[0], "latest", args[2], nil

	default:
		return "", "", "", errors.New("invalid arguments")
	}
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
