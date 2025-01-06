/*
SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company and image-mapper contributors
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"context"
	"flag"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"

	"github.com/sap/admission-webhook-runtime/pkg/admission"

	"github.com/sap/image-mapper/internal/webhook"
)

func main() {
	admission.InitFlags(nil)
	defineFlags()
	klog.InitFlags(nil)
	flag.Parse()

	config, err := buildConfigFromFlags()
	if err != nil {
		klog.Fatalf("error loading config: %s", err)
	}

	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		klog.Fatalf("error popuplating scheme: %s", err)
	}

	webhook := &webhook.PodWebhook{Config: config}
	if err := admission.RegisterMutatingWebhook[*corev1.Pod](webhook, scheme, klogr.New()); err != nil {
		klog.Fatalf("error registering webhook: %s", err)
	}

	if err := admission.Serve(context.Background(), nil); err != nil {
		klog.Fatalf("error serving webhook: %s", err)
	}
}
