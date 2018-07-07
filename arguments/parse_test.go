// Copyright 2018 Josh Komoroske. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file.

package arguments

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {

	tests := []struct {
		title                   string
		args                    []string
		expectedName            string
		expectedSourceReference string
		expectedTargetReference string
		expectedError           error
	}{
		{
			title:         "Nil args",
			expectedError: errorInvalidArguments,
		},
		{
			title:         "Empty args",
			args:          []string{},
			expectedError: errorInvalidArguments,
		},
		{
			title:         "Too many components",
			args:          []string{"docker.io/library/org/example:1.2.3", ":4.5.6"},
			expectedError: errorInvalidImageName,
		},
		{
			title:         "Target digest",
			args:          []string{"org/example", ":1.2.3", "@00000000aaaaaaaa11111111bbbbbbbb22222222cccccccc33333333dddddddd"},
			expectedError: errorInvalidTargetReference,
		},
		{
			title:         "Source digest with colon",
			args:          []string{"org/example", ":sha256:00000000aaaaaaaa11111111bbbbbbbb22222222cccccccc33333333dddddddd", "4.5.6"},
			expectedError: errorInvalidSourceReference,
		},
		{
			"Tags",
			[]string{"org/example", "1.2.3", "4.5.6"},
			"org/example", "1.2.3", "4.5.6",
			nil,
		},
		{
			"Tags with colons",
			[]string{"org/example", ":1.2.3", ":4.5.6"},
			"org/example", "1.2.3", "4.5.6",
			nil,
		},
		{
			"Digest with @",
			[]string{"org/example", "@00000000aaaaaaaa11111111bbbbbbbb22222222cccccccc33333333dddddddd", ":4.5.6"},
			"org/example", "sha256:00000000aaaaaaaa11111111bbbbbbbb22222222cccccccc33333333dddddddd", "4.5.6",
			nil,
		},
		{
			"Digest with sha256",
			[]string{"org/example", "sha256:00000000aaaaaaaa11111111bbbbbbbb22222222cccccccc33333333dddddddd", ":4.5.6"},
			"org/example", "sha256:00000000aaaaaaaa11111111bbbbbbbb22222222cccccccc33333333dddddddd", "4.5.6",
			nil,
		},
		{
			"Digest with @sha256",
			[]string{"org/example", "@sha256:00000000aaaaaaaa11111111bbbbbbbb22222222cccccccc33333333dddddddd", ":4.5.6"},
			"org/example", "sha256:00000000aaaaaaaa11111111bbbbbbbb22222222cccccccc33333333dddddddd", "4.5.6",
			nil,
		},
		{
			"No tag",
			[]string{"org/example", ":4.5.6"},
			"org/example", "latest", "4.5.6",
			nil,
		},
		{
			"Joined tag",
			[]string{"org/example:1.2.3", ":4.5.6"},
			"org/example", "1.2.3", "4.5.6",
			nil,
		},
		{
			"Joined digest with @",
			[]string{"org/example@00000000aaaaaaaa11111111bbbbbbbb22222222cccccccc33333333dddddddd", ":4.5.6"},
			"org/example", "sha256:00000000aaaaaaaa11111111bbbbbbbb22222222cccccccc33333333dddddddd", "4.5.6",
			nil,
		},
		{
			"Joined digest with @sha256",
			[]string{"org/example@sha256:00000000aaaaaaaa11111111bbbbbbbb22222222cccccccc33333333dddddddd", ":4.5.6"},
			"org/example", "sha256:00000000aaaaaaaa11111111bbbbbbbb22222222cccccccc33333333dddddddd", "4.5.6",
			nil,
		},
	}

	for index, test := range tests {
		name := fmt.Sprintf("Case #%d %s", index+1, test.title)
		t.Run(name, func(t *testing.T) {

			name, sourceReference, targetReference, err := Parse(test.args)
			if err != test.expectedError {
				panic(fmt.Sprintf("Expected %v, actual %v", test.expectedError, err))
			}
			if name != test.expectedName {
				panic(fmt.Sprintf("Expected %s, actual %s", test.expectedName, name))
			}
			if sourceReference != test.expectedSourceReference {
				panic(fmt.Sprintf("Expected %s, actual %s", test.expectedSourceReference, sourceReference))
			}
			if targetReference != test.expectedTargetReference {
				panic(fmt.Sprintf("Expected %s, actual %s", test.expectedTargetReference, targetReference))
			}
		})
	}
}

