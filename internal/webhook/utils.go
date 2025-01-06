/*
SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company and image-mapper contributors
SPDX-License-Identifier: Apache-2.0
*/

package webhook

import (
	"fmt"
	"regexp"
)

func normalizeImage(image string) (string, error) {
	match := regexp.MustCompile(`^(?:(?P<registry>(?:[a-z0-9-]+\.)+[a-z0-9-]+(?::(?:[0-9]+))?)/)?(?P<repository>[a-z0-9_\-./]+)(?::(?P<tag>[a-zA-Z0-9_\-.]+))?(?:@(?P<digest>sha256:[a-f0-9]+))?$`).FindAllStringSubmatch(image, -1)
	if match != nil {
		registry := match[0][1]
		if registry == "" {
			registry = "docker.io"
		}
		repository := match[0][2]
		tag := match[0][3]
		if tag == "" {
			tag = "latest"
		}
		return registry + "/" + repository + ":" + tag, nil
	} else {
		return "", fmt.Errorf("image %s has invalid format", image)
	}
}
