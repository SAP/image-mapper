/*
SPDX-FileCopyrightText: 2026 SAP SE or an SAP affiliate company and image-mapper contributors
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/sap/image-mapper/internal/webhook"
)

type StringArray []string

func (v *StringArray) String() string {
	if v == nil {
		return ""
	} else {
		return strings.Join(*v, ",")
	}
}

func (v *StringArray) Set(value string) error {
	*v = append(*v, value)
	return nil
}

type StringPair struct {
	A string
	B string
}

type MappingRule struct {
	Pattern     string
	Replacement string
}

var mappingFile string
var labelsIfModified StringArray
var annotationsIfModified StringArray

func defineFlags() {
	flag.StringVar(&mappingFile, "mapping-file", mappingFile, "File containing the image mappings")
	flag.Var(&labelsIfModified, "add-label-if-modified", "Add specified label (given as key=value) to pod if some image was changed (can be repeated)")
	flag.Var(&annotationsIfModified, "add-annotation-if-modified", "Add specified annotation (given as key=value) to pod if some image was changed (can be repeated)")
}

func buildConfigFromFlags() (*webhook.PodWebhookConfig, error) {
	config := &webhook.PodWebhookConfig{}

	if mappingFile == "" {
		return nil, fmt.Errorf("missing flag: -mapping-file")
	}
	mappingData, err := os.ReadFile(mappingFile)
	if err != nil {
		return nil, errors.Wrap(err, "error loading mapping file")
	}
	var mappingRules []MappingRule
	if err := json.Unmarshal(mappingData, &mappingRules); err != nil {
		return nil, errors.Wrap(err, "parsing mapping file")
	}
	for i, mappingRule := range mappingRules {
		mapping, err := webhook.ParseMapping(mappingRule.Pattern, mappingRule.Replacement)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid mapping [%d]", i)
		}
		config.Mappings = append(config.Mappings, mapping)
	}

	for _, s := range labelsIfModified {
		label, err := webhook.ParseLabel(s)
		if err != nil {
			return nil, err
		}
		config.LabelsIfModified = append(config.LabelsIfModified, label)
	}

	for _, s := range annotationsIfModified {
		annotation, err := webhook.ParseAnnotation(s)
		if err != nil {
			return nil, err
		}
		config.AnnotationsIfModified = append(config.AnnotationsIfModified, annotation)
	}

	return config, nil
}
