// Copyright 2018 Josh Komoroske. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file.

package main

import (
	"fmt"
	"os"
)

func main() {
	if err := mainCmd(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "docker-retag: %s\n", err.Error())
		os.Exit(1)
	}
}

func mainCmd(args []string) error {
	return nil
}
