// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package fieldspec_test

import (
	"testing"

	"sigs.k8s.io/kustomize/kyaml/kustomizefilters/fieldspec"
)

func TestSelectorsFilter(t *testing.T) {
	doTestCases(t, selectorsTestCases)
}

var selectorsTestCases = []fieldSpecTestCase{

	// Test Case
	{
		name: "crd-spec-selector-matchLabels",
		input: `
apiVersion: batch/v1
kind: Foo
metadata:
  name: instance-1
spec:
  selector:
    matchLabels:
      a: b
---
apiVersion: batch/v1
kind: Foo
metadata:
  name: instance-2
spec:
  selector:
    matchLabels:
      e: d
---
apiVersion: batch/v1
kind: Foo
metadata:
  name: instance-3
`,
		expected: `
apiVersion: batch/v1
kind: Foo
metadata:
  name: instance-1
spec:
  selector:
    matchLabels:
      a: b
      e: f
---
apiVersion: batch/v1
kind: Foo
metadata:
  name: instance-2
spec:
  selector:
    matchLabels:
      e: f
---
apiVersion: batch/v1
kind: Foo
metadata:
  name: instance-3
`,
		instance: fieldspec.KustomizeSelectorsFilter{Selectors: map[string]*string{"e": &f}},
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
kind: DaemonSet
metadata:
  name: instance-3
---
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: instance-4
---
apiVersion: apps/v1
kind: Service
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
spec:
  selector:
    matchLabels:
      e: f
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: instance-2
spec:
  selector:
    matchLabels:
      e: f
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: instance-3
spec:
  selector:
    matchLabels:
      e: f
---
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: instance-4
spec:
  selector:
    matchLabels:
      e: f
---
apiVersion: apps/v1
kind: Service
metadata:
  name: instance-5
spec:
  selector:
    e: f
---
apiVersion: apps/v1
kind: ReplicationController
metadata:
  name: instance-6
spec:
  selector:
    e: f
`,
		instance: fieldspec.KustomizeSelectorsFilter{Selectors: map[string]*string{"e": &f}},
	},

	// Test Case
	{
		name: "cronjob",
		input: `
apiVersion: batch/v1
kind: CronJob
metadata:
  name: instance-1
spec:
  jobTemplate:
    spec:
      selector:
        matchLabels:
          a: b
---
apiVersion: batch/v1
kind: CronJob
metadata:
  name: instance-2
spec:
  jobTemplate:
    spec:
      selector:
        matchLabels:
          e: b
---
apiVersion: batch/v1
kind: CronJob
metadata:
  name: instance-3
`,
		expected: `
apiVersion: batch/v1
kind: CronJob
metadata:
  name: instance-1
spec:
  jobTemplate:
    spec:
      selector:
        matchLabels:
          a: b
          e: f
---
apiVersion: batch/v1
kind: CronJob
metadata:
  name: instance-2
spec:
  jobTemplate:
    spec:
      selector:
        matchLabels:
          e: f
---
apiVersion: batch/v1
kind: CronJob
metadata:
  name: instance-3
`,
		instance: fieldspec.KustomizeSelectorsFilter{Selectors: map[string]*string{"e": &f}},
	},

	// Test Case
	{
		name: "network-policy",
		input: `
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: instance-1
spec:
  podSelector:
    matchLabels: {}
  ingress:
    from:
      podSelector:
        matchLabels: {}
  egress:
    to:
      podSelector:
        matchLabels: {}
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: instance-2
`,
		expected: `
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: instance-1
spec:
  podSelector:
    matchLabels:
      e: f
  ingress:
    from:
      podSelector:
        matchLabels:
          e: f
  egress:
    to:
      podSelector:
        matchLabels:
          e: f
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: instance-2
`,
		instance: fieldspec.KustomizeSelectorsFilter{Selectors: map[string]*string{"e": &f}},
	},
}
