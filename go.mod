module github.com/shipyard-run/cli

go 1.13

require (
	github.com/DATA-DOG/godog v0.7.13
	github.com/docker/docker v1.13.1
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/hashicorp/hcl2 v0.0.0-20191002203319-fb75b3253c80
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.5.0
	github.com/stretchr/testify v1.2.2
	github.com/zclconf/go-cty v1.1.1
	golang.org/x/xerrors v0.0.0-20191011141410-1b5146add898
)

replace github.com/docker/docker => github.com/moby/moby v1.13.1
