// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	"sigs.k8s.io/kustomize/kyaml/cmd/config/complete"
	"sigs.k8s.io/kustomize/kyaml/cmd/config/configcobra"
	"sigs.k8s.io/kustomize/kyaml/commandutil"
)

func main() {
	// enable the config commands
	os.Setenv(commandutil.EnableAlphaCommmandsEnvName, "true")
	cmd := configcobra.NewConfigCommand("")
	complete.Complete(cmd).Complete("config")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
