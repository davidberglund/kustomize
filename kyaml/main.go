// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/kustomize/kyaml/cmd"
)

var root = &cobra.Command{
	Use:   "kyaml",
	Short: "kyaml reference comand",
	Long: `Description:
  Reference implementation for using the kyaml libraries.
`,
	Example: ``,
}

func main() {
	root.AddCommand(cmd.GrepCommand())
	root.AddCommand(cmd.TreeCommand())
	root.AddCommand(cmd.CatCommand())
	root.AddCommand(cmd.FmtCommand())
	root.AddCommand(cmd.KustomizeCommand())
	root.AddCommand(cmd.MergeCommand())
	root.AddCommand(cmd.CountCommand())
	root.AddCommand(cmd.CreateBlueprintKustomizationCommand())

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
