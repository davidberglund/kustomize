// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package kustomizefilters

import (
	"fmt"
	"sort"

	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/kustomizefilters/fieldreference"
	"sigs.k8s.io/kustomize/kyaml/kustomizefilters/fieldspec"
	"sigs.k8s.io/kustomize/kyaml/kustomizefilters/generators"
	"sigs.k8s.io/kustomize/kyaml/kustomizefilters/patches"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// BlueprintKustomizationFilter implements the BlueprintKustomization API
type BlueprintKustomizationFilter struct {
	KustomizeFilePath string

	yaml.ResourceMeta `yaml:",inline,omitempty"`

	Spec BlueprintKustomizationSpec `yaml:"spec,omitempty"`

	Status BlueprintKustomizationStatus `yaml:"status,omitempty"`
}

func (bk BlueprintKustomizationFilter) Filter(input []*yaml.RNode) ([]*yaml.RNode, error) {
	if bk.Name == "" {
		return nil, fmt.Errorf("must specify BlueprintKustomization metadata.name")
	}
	bk.Spec.SetKustomizeNameFilter.NameKustomizeName = bk.ObjectMeta.Name
	bk.Spec.ConfigMapGeneratorFilter.GeneratorKustomizeName = bk.ObjectMeta.Name
	bk.Spec.ConfigMapGeneratorFilter.Root = bk.KustomizeFilePath

	var in []*yaml.RNode
	var keep []*yaml.RNode
	var k *yaml.RNode
	for i := range input {
		meta, err := input[i].GetMeta()
		if err != nil {
			return nil, err
		}
		if meta.Kind == "BlueprintKustomization" && meta.Name == bk.Name {
			k = input[i]
			err = k.PipeE(yaml.Clear("status"))
			if err != nil {
				return nil, err
			}
			err := k.PipeE(
				yaml.SetAnnotation("config.kubernetes.io/local-config", "true"))
			if err != nil {
				return nil, err
			}
			keep = append(keep, k)
			continue
		}
		if _, found := meta.Annotations["config.kubernetes.io/local-config"]; found {
			keep = append(keep, input[i])
			continue
		}
		in = append(in, input[i])
	}

	reset := fieldreference.ResetKustomizeNameFilter{
		From:              map[fieldreference.NameKey]string{},
		NameKustomizeName: bk.Spec.SetKustomizeNameFilter.NameKustomizeName,
	}

	buff := &kio.PackageBuffer{Nodes: in}

	err := kio.Pipeline{
		Inputs: []kio.Reader{buff},
		Filters: []kio.Filter{
			&reset, // reset the names first
			&bk.Spec.ConfigMapGeneratorFilter,
			&bk.Spec.PatchFilter,
			&bk.Spec.KustomizeAnnotationsFilter,
			&bk.Spec.KustomizeLabelsFilter,
			&bk.Spec.KustomizeSelectorsFilter,
			&bk.Spec.KustomizeNamespaceFilter,
			&bk.Spec.SetKustomizeNameFilter,
		},
		Outputs: []kio.Writer{buff},
	}.Execute()
	if err != nil {
		return nil, err
	}

	// update status
	if k != nil {
		status := BlueprintKustomizationStatus{}
		for k, v := range bk.Spec.SetKustomizeNameFilter.To {
			status.NameMappings = append(status.NameMappings, Mapping{
				API:     k.Kind,
				NewName: v,
				Name:    k.Name,
			})
		}
		sort.Sort(status.NameMappings)

		b, err := yaml.Marshal(status)
		if err != nil {
			return nil, err
		}
		n := yaml.NewRNode(&yaml.Node{})
		err = yaml.Unmarshal(b, n.YNode())
		if err != nil {
			return nil, err
		}
		err = k.PipeE(yaml.SetField("status", n))
		if err != nil {
			return nil, err
		}
	}

	return append(buff.Nodes, keep...), nil
}

type BlueprintKustomizationSpec struct {
	patches.PatchFilter `yaml:",inline,omitempty"`

	fieldspec.KustomizeNamespaceFilter `yaml:",inline,omitempty"`

	fieldreference.SetKustomizeNameFilter `yaml:",inline,omitempty"`

	fieldspec.KustomizeLabelsFilter `yaml:",inline,omitempty"`

	fieldspec.KustomizeSelectorsFilter `yaml:",inline,omitempty"`

	fieldspec.KustomizeAnnotationsFilter `yaml:",inline,omitempty"`

	generators.ConfigMapGeneratorFilter `yaml:",inline,omitempty"`
}

type BlueprintKustomizationStatus struct {
	NameMappings Mappings `yaml:"nameMappings,omitempty"`
}

type Mappings []Mapping

// Len is part of sort.Interface.
func (s Mappings) Len() int {
	return len(s)
}

// Swap is part of sort.Interface.
func (s Mappings) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s Mappings) Less(i, j int) bool {
	if s[i].Name != s[j].Name {
		return s[i].Name < s[j].Name
	}
	if s[i].API != s[j].API {
		return s[i].API < s[j].API
	}
	if s[i].NewName != s[j].NewName {
		return s[i].NewName < s[j].NewName
	}
	return false
}

type Mapping struct {
	API     string `yaml:"api,omitempty"`
	NewName string `yaml:"newName,omitempty"`
	Name    string `yaml:"name,omitempty"`
}
