// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package patches

import (
	"strings"

	"sigs.k8s.io/kustomize/kyaml/yaml"
	"sigs.k8s.io/kustomize/kyaml/yaml/merge2"
)

type PatchFilter struct {
	Patches []Patch `yaml:"patches"`
}

type Patch struct {
	Patch  yaml.Node `yaml:"patch"`
	Target Target    `yaml:"target,omitempty"`
}

type Target struct {
	Kind               string            `yaml:"kind,omitempty"`
	ApiVersion         string            `yaml:"apiVersion,omitempty"`
	Name               string            `yaml:"name,omitempty"`
	Namespace          string            `yaml:"namespace,omitempty"`
	LabelSelector      map[string]string `yaml:"labelSelector,omitempty"`
	LabelPrefix        map[string]string `yaml:"labelPrefix,omitempty"`
	AnnotationSelector map[string]string `yaml:"annotationSelector,omitempty"`
	AnnotationPrefix   map[string]string `yaml:"annotationPrefix,omitempty"`
}

func (pf *PatchFilter) Filter(input []*yaml.RNode) ([]*yaml.RNode, error) {
	for i := range input {
		for j := range pf.Patches {
			if match, err := pf.match(input[i], pf.Patches[j].Target); err != nil {
				return nil, err
			} else if !match {
				continue
			}

			// make a copy of the patch so when fields are copied to the destination,
			// each has a unique copy
			b, err := yaml.Marshal(&pf.Patches[j].Patch)
			if err != nil {
				return nil, err
			}
			var patch yaml.Node
			err = yaml.Unmarshal(b, &patch)
			if err != nil {
				return nil, err
			}

			// merge the patch into the node
			input[i], err = merge2.Merge(yaml.NewRNode(patch.Content[0]), input[i])
			if err != nil {
				return nil, err
			}
		}
	}
	return input, nil
}

func (pf *PatchFilter) match(node *yaml.RNode, target Target) (bool, error) {
	meta, err := node.GetMeta()
	if err != nil {
		return false, err
	}
	if meta.Annotations == nil {
		meta.Annotations = map[string]string{}
	}
	if meta.Labels == nil {
		meta.Labels = map[string]string{}
	}

	if target.Name != "" && target.Name != meta.Name {
		return false, nil
	}
	if target.Kind != "" && target.Kind != meta.Kind {
		return false, nil
	}
	if target.ApiVersion != "" && target.ApiVersion != meta.ApiVersion {
		return false, nil
	}
	if target.Namespace != "" && target.Namespace != meta.Namespace {
		return false, nil
	}

	for k, v := range target.AnnotationSelector {
		if meta.Annotations[k] != v {
			return false, nil
		}
	}
	for k, v := range target.AnnotationPrefix {
		if !strings.HasPrefix(meta.Annotations[k], v) {
			return false, nil
		}
	}
	for k, v := range target.LabelSelector {
		if meta.Labels[k] != v {
			return false, nil
		}
	}
	for k, v := range target.LabelPrefix {
		if !strings.HasPrefix(meta.Labels[k], v) {
			return false, nil
		}
	}
	return true, nil
}
