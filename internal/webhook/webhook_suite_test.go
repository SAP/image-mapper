/*
SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company and image-mapper contributors
SPDX-License-Identifier: Apache-2.0
*/

package webhook_test

import (
	"context"
	"testing"

	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"

	"github.com/sap/image-mapper/internal/webhook"
)

func TestWebhook(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Webhook Suite")
}

var _ = Describe("Webhook", func() {
	Context("Generic Webhook", func() {
		var mutate func(pod *corev1.Pod) error
		var pod *corev1.Pod
		var container1 *corev1.Container
		var container2 *corev1.Container
		var initContainer1 *corev1.Container
		var initContainer2 *corev1.Container

		BeforeEach(func() {
			podWebhook := &webhook.PodWebhook{
				Config: &webhook.PodWebhookConfig{
					Mappings: []webhook.Mapping{
						parseMapping(`docker\.io/.*`, `my-docker-cache.io/${repository}:${tag}`),
						parseMapping(`registry\.io/.*`, `my-registry-cache.io/${repository}:${tag}`),
						parseMapping(`^other-registry\.io/([^/]+)/(.+)$`, `my-registry-cache.io/$1/foo/$2`),
					},
					LabelsIfModified: []webhook.Label{
						parseLabel("mutated-by-test-label=true"),
					},
					AnnotationsIfModified: []webhook.Annotation{
						parseAnnotation("mutated-by-test-annotation=yes"),
					},
				},
			}

			mutate = func(pod *corev1.Pod) error {
				ctx := logr.NewContext(context.TODO(), logr.Discard())
				return podWebhook.MutateCreate(ctx, pod)
			}

			pod = &corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "container1",
							Image: "not.to.be.replaced/repo:tag",
						},
						{
							Name:  "container2",
							Image: "not.to.be.replaced/repo:tag",
						},
					},
					InitContainers: []corev1.Container{
						{
							Name:  "initContainer1",
							Image: "not.to.be.replaced/repo:tag",
						},
						{
							Name:  "initContainer2",
							Image: "not.to.be.replaced/repo:tag",
						},
					},
				},
			}
			container1 = &pod.Spec.Containers[0]
			container2 = &pod.Spec.Containers[1]
			initContainer1 = &pod.Spec.InitContainers[0]
			initContainer2 = &pod.Spec.InitContainers[1]
		})

		It("should replace images correctly (1)", func() {
			container1.Image = "nginx"
			initContainer2.Image = "registry.io/repo/nginx:v123"

			err := mutate(pod)
			Expect(err).NotTo(HaveOccurred())

			Expect(container1.Image).To(Equal("my-docker-cache.io/nginx:latest"))
			Expect(container2.Image).To(Equal("not.to.be.replaced/repo:tag"))
			Expect(initContainer1.Image).To(Equal("not.to.be.replaced/repo:tag"))
			Expect(initContainer2.Image).To(Equal("my-registry-cache.io/repo/nginx:v123"))
		})

		It("should replace images correctly (2)", func() {
			container1.Image = "registry.io/nginx"
			container2.Image = "other-registry.io/repo/alpine:v1"
			initContainer1.Image = "registry.io/repo/ubuntu:v2"

			err := mutate(pod)
			Expect(err).NotTo(HaveOccurred())

			Expect(container1.Image).To(Equal("my-registry-cache.io/nginx:latest"))
			Expect(container2.Image).To(Equal("my-registry-cache.io/repo/foo/alpine:v1"))
			Expect(initContainer1.Image).To(Equal("my-registry-cache.io/repo/ubuntu:v2"))
			Expect(initContainer2.Image).To(Equal("not.to.be.replaced/repo:tag"))
		})
	})
})

func parseMapping(pattern string, replacement string) webhook.Mapping {
	mapping, err := webhook.ParseMapping(pattern, replacement)
	if err != nil {
		panic(err)
	}
	return mapping
}

func parseLabel(s string) webhook.Label {
	label, err := webhook.ParseLabel(s)
	if err != nil {
		panic(err)
	}
	return label
}

func parseAnnotation(s string) webhook.Annotation {
	annotation, err := webhook.ParseAnnotation(s)
	if err != nil {
		panic(err)
	}
	return annotation
}
