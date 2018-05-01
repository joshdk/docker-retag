// Copyright 2018 Josh Komoroske. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file.

package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
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

	fmt.Println(repository, oldTag, newTag)

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
