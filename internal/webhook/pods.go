/*
SPDX-FileCopyrightText: 2026 SAP SE or an SAP affiliate company and image-mapper contributors
SPDX-License-Identifier: Apache-2.0
*/

package webhook

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
)

type PodWebhookConfig struct {
	Mappings              []Mapping
	LabelsIfModified      []Label
	AnnotationsIfModified []Annotation
}

type PodWebhook struct {
	Config *PodWebhookConfig
}

func (w *PodWebhook) MutateCreate(ctx context.Context, pod *corev1.Pod) error {
	return w.mutate(ctx, pod)
}

func (w *PodWebhook) MutateUpdate(ctx context.Context, oldPod *corev1.Pod, newPod *corev1.Pod) error {
	return w.mutate(ctx, newPod)
}

func (w *PodWebhook) mutate(ctx context.Context, pod *corev1.Pod) error {
	log, err := logr.FromContext(ctx)
	if err != nil {
		panic(err)
	}

	namespace := pod.Namespace
	name := pod.Name
	if name == "" {
		name = pod.GenerateName + "<new>"
	}
	numChanges := 0

	for i := 0; i < len(pod.Spec.Containers); i++ {
		container := &pod.Spec.Containers[i]
		currentImage, err := normalizeImage(container.Image)
		if err != nil {
			log.Error(err, "error norrmalizing image for container", "namespace", namespace, "pod", name, "container", container.Name)
			return errors.Wrapf(err, "error norrmalizing image for container %s/%s/%s", namespace, name, container.Name)
		}
		newImage := w.getImage(currentImage)
		if newImage == currentImage {
			log.V(2).Info("not changing image for container", "namespace", namespace, "pod", name, "container", container.Name)
		} else {
			log.Info("changing image for container", namespace, "pod", name, "container", container.Name, "oldImage", currentImage, "newImage", newImage)
			container.Image = newImage
			numChanges++
		}
	}

	for i := 0; i < len(pod.Spec.InitContainers); i++ {
		container := &pod.Spec.InitContainers[i]
		currentImage, err := normalizeImage(container.Image)
		if err != nil {
			log.Error(err, "error norrmalizing image for init container", "namespace", namespace, "pod", name, "container", container.Name)
			return errors.Wrapf(err, "error norrmalizing image for init container %s/%s/%s", namespace, name, container.Name)
		}
		newImage := w.getImage(currentImage)
		if newImage == currentImage {
			log.V(2).Info("not changing image for init container", "namespace", namespace, "pod", name, "container", container.Name)
		} else {
			log.Info("changing image for init container", namespace, "pod", name, "container", container.Name, "oldImage", currentImage, "newImage", newImage)
			container.Image = newImage
			numChanges++
		}
	}

	if numChanges > 0 {
		for _, label := range w.Config.LabelsIfModified {
			if pod.Labels == nil {
				pod.Labels = make(map[string]string)
			}
			pod.Labels[label.Key] = label.Value
		}
		for _, annotation := range w.Config.AnnotationsIfModified {
			if pod.Annotations == nil {
				pod.Annotations = make(map[string]string)
			}
			pod.Annotations[annotation.Key] = annotation.Value
		}
	}

	return nil
}

func (w *PodWebhook) getImage(image string) string {
	for _, mapping := range w.Config.Mappings {
		if mapping.PatternRegexp.MatchString(image) {
			return mapping.ReplacementRegexp.ReplaceAllString(image, mapping.Replacement)
		}
	}
	return image
}
