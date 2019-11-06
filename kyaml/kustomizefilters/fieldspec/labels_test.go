// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package fieldspec_test

import (
	"testing"

	"sigs.k8s.io/kustomize/kyaml/kustomizefilters/fieldspec"
)

func TestLabelsFilter(t *testing.T) {
	doTestCases(t, labelsTestCases)
}

var f = "f"

var labelsTestCases = []fieldSpecTestCase{

	// Test Case
	{
		name: "crd",
		input: `
apiVersion: example.com/v1
kind: Foo
metadata:
  name: instance-1
spec:
  template:
    metadata:
      labels: {}
---
apiVersion: example.com/v1
kind: Foo
metadata:
  name: instance-2
`,
		expected: `
apiVersion: example.com/v1
kind: Foo
metadata:
  name: instance-1
  labels:
    e: f
spec:
  template:
    metadata:
      labels:
        e: f
---
apiVersion: example.com/v1
kind: Foo
metadata:
  name: instance-2
  labels:
    e: f
`,
		instance: fieldspec.KustomizeLabelsFilter{Labels: map[string]*string{"e": &f}},
	},

	// Test Case
	{
		name: "builtin-spec-selector-matchLabels",
		input: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: instance-1
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: instance-2
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: instance-3
spec:
  volumeClaimTemplates:
  - metadata:
      name: foo-1
  - metadata:
      name: foo-2
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: instance-4
---
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: instance-5
---
apiVersion: apps/v1
kind: ReplicationController
metadata:
  name: instance-6
`,
		expected: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: instance-1
  labels:
    e: f
spec:
  template:
    metadata:
      labels:
        e: f
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: instance-2
  labels:
    e: f
spec:
  template:
    metadata:
      labels:
        e: f
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: instance-3
  labels:
    e: f
spec:
  volumeClaimTemplates:
  - metadata:
      name: foo-1
      labels:
        e: f
  - metadata:
      name: foo-2
      labels:
        e: f
  template:
    metadata:
      labels:
        e: f
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: instance-4
  labels:
    e: f
spec:
  template:
    metadata:
      labels:
        e: f
---
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: instance-5
  labels:
    e: f
spec:
  template:
    metadata:
      labels:
        e: f
---
apiVersion: apps/v1
kind: ReplicationController
metadata:
  name: instance-6
  labels:
    e: f
spec:
  template:
    metadata:
      labels:
        e: f
`,
		instance: fieldspec.KustomizeLabelsFilter{Labels: map[string]*string{"e": &f}},
	},
}
