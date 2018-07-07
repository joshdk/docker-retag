// Copyright 2018 Josh Komoroske. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file.

package arguments

import (
	"errors"
	"regexp"
	"strings"
)

var (
	errorInvalidArguments       = errors.New("invalid arguments")
	errorInvalidImageName       = errors.New("invalid image name")
	errorInvalidSourceReference = errors.New("invalid source reference")
	errorInvalidTargetReference = errors.New("invalid target reference")
)

func Parse(args []string) (string, string, string, error) {
	var (
		rawName                  string
		rawSourceReference       string
		rawTargetReference       string
		canonicalName            string
		canonicalSourceReference string
		canonicalTargetReference string
	)

	switch len(args) {
	case 3:
		rawName, rawSourceReference, rawTargetReference = args[0], args[1], args[2]
	case 2:
		rawName, rawSourceReference = splitImageRef(args[0])
		rawTargetReference = args[1]
	default:
		return "", "", "", errorInvalidArguments
	}

	if name, ok := maybeName(rawName); ok {
		canonicalName = name
	} else {
		return "", "", "", errorInvalidImageName
	}

	if digest, ok := maybeSHA256Digest(rawSourceReference); ok {
		canonicalSourceReference = digest
	} else if tag, ok := maybeTag(rawSourceReference); ok {
		canonicalSourceReference = tag
	} else {
		return "", "", "", errorInvalidSourceReference
	}

	if tag, ok := maybeTag(rawTargetReference); ok {
		canonicalTargetReference = tag
	} else {
		return "", "", "", errorInvalidTargetReference
	}

	return canonicalName, canonicalSourceReference, canonicalTargetReference, nil
}

// splitImageRef takes the given arg, and splits it into the name and
// reference components. If there is no reference, the tag :latest is used.
//
// The name and reference components are not suitable for direct consumption,
// and must be passed to maybeName or maybeTag/maybeSHA256Digest respectively.
func splitImageRef(arg string) (name string, reference string) {
	// Do we have a reference with a digest?
	// example/image@sha256:acd...def
	if chunks := strings.SplitN(arg, "@", 2); len(chunks) == 2 {
		return chunks[0], "@" + chunks[1]
	}

	// Do we have a reference with a tag?
	// example/image:1.2.3
	if chunks := strings.SplitN(arg, ":", 2); len(chunks) == 2 {
		return chunks[0], ":" + chunks[1]
	}

	return arg, ":latest"
}

// maybeName takes the given arg, and determines if it looks like an image name.
// If successful, a canonical name is returned, and a true value indicating
// success.
//
// The canonical name is acceptable to be used in Docker Hup API urls.
//
// This function attempts to preform some helpful input handling. Such as
// stripping any docker.io image name prefix, and handling library images that
// do not require an organization prefix.
func maybeName(arg string) (image string, isImage bool) {
	var (
		chunks = strings.Split(strings.TrimPrefix(arg, "docker.io/"), "/")
		org    string
		repo   string
	)

	switch len(chunks) {
	case 2:
		org, repo = chunks[0], chunks[1]
	case 1:
		org, repo = "library", chunks[0]
	}

	if len(org) == 0 || len(repo) == 0 {
		return "", false
	}

	return org + "/" + repo, true
}

// maybeTag takes the given arg, and determines if it looks like a tag based
// image reference. If successful, a canonical tag is returned, and a true
// value indicating success.
//
// The canonical tag is acceptable to be used in Docker Hup API urls.
func maybeTag(arg string) (tag string, isTag bool) {
	//https://github.com/docker/distribution/blob/master/reference/regexp.go#L36
	tagRegex := regexp.MustCompile(`^:?([\w][\w.-]{0,127})$`)
	chunks := tagRegex.FindStringSubmatch(arg)

	if len(chunks) == 2 {
		return chunks[1], true
	}
	return "", false
}

// maybeSHA256Digest takes the given arg, and determines if it looks like a
// digest based image reference. If successful, a canonical digest is returned,
// and a true value indicating success.
//
// The canonical digest is acceptable to be used in Docker Hup API urls.
//
// This function is specifically more restrictive with allowed digests compared
// to what is allowed by the specification. The sha256 digest type is the only
// implementation at the moment, so it is safer to limit against that.
func maybeSHA256Digest(arg string) (digest string, isDigest bool) {
	//https://github.com/docker/distribution/blob/master/reference/regexp.go#L43
	digestRegex := regexp.MustCompile(`^(@|sha256:|@sha256:)([0-9a-f]{64})$`)
	chunks := digestRegex.FindStringSubmatch(arg)

	if len(chunks) == 3 {
		return "sha256:" + chunks[2], true
	}
	return "", false
}