func TestSplitImageRef(t *testing.T) {

	tests := []struct {
		title             string
		arg               string
		expectedName      string
		expectedReference string
	}{
		{
			"Empty string",
			"",
			"", ":latest",
		},
		{
			"Only colon",
			":",
			"", ":",
		},
		{
			"Only @",
			"@",
			"", "@",
		},
		{
			"Bare image",
			"org/example",
			"org/example", ":latest",
		},
		{
			"Image tagged latest",
			"org/example:latest",
			"org/example", ":latest",
		},
		{
			"Image tagged with semver",
			"org/example:1.2.3",
			"org/example", ":1.2.3",
		},
		{
			"Image tagged with double colon",
			"org/example::1.2.3",
			"org/example", "::1.2.3",
		},
		{
			"Image with digest",
			"org/example@00000000aaaaaaaa11111111bbbbbbbb22222222cccccccc33333333dddddddd",
			"org/example", "@00000000aaaaaaaa11111111bbbbbbbb22222222cccccccc33333333dddddddd",
		},
		{
			"Image with sha256 digest",
			"org/example@sha256:00000000aaaaaaaa11111111bbbbbbbb22222222cccccccc33333333dddddddd",
			"org/example", "@sha256:00000000aaaaaaaa11111111bbbbbbbb22222222cccccccc33333333dddddddd",
		},
		{
			"Image tagged with double @",
			"org/example@@00000000aaaaaaaa11111111bbbbbbbb22222222cccccccc33333333dddddddd",
			"org/example", "@@00000000aaaaaaaa11111111bbbbbbbb22222222cccccccc33333333dddddddd",
		},
	}

	for index, test := range tests {
		name := fmt.Sprintf("Case #%d %s", index+1, test.title)
		t.Run(name, func(t *testing.T) {

			name, reference := splitImageRef(test.arg)
			if name != test.expectedName {
				panic(fmt.Sprintf("Expected %s, actual %s", test.expectedName, name))
			}
			if reference != test.expectedReference {
				panic(fmt.Sprintf("Expected %s, actual %s", test.expectedReference, reference))
			}
		})
	}
}

func TestMaybeName(t *testing.T) {

	tests := []struct {
		title          string
		arg            string
		expectedName   string
		expectedIsName bool
	}{
		{
			title: "Empty string",
		},
		{
			title: "Too many components",
			arg:   "docker.io/library/org/example",
		},
		{
			title: "Empty org",
			arg:   "/example",
		},
		{
			title: "Empty repo",
			arg:   "org/",
		},
		{
			title: "Empty org and repo",
			arg:   "/",
		},
		{
			title: "Too many slashes",
			arg:   "org//example",
		},
		{
			title: "Leading slash",
			arg:   "/org/example",
		},
		{
			title: "Trailing slash",
			arg:   "org/example/",
		},
		{
			"Bare library image",
			"example", "library/example",
			true,
		},
		{
			"Library image with library prefix",
			"library/example", "library/example",
			true,
		},
		{
			"Library image with docker.io prefix",
			"docker.io/example", "library/example",
			true,
		},
		{
			"Library image with both docker.io and library prefix",
			"docker.io/library/example", "library/example",
			true,
		},
		{
			"Org image",
			"org/example", "org/example",
			true,
		},
		{
			"Org image with docker.io prefix",
			"docker.io/org/example", "org/example",
			true,
		},
	}

	for index, test := range tests {
		name := fmt.Sprintf("Case #%d %s", index+1, test.title)
		t.Run(name, func(t *testing.T) {

			name, isName := maybeName(test.arg)
			if isName != test.expectedIsName {
				panic(fmt.Sprintf("Expected %v, actual %v", test.expectedIsName, isName))
			}
			if name != test.expectedName {
				panic(fmt.Sprintf("Expected %s, actual %s", test.expectedName, name))
			}
		})
	}
}

