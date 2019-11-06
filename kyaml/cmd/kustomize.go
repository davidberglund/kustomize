// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/kio/filters"
	"sigs.k8s.io/kustomize/kyaml/kustomizefilters"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func GetKustomizeRunner() *KustomizeRunner {
	r := &KustomizeRunner{}
	c := &cobra.Command{
		Use:     "kustomize DIR",
		Short:   "Kustomize a package",
		Long:    ``,
		Example: ``,
		RunE:    r.runE,
	}
	c.Flags().BoolVar(&r.IncludeSubpackages, "include-subpackages", false,
		"also kustomize resources from subpackages.")
	c.Flags().BoolVar(&r.Format, "format", true,
		"format resource config yaml after kustomizing.")
	c.Flags().BoolVar(&r.KeepAnnotations, "annotate", false,
		"annotate resources with their file origins.")
	c.Flags().BoolVar(&r.DryRun, "dry-run", false,
		"print kustomizations rather than writing them back.")
	r.Command = c
	_ = r.Command.MarkFlagRequired("kustomization")
	r.Command.Args = cobra.ExactArgs(1)
	return r
}

func KustomizeCommand() *cobra.Command {
	return GetKustomizeRunner().Command
}

// KustomizeRunner contains the run function
type KustomizeRunner struct {
	IncludeSubpackages bool
	Format             bool
	KeepAnnotations    bool
	DryRun             bool
	Command            *cobra.Command
}

func (r *KustomizeRunner) runE(c *cobra.Command, args []string) error {
	var dirs []string
	err := filepath.Walk(args[0], func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		k := filepath.Join(path, "kustomization.yaml")
		if _, err := os.Stat(k); err == nil {
			dirs = append(dirs, k)
			return nil
		}
		return nil
	})
	if err != nil {
		return err
	}

	// kustomize depth-first
	for i := len(dirs) - 1; i >= 0; i-- {
		if err := r.kustomize(c, dirs[i]); err != nil {
			return err
		}
	}
	return nil
}

func (r *KustomizeRunner) kustomize(c *cobra.Command, path string) error {
	// setup input / output
	var rw kio.ReaderWriter
	var fltrs []kio.Filter
	// read package from fs
	rw = &kio.LocalPackageReadWriter{
		PackagePath:           filepath.Dir(path),
		IncludeSubpackages:    r.IncludeSubpackages,
		KeepReaderAnnotations: r.KeepAnnotations,
	}

	// setup kustomization filter
	k, err := getKustomization(path)
	if err != nil {
		return err
	}
	fltrs = append(fltrs, k)

	// setup format filter
	if r.Format {
		fltrs = append(fltrs, filters.FormatFilter{})
	}

	return kio.Pipeline{Inputs: []kio.Reader{rw}, Filters: fltrs, Outputs: []kio.Writer{rw}}.Execute()
}

func getKustomization(path string) (*kustomizefilters.BlueprintKustomizationFilter, error) {
	k := &kustomizefilters.BlueprintKustomizationFilter{KustomizeFilePath: filepath.Dir(path)}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	d := yaml.NewDecoder(bytes.NewBuffer(b))
	d.KnownFields(true)
	if err := d.Decode(k); err != nil {
		return nil, err
	}
	return k, nil
}
