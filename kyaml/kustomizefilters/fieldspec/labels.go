// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package fieldspec

import (
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// KustomizeLabelsFilter sets labels on all Resources in a package, including sub-field
// labels -- e.g. spec.template.metadata.labels.
// Overrides existing labels iff the keys match.
type KustomizeLabelsFilter struct {
	// commonLabels are the labels to set
	Labels map[string]*string `yaml:"commonLabels,omitempty"`
}

var _ kio.Filter = KustomizeLabelsFilter{}

func (lf KustomizeLabelsFilter) Filter(input []*yaml.RNode) ([]*yaml.RNode, error) {
	for k, v := range lf.Labels {
		f := labelsFilter(k, v)
		_, err := kio.FilterAll(f).Filter(input)
		if err != nil {
			return nil, err
		}
	}
	return input, nil
}

func init() {
	err := yaml.Unmarshal([]byte(commonLabelFieldSpecs), &labelReferenceFieldSpecs)
	if err != nil {
		panic(err)
	}
}

var labelReferenceFieldSpecs FieldSpecList

func labelsFilter(key string, value *string) *FieldSpecListFilter {
	if value != nil {
		return &FieldSpecListFilter{
			FieldSpecList: labelReferenceFieldSpecs,
			SetValue: func(node *yaml.RNode) error {
				_, err := node.Pipe(yaml.FieldSetter{Name: key, StringValue: *value})
				return err
			},
			CreateKind: yaml.MappingNode,
		}
	}

	return &FieldSpecListFilter{
		FieldSpecList: labelReferenceFieldSpecs,
		SetValue: func(node *yaml.RNode) error {
			_, err := node.Pipe(yaml.FieldClearer{Name: key})
			return err
		},
		CreateKind: yaml.MappingNode,
	}
}

const commonLabelFieldSpecs = `
items:
# duck-type supported labels
#
- path: metadata/labels
  create: true
- path: spec/template/metadata/labels
  create: false

# non-duck-type supported labels
#
- path: spec/template/metadata/labels
  create: true
  version: v1
  kind: ReplicationController

- path: spec/template/metadata/labels
  create: true
  kind: Deployment

- path: spec/template/metadata/labels
  create: true
  kind: ReplicaSet

- path: spec/template/metadata/labels
  create: true
  kind: DaemonSet

- path: spec/template/metadata/labels
  create: true
  group: apps
  kind: StatefulSet

- path: spec/volumeClaimTemplates[]/metadata/labels
  create: true
  group: apps
  kind: StatefulSet

- path: spec/template/metadata/labels
  create: true
  group: batch
  kind: Job

- path: spec/jobTemplate/metadata/labels
  create: true
  group: batch
  kind: CronJob

- path: spec/jobTemplate/spec/template/metadata/labels
  create: true
  group: batch
  kind: CronJob
`
