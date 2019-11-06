// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package fieldspec

import (
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type KustomizeNamespaceFilter struct {
	// commonNamespace is set on metadata.namespace for all Resources
	KustomizeNamespace string `yaml:"commonNamespace,omitempty"`
}

var _ kio.Filter = KustomizeNamespaceFilter{}

func (ns KustomizeNamespaceFilter) Filter(input []*yaml.RNode) ([]*yaml.RNode, error) {
	f := namespaceFilter(ns.KustomizeNamespace)
	_, err := kio.FilterAll(f).Filter(input)
	return input, err
}

func init() {
	err := yaml.Unmarshal([]byte(namespaceFieldSpecs), &namespaceReferenceFieldSpecs)
	if err != nil {
		panic(err)
	}
}

func namespaceFilter(value string) *FieldSpecListFilter {
	return &FieldSpecListFilter{
		FieldSpecList: namespaceReferenceFieldSpecs,
		SetValue: func(node *yaml.RNode) error {
			_, err := node.Pipe(yaml.FieldSetter{StringValue: value})
			return err
		},
		CreateKind: yaml.ScalarNode,
	}
}

var namespaceReferenceFieldSpecs FieldSpecList

const namespaceFieldSpecs = `
items:
- path: metadata/namespace
  create: true
`
