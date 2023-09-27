/*
SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and image-mapper contributors
SPDX-License-Identifier: Apache-2.0
*/

package webhook

import (
	"fmt"
	"regexp"
	"strings"
)

type Mapping struct {
	Pattern           string `json:"pattern"`
	Replacement       string `json:"replacement"`
	PatternRegexp     *regexp.Regexp
	ReplacementRegexp *regexp.Regexp
}

func ParseMapping(pattern string, replacement string) (Mapping, error) {
	mapping := Mapping{
		Pattern:     pattern,
		Replacement: replacement,
	}
	var err error
	if mapping.Pattern == "" {
		return Mapping{}, fmt.Errorf("empty pattern")
	}
	if !strings.HasPrefix(mapping.Pattern, "^") {
		mapping.Pattern = "^" + mapping.Pattern
	}
	if !strings.HasSuffix(mapping.Pattern, "$") {
		mapping.Pattern = mapping.Pattern + "$"
	}
	if strings.Contains(mapping.Pattern, "(") {
		mapping.PatternRegexp, err = regexp.Compile(mapping.Pattern)
		if err != nil {
			return Mapping{}, err
		}
		mapping.ReplacementRegexp = mapping.PatternRegexp
	} else {
		mapping.PatternRegexp, err = regexp.Compile(mapping.Pattern)
		if err != nil {
			return Mapping{}, err
		}
		mapping.ReplacementRegexp, err = regexp.Compile(`^(?:(?P<registry>(?:[a-z0-9-]+\.)+[a-z0-9-]+(?::(?:[0-9]+))?)/)?(?P<repository>[a-z0-9_\-./]+)(?::(?P<tag>[a-zA-Z0-9_\-.]+))?(?:@(?P<digest>sha256:[a-f0-9]+))?$`)
		if err != nil {
			return Mapping{}, err
		}
	}
	if mapping.Replacement == "" {
		return Mapping{}, fmt.Errorf("empty replacement")
	}
	return mapping, nil
}
