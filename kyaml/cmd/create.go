// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0
//
package cmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"sigs.k8s.io/kustomize/kyaml/kustomizefilters"
	"sigs.k8s.io/yaml"
)

func GetCreateBlueprintRunner() *CreateBlueprintKustomizationRunner {
	r := &CreateBlueprintKustomizationRunner{}
	c := &cobra.Command{
		Use:   "blueprint-kustomization DIR",
		Short: "Create a BlueprintKustomization for a directory",
		Long: `Create a BlueprintKustomization for a directory

  DIR:
    Path to local directory.
`,
		Example: `# create a new blueprint kustomization
kyaml create blueprint-kustomization my-dir/ --namespace=my-ns --annotations=foo=bar
`,
		RunE:    r.runE,
		PreRunE: r.preRunE,
		Args:    cobra.ExactArgs(1),
	}

	c.Flags().StringVar(&r.Name, "name", "", "name of the BlueprintKustomization")
	c.MarkFlagRequired("name")
	c.Flags().StringVar(&r.Spec.KustomizeNamespace, "namespace", "default", "common namespace")
	c.Flags().StringVar(&r.Spec.NamePrefix, "name-prefix", "", "common name prefix")
	c.Flags().StringVar(&r.Spec.NameSuffix, "name-suffix", "", "common name suffix")
	c.Flags().StringSliceVar(&r.Annotations, "annotation", []string{}, "common annotation")
	c.Flags().StringSliceVar(&r.Labels, "label", []string{}, "common label")
	c.Flags().StringSliceVar(&r.Selectors, "selector", []string{}, "common selector")

	r.Command = &cobra.Command{
		Use:   "create",
		Short: "Creation yaml files for common types.",
	}

	r.Command.AddCommand(c)
	return r
}

func CreateBlueprintKustomizationCommand() *cobra.Command {
	return GetCreateBlueprintRunner().Command
}

// CountRunner contains the run function
type CreateBlueprintKustomizationRunner struct {
	Command     *cobra.Command
	Annotations []string
	Labels      []string
	Selectors   []string
	kustomizefilters.BlueprintKustomizationFilter
}

func (r *CreateBlueprintKustomizationRunner) preRunE(c *cobra.Command, args []string) error {
	r.BlueprintKustomizationFilter.Kind = "BlueprintKustomization"
	r.BlueprintKustomizationFilter.ApiVersion = "kustomize.io/v1alpha1"

	for i := range r.Annotations {
		parts := strings.Split(r.Annotations[i], "=")
		if len(parts) != 2 {
			return fmt.Errorf("--annotation must have a value of the form 'key=value': %s",
				r.Annotations[i])
		}
		r.BlueprintKustomizationFilter.Spec.Annotations[parts[0]] = &parts[1]
	}
	for i := range r.Labels {
		parts := strings.Split(r.Labels[i], "=")
		if len(parts) != 2 {
			return fmt.Errorf("--label must have a value of the form 'key=value': %s",
				r.Labels[i])
		}
		r.BlueprintKustomizationFilter.Spec.Labels[parts[0]] = &parts[1]
	}

	for i := range r.Selectors {
		parts := strings.Split(r.Selectors[i], "=")
		if len(parts) != 2 {
			return fmt.Errorf("--label must have a value of the form 'key=value': %s",
				r.Labels[i])
		}
		r.BlueprintKustomizationFilter.Spec.Selectors[parts[0]] = &parts[1]
	}

	return nil
}

func (r *CreateBlueprintKustomizationRunner) runE(c *cobra.Command, args []string) error {
	b, err := yaml.Marshal(r.BlueprintKustomizationFilter)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(args[0], "kustomization.blueprint.yaml"), b, 0600)
}