func TestMaybeTag(t *testing.T) {

	tests := []struct {
		title         string
		arg           string
		expectedTag   string
		expectedIsTag bool
	}{
		{
			title: "Empty string",
		},
		{
			title: "Only colon",
			arg:   ":",
		},
		{
			title: "Double colon",
			arg:   "::1.2.3",
		},
		{
			"Word tag",
			"latest", "latest",
			true,
		},
		{
			"Word tag with colon",
			":latest", "latest",
			true,
		},
		{
			"Semver tag",
			"1.2.3", "1.2.3",
			true,
		},
		{
			"Semver tag with colon",
			":1.2.3", "1.2.3",
			true,
		},
		{
			"Complex example tag 1",
			":7u181-jre-alpine3.7", "7u181-jre-alpine3.7",
			true,
		},
		{
			"Complex example tag 2",
			":1.10rc2-stretch", "1.10rc2-stretch",
			true,
		},
		{
			"Complex example tag 3",
			":4.0.0-rc6-xenial", "4.0.0-rc6-xenial",
			true,
		},
	}

	for index, test := range tests {
		name := fmt.Sprintf("Case #%d %s", index+1, test.title)
		t.Run(name, func(t *testing.T) {

			tag, isTag := maybeTag(test.arg)
			if isTag != test.expectedIsTag {
				panic(fmt.Sprintf("Expected %v, actual %v", test.expectedIsTag, isTag))
			}
			if tag != test.expectedTag {
				panic(fmt.Sprintf("Expected %s, actual %s", test.expectedTag, tag))
			}
		})
	}
}

func TestMaybeDigest(t *testing.T) {

	tests := []struct {
		title            string
		arg              string
		expectedDigest   string
		expectedIsDigest bool
	}{
		{
			title: "Empty string",
		},
		{
			title: "Only @",
			arg:   "@",
		},
		{
			title: "Double @",
			arg:   "@@00000000aaaaaaaa11111111bbbbbbbb22222222cccccccc33333333dddddddd",
		},
		{
			title: "Digest too short",
			arg:   "@00000000aaaaaaaa11111111bbbbbbbb22222222cccccccc33333333ddddddd",
		},
		{
			title: "Digest too long",
			arg:   "@00000000aaaaaaaa11111111bbbbbbbb22222222cccccccc33333333dddddddd4",
		},
		{
			title: "Digest uppercase hex",
			arg:   "@00000000aaaaaaaa11111111bbbbbbbb22222222cccccccc33333333dddddddD",
		},
		{
			title: "Digest illegal letter",
			arg:   "@O0000000aaaaaaaal1111111bbbbbbbb22222222cccccccc33333333dddddddd",
		},
		{
			title: "Wrong digest type",
			arg:   "@sha512:00000000aaaaaaaa11111111bbbbbbbb22222222cccccccc33333333dddddddd",
		},
		{
			"Digest with @",
			"@00000000aaaaaaaa11111111bbbbbbbb22222222cccccccc33333333dddddddd",
			"sha256:00000000aaaaaaaa11111111bbbbbbbb22222222cccccccc33333333dddddddd",
			true,
		},
		{
			"Digest with sha256",
			"sha256:00000000aaaaaaaa11111111bbbbbbbb22222222cccccccc33333333dddddddd",
			"sha256:00000000aaaaaaaa11111111bbbbbbbb22222222cccccccc33333333dddddddd",
			true,
		},
		{
			"Digest with @sha256",
			"@sha256:00000000aaaaaaaa11111111bbbbbbbb22222222cccccccc33333333dddddddd",
			"sha256:00000000aaaaaaaa11111111bbbbbbbb22222222cccccccc33333333dddddddd",
			true,
		},
	}

	for index, test := range tests {
		name := fmt.Sprintf("Case #%d %s", index+1, test.title)
		t.Run(name, func(t *testing.T) {

			digest, isDigest := maybeSHA256Digest(test.arg)
			if isDigest != test.expectedIsDigest {
				panic(fmt.Sprintf("Expected %v, actual %v", test.expectedIsDigest, isDigest))
			}
			if digest != test.expectedDigest {
				panic(fmt.Sprintf("Expected %s, actual %s", test.expectedDigest, digest))
			}
		})
	}
}
