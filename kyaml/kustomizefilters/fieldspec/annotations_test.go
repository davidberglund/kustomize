// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package fieldspec_test

import (
	"testing"

	"sigs.k8s.io/kustomize/kyaml/kustomizefilters/fieldspec"
)

func TestAnnotationFilter(t *testing.T) {
	doTestCases(t, annotationTestCases)
}

var bar = "bar"

var annotationTestCases = []fieldSpecTestCase{

	// Test Case
	{
		name: "add-annotations-crd",
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
  annotations:
    # keep this annotation
    a: b
spec:
  template:
    metadata:
      annotations:
        c: d
---
apiVersion: example.com/v1
kind: Foo
metadata:
  name: instance-2
  annotations: {}
`,
		expected: `
apiVersion: example.com/v1
kind: Foo
metadata:
  name: instance
  annotations:
    foo: bar
---
apiVersion: example.com/v1
kind: Bar
metadata:
  name: instance
  annotations:
    # keep this annotation
    a: b
    foo: bar
spec:
  template:
    metadata:
      annotations:
        c: d
        foo: bar
---
apiVersion: example.com/v1
kind: Foo
metadata:
  name: instance-2
  annotations:
    foo: bar
`,
		instance: fieldspec.KustomizeAnnotationsFilter{Annotations: map[string]*string{"foo": &bar}},
	},

	// Test Case
	{
		name: "update-annotation-crd",
		input: `
apiVersion: example.com/v1
kind: Foo
metadata:
  name: instance
  annotations:
    foo: baz
---
apiVersion: example.com/v1
kind: Bar
metadata:
  name: instance
  annotations:
    foo: baz
    a: b
spec:
  template:
    metadata:
      annotations:
        c: d
        foo: baz
`,
		expected: `
apiVersion: example.com/v1
kind: Foo
metadata:
  name: instance
  annotations:
    foo: bar
---
apiVersion: example.com/v1
kind: Bar
metadata:
  name: instance
  annotations:
    foo: bar
    a: b
spec:
  template:
    metadata:
      annotations:
        c: d
        foo: bar
`,
		instance: fieldspec.KustomizeAnnotationsFilter{Annotations: map[string]*string{"foo": &bar}},
	},

	// Test Case
	{
		name: "add-annotation-replication-controller",
		input: `
apiVersion: v1
kind: ReplicationController
metadata:
  name: instance
`,
		expected: `
apiVersion: v1
kind: ReplicationController
metadata:
  name: instance
  annotations:
    foo: bar
spec:
  template:
    metadata:
      annotations:
        foo: bar
`,
		instance: fieldspec.KustomizeAnnotationsFilter{Annotations: map[string]*string{"foo": &bar}},
	},

	// Test Case
	{
		name: "add-annotation-deployment",
		input: `
apiVersion: example.com/v1alpha17
kind: Deployment
metadata:
  name: instance
`,
		expected: `
apiVersion: example.com/v1alpha17
kind: Deployment
metadata:
  name: instance
  annotations:
    foo: bar
spec:
  template:
    metadata:
      annotations:
        foo: bar
`,
		instance: fieldspec.KustomizeAnnotationsFilter{Annotations: map[string]*string{"foo": &bar}},
	},

	// Test Case
	{
		name: "add-annotation-replica-set",
		input: `
apiVersion: example.com/v1alpha17
kind: ReplicaSet
metadata:
  name: instance
`,
		expected: `
apiVersion: example.com/v1alpha17
kind: ReplicaSet
metadata:
  name: instance
  annotations:
    foo: bar
spec:
  template:
    metadata:
      annotations:
        foo: bar
`,
		instance: fieldspec.KustomizeAnnotationsFilter{Annotations: map[string]*string{"foo": &bar}},
	},

	// Test Case
	{
		name: "add-annotation-stateful-set",
		input: `
apiVersion: example.com/v1alpha17
kind: StatefulSet
metadata:
  name: instance
`,
		expected: `
apiVersion: example.com/v1alpha17
kind: StatefulSet
metadata:
  name: instance
  annotations:
    foo: bar
spec:
  template:
    metadata:
      annotations:
        foo: bar
`,
		instance: fieldspec.KustomizeAnnotationsFilter{Annotations: map[string]*string{"foo": &bar}},
	},

	// Test Case
	{
		name: "add-annotation-daemon-set",
		input: `
apiVersion: example.com/v1alpha17
kind: DaemonSet
metadata:
  name: instance
`,
		expected: `
apiVersion: example.com/v1alpha17
kind: DaemonSet
metadata:
  name: instance
  annotations:
    foo: bar
spec:
  template:
    metadata:
      annotations:
        foo: bar
`,
		instance: fieldspec.KustomizeAnnotationsFilter{Annotations: map[string]*string{"foo": &bar}},
	},

	// Test Case
	{
		name: "add-annotation-job",
		input: `
apiVersion: batch/v1alpha17
kind: Job
metadata:
  name: instance
`,
		expected: `
apiVersion: batch/v1alpha17
kind: Job
metadata:
  name: instance
  annotations:
    foo: bar
spec:
  template:
    metadata:
      annotations:
        foo: bar
`,
		instance: fieldspec.KustomizeAnnotationsFilter{Annotations: map[string]*string{"foo": &bar}},
	},

	// Test Case
	{
		name: "add-annotation-job",
		input: `
apiVersion: batch/v1alpha17
kind: CronJob
metadata:
  name: instance
`,
		expected: `
apiVersion: batch/v1alpha17
kind: CronJob
metadata:
  name: instance
  annotations:
    foo: bar
spec:
  jobTemplate:
    metadata:
      annotations:
        foo: bar
    spec:
      template:
        metadata:
          annotations:
            foo: bar
`,
		instance: fieldspec.KustomizeAnnotationsFilter{Annotations: map[string]*string{"foo": &bar}},
	},
}
