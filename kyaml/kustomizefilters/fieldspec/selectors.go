// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package fieldspec

import (
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// KustomizeSelectorsFilter sets selectors for all Resources in a package.
// Overrides existing selector values iff the keys match.
type KustomizeSelectorsFilter struct {
	// commonSelectors are the selectors to set
	Selectors map[string]*string `yaml:"commonSelectors,omitempty"`
}

var _ kio.Filter = KustomizeSelectorsFilter{}

func (sf KustomizeSelectorsFilter) Filter(input []*yaml.RNode) ([]*yaml.RNode, error) {
	for k, v := range sf.Selectors {
		f := selectorsFilter(k, v)
		_, err := kio.FilterAll(f).Filter(input)
		if err != nil {
			return nil, err
		}
	}
	return input, nil
}

func init() {
	err := yaml.Unmarshal([]byte(commonSelectorFieldSpecs), &selectorReferenceFieldSpecs)
	if err != nil {
		panic(err)
	}
}

func selectorsFilter(key string, value *string) *FieldSpecListFilter {
	if value != nil {
		return &FieldSpecListFilter{
			FieldSpecList: selectorReferenceFieldSpecs,
			SetValue: func(node *yaml.RNode) error {
				_, err := node.Pipe(yaml.FieldSetter{Name: key, StringValue: *value})
				return err
			},
			CreateKind: yaml.MappingNode,
		}
	}

	return &FieldSpecListFilter{
		FieldSpecList: selectorReferenceFieldSpecs,
		SetValue: func(node *yaml.RNode) error {
			_, err := node.Pipe(yaml.FieldClearer{Name: key})
			return err
		},
		CreateKind: yaml.MappingNode,
	}
}

var selectorReferenceFieldSpecs FieldSpecList

const commonSelectorFieldSpecs = `
items:

# duck-type supported selectors
#
- path: spec/selector/matchLabels
  create: false

# non-duck-type supported selectors
#
- path: spec/selector/matchLabels
  create: true
  group: apps
  kind: StatefulSet
- path: spec/selector/matchLabels
  create: true
  kind: DaemonSet
- path: spec/selector/matchLabels
  create: true
  kind: ReplicaSet
- path: spec/selector/matchLabels
  create: true
  kind: Deployment
- path: spec/selector
  create: true
  version: v1
  kind: Service
- path: spec/selector
  create: true
  version: v1
  kind: ReplicationController
- path: spec/jobTemplate/spec/selector/matchLabels
  create: false
  group: batch
  kind: CronJob
- path: spec/podSelector/matchLabels
  create: false
  group: networking.k8s.io
  kind: NetworkPolicy
- path: spec/ingress/from/podSelector/matchLabels
  create: false
  group: networking.k8s.io
  kind: NetworkPolicy
- path: spec/egress/to/podSelector/matchLabels
  create: false
  group: networking.k8s.io
  kind: NetworkPolicy
`
