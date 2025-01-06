/*
SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company and image-mapper contributors
SPDX-License-Identifier: Apache-2.0
*/

package webhook

import (
	"fmt"
	"regexp"
)

type Annotation struct {
	Key   string
	Value string
}

func ParseAnnotation(annotation string) (Annotation, error) {
	// todo: make the regular expression more accurate, and check for the exact syntax: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations
	if match := regexp.MustCompile(`^([a-zA-Z0-9.\-_/]+)=(.*)$`).FindAllStringSubmatch(annotation, -1); match != nil {
		return Annotation{Key: match[0][1], Value: match[0][2]}, nil
	} else {
		return Annotation{}, fmt.Errorf("invalid annotation expression: %s", annotation)
	}
}
