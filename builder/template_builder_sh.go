package builder

func init() {
	templateMap["{{.House}}/{{.Name}}_{{.Version}}/builder.sh"] = builderBash
}

const builderBash = `#!/bin/sh
export CGO_ENABLED=1
export GO111MODULE=on
export GOPRIVATE={{.Module}}
go mod tidy
go build -o {{.House}}/{{.Name}}_{{.Version}}/libcgo_{{.Name}}_{{.Version}}.so -buildmode=c-shared {{.House}}/{{.Name}}_{{.Version}}/libcgo_{{.Name}}_{{.Version}}
go build -o {{.House}}/{{.Name}}_{{.Version}}/libgo_{{.Name}}_{{.Version}}.so -buildmode=plugin -ldflags="-r {{.House}}/{{.Name}}_{{.Version}}/" {{.House}}/{{.Name}}_{{.Version}}/libgo_{{.Name}}_{{.Version}}
`
