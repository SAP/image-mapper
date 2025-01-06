/*
SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company and image-mapper contributors
SPDX-License-Identifier: Apache-2.0
*/

package webhook

import (
	"fmt"
	"regexp"
)

type Label struct {
	Key   string
	Value string
}

func ParseLabel(label string) (Label, error) {
	// todo: make the regular expression more accurate, and check for the exact syntax: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels
	if match := regexp.MustCompile(`^([a-zA-Z0-9.\-_/]+)=(.*)$`).FindAllStringSubmatch(label, -1); match != nil {
		return Label{Key: match[0][1], Value: match[0][2]}, nil
	} else {
		return Label{}, fmt.Errorf("invalid label expression: %s", label)
	}
}
