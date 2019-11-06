// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package cmd_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/kustomize/kyaml/cmd"
)

// TODO(pwittrock): write tests for reading / writing ResourceLists

func TestCmd_files(t *testing.T) {
	d, err := ioutil.TempDir("", "kustomize-cat-test")
	if !assert.NoError(t, err) {
		return
	}
	defer os.RemoveAll(d)

	err = ioutil.WriteFile(filepath.Join(d, "f1.yaml"), []byte(`
kind: Deployment
metadata:
  labels:
    app: nginx2
  name: foo
  annotations:
    app: nginx2
spec:
  replicas: 1
---
kind: Service
metadata:
  name: foo
  annotations:
    app: nginx
spec:
  selector:
    app: nginx
`), 0600)
	if !assert.NoError(t, err) {
		return
	}
	err = ioutil.WriteFile(filepath.Join(d, "f2.yaml"), []byte(`
apiVersion: gcr.io/example/image:version
kind: Abstraction
metadata:
  name: foo
spec:
  replicas: 3
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: nginx
  name: bar
  annotations:
    app: nginx
spec:
  replicas: 3
`), 0600)
	if !assert.NoError(t, err) {
		return
	}

	// fmt the files
	b := &bytes.Buffer{}
	r := cmd.GetCatRunner()
	r.Command.SetArgs([]string{d})
	r.Command.SetOut(b)
	if !assert.NoError(t, r.Command.Execute()) {
		return
	}

	if !assert.Equal(t, `kind: Deployment
metadata:
  labels:
    app: nginx2
  name: foo
  annotations:
    app: nginx2
    config.kubernetess.io/package: .
    config.kubernetess.io/path: f1.yaml
spec:
  replicas: 1
---
kind: Service
metadata:
  name: foo
  annotations:
    app: nginx
    config.kubernetess.io/package: .
    config.kubernetess.io/path: f1.yaml
spec:
  selector:
    app: nginx
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: bar
  labels:
    app: nginx
  annotations:
    app: nginx
    config.kubernetess.io/package: .
    config.kubernetess.io/path: f2.yaml
spec:
  replicas: 3
`, b.String()) {
		return
	}
}

func TestCmd_filesWithReconcilers(t *testing.T) {
	d, err := ioutil.TempDir("", "kustomize-cat-test")
	if !assert.NoError(t, err) {
		return
	}
	defer os.RemoveAll(d)

	err = ioutil.WriteFile(filepath.Join(d, "f1.yaml"), []byte(`
kind: Deployment
metadata:
  labels:
    app: nginx2
  name: foo
  annotations:
    app: nginx2
spec:
  replicas: 1
---
kind: Service
metadata:
  name: foo
  annotations:
    app: nginx
spec:
  selector:
    app: nginx
`), 0600)
	if !assert.NoError(t, err) {
		return
	}
	err = ioutil.WriteFile(filepath.Join(d, "f2.yaml"), []byte(`
apiVersion: gcr.io/example/image:version
kind: Abstraction
metadata:
  name: foo
spec:
  replicas: 3
---
kind: Deployment
metadata:
  labels:
    app: nginx
  name: bar
  annotations:
    app: nginx
spec:
  replicas: 3
`), 0600)
	if !assert.NoError(t, err) {
		return
	}

	// fmt the files
	b := &bytes.Buffer{}
	r := cmd.GetCatRunner()
	r.Command.SetArgs([]string{d, "--include-reconcilers"})
	r.Command.SetOut(b)
	if !assert.NoError(t, r.Command.Execute()) {
		return
	}

	if !assert.Equal(t, `kind: Deployment
metadata:
  labels:
    app: nginx2
  name: foo
  annotations:
    app: nginx2
    config.kubernetess.io/package: .
    config.kubernetess.io/path: f1.yaml
spec:
  replicas: 1
---
kind: Service
metadata:
  name: foo
  annotations:
    app: nginx
    config.kubernetess.io/package: .
    config.kubernetess.io/path: f1.yaml
spec:
  selector:
    app: nginx
---
apiVersion: gcr.io/example/image:version
kind: Abstraction
metadata:
  name: foo
  annotations:
    config.kubernetess.io/package: .
    config.kubernetess.io/path: f2.yaml
spec:
  replicas: 3
---
kind: Deployment
metadata:
  labels:
    app: nginx
  name: bar
  annotations:
    app: nginx
    config.kubernetess.io/package: .
    config.kubernetess.io/path: f2.yaml
spec:
  replicas: 3
`, b.String()) {
		return
	}
}

func TestCmd_filesWithoutNonReconcilers(t *testing.T) {
	d, err := ioutil.TempDir("", "kustomize-cat-test")
	if !assert.NoError(t, err) {
		return
	}
	defer os.RemoveAll(d)

	err = ioutil.WriteFile(filepath.Join(d, "f1.yaml"), []byte(`
kind: Deployment
metadata:
  labels:
    app: nginx2
  name: foo
  annotations:
    app: nginx2
spec:
  replicas: 1
---
kind: Service
metadata:
  name: foo
  annotations:
    app: nginx
spec:
  selector:
    app: nginx
`), 0600)
	if !assert.NoError(t, err) {
		return
	}
	err = ioutil.WriteFile(filepath.Join(d, "f2.yaml"), []byte(`
apiVersion: gcr.io/example/image:version
kind: Abstraction
metadata:
  name: foo
spec:
  replicas: 3
---
kind: Deployment
metadata:
  labels:
    app: nginx
  name: bar
  annotations:
    app: nginx
spec:
  replicas: 3
`), 0600)
	if !assert.NoError(t, err) {
		return
	}

	// fmt the files
	b := &bytes.Buffer{}
	r := cmd.GetCatRunner()
	r.Command.SetArgs([]string{d, "--include-reconcilers", "--exclude-non-reconcilers"})
	r.Command.SetOut(b)
	if !assert.NoError(t, r.Command.Execute()) {
		return
	}

	if !assert.Equal(t, `apiVersion: gcr.io/example/image:version
kind: Abstraction
metadata:
  name: foo
  annotations:
    config.kubernetess.io/package: .
    config.kubernetess.io/path: f2.yaml
spec:
  replicas: 3
`, b.String()) {
		return
	}
}
