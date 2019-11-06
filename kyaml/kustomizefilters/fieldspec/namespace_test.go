// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package fieldspec_test

import (
	"testing"

	"sigs.k8s.io/kustomize/kyaml/kustomizefilters/fieldspec"
)

func TestNamespaceFilter(t *testing.T) {
	doTestCases(t, namespaceTestCases)
}

var namespaceTestCases = []fieldSpecTestCase{
	// Test Case
	{
		name: "add-namespace",
		input: `
apiVersion: example.com/v1
kind: Foo
metadata:
  name: instance
---
apiVersion: example.com/v1
kind: Bar
metadata:
  name: instance
`,
		expected: `
apiVersion: example.com/v1
kind: Foo
metadata:
  name: instance
  namespace: foo
---
apiVersion: example.com/v1
kind: Bar
metadata:
  name: instance
  namespace: foo
`,
		instance: fieldspec.KustomizeNamespaceFilter{KustomizeNamespace: "foo"},
	},

	// Test Case
	{
		name: "update-namespace",
		input: `
apiVersion: example.com/v1
kind: Foo
metadata:
  name: instance
  # update this namespace
  namespace: bar
---
apiVersion: example.com/v1
kind: Bar
metadata:
  name: instance
  namespace: bar
`,
		expected: `
apiVersion: example.com/v1
kind: Foo
metadata:
  name: instance
  # update this namespace
  namespace: foo
---
apiVersion: example.com/v1
kind: Bar
metadata:
  name: instance
  namespace: foo
`,
		instance: fieldspec.KustomizeNamespaceFilter{KustomizeNamespace: "foo"},
	},
}
