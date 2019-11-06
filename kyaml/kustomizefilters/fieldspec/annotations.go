// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package fieldspec

import (
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// KustomizeAnnotationsFilter sets Annotations on all Resources in a package, including sub-field
// annotations -- e.g. spec.template.metadata.annotations.
// Overrides existing annotations iff the keys match.
type KustomizeAnnotationsFilter struct {
	// commonAnnotations are the annotations to set
	Annotations map[string]*string `yaml:"commonAnnotations,omitempty"`
}

var _ kio.Filter = KustomizeAnnotationsFilter{}

func (af KustomizeAnnotationsFilter) Filter(input []*yaml.RNode) ([]*yaml.RNode, error) {
	for k, v := range af.Annotations {
		f := AnnotationFilter(k, v)
		_, err := kio.FilterAll(f).Filter(input)
		if err != nil {
			return nil, err
		}
	}
	return input, nil
}

func init() {
	err := yaml.Unmarshal([]byte(commonAnnotationFieldSpecs), &annotationReferenceFieldSpecs)
	if err != nil {
		panic(err)
	}
}

func AnnotationFilter(key string, value *string) *FieldSpecListFilter {
	if value != nil {
		return &FieldSpecListFilter{
			FieldSpecList: annotationReferenceFieldSpecs,
			SetValue: func(node *yaml.RNode) error {
				_, err := node.Pipe(yaml.FieldSetter{Name: key, StringValue: *value})
				return err
			},
			CreateKind: yaml.MappingNode,
		}
	}

	return &FieldSpecListFilter{
		FieldSpecList: annotationReferenceFieldSpecs,
		SetValue: func(node *yaml.RNode) error {
			_, err := node.Pipe(yaml.FieldClearer{Name: key})
			return err
		},
		CreateKind: yaml.MappingNode,
	}
}

var annotationReferenceFieldSpecs FieldSpecList

const commonAnnotationFieldSpecs = `
items:

# duck-type supported annotations
#
- path: metadata/annotations
  create: true
- path: spec/template/metadata/annotations
  create: false

# non-duck-type supported annotations
#
- path: spec/template/metadata/annotations
  create: true
  version: v1
  kind: ReplicationController

- path: spec/template/metadata/annotations
  create: true
  kind: Deployment

- path: spec/template/metadata/annotations
  create: true
  kind: ReplicaSet

- path: spec/template/metadata/annotations
  create: true
  kind: DaemonSet

- path: spec/template/metadata/annotations
  create: true
  kind: StatefulSet

- path: spec/template/metadata/annotations
  create: true
  group: batch
  kind: Job

- path: spec/jobTemplate/metadata/annotations
  create: true
  group: batch
  kind: CronJob

- path: spec/jobTemplate/spec/template/metadata/annotations
  create: true
  group: batch
  kind: CronJob

`
