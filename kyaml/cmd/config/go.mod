module sigs.k8s.io/kustomize/cmd/config

go 1.13

require (
	github.com/go-errors/errors v1.0.1
	github.com/go-openapi/spec v0.19.5
	github.com/olekukonko/tablewriter v0.0.4
	github.com/posener/complete/v2 v2.0.1-alpha.12
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.4.0
	gopkg.in/inf.v0 v0.9.1
	sigs.k8s.io/kustomize/kyaml v0.1.7
)

replace sigs.k8s.io/kustomize/kyaml v0.1.7 => ../../kyaml
