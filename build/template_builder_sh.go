package build

func init() {
	templateMap["{{.Dir}}/builder.sh"] = templateBuilderSh
}

const templateBuilderSh = `#!/bin/bash
set -e
mkdir -p {{.Dir}}
cd {{.Dir}}
export CGO_ENABLED=1
export GO111MODULE=on
export GOPRIVATE={{.Module}}
go clean --modcache
go mod tidy

# Meta variables - source
meta_source_module="{{.Module}}"
meta_source_version=$(go list -m -f '{{"{{"}}.Version{{"}}"}}' "$meta_source_module" 2>/dev/null || echo "unknown")
meta_source_built=$(TZ='Asia/Shanghai' date '+%Y-%m-%d_%H:%M:%S_CST+0800')

# Meta variables - toolchain
meta_toolchain_os="{{.OS}}"
meta_toolchain_arch="{{.Arch}}"
meta_toolchain_compiler="{{.Compiler}}"
meta_toolchain_variant="{{.Variant}}"

# Build shared library with CGO enabled.
cgo_ldflags="-s -w -extldflags=-Wl,-rpath,\$ORIGIN -X 'main.MetaSourceModule=$meta_source_module' -X 'main.MetaSourceVersion=$meta_source_version' -X 'main.MetaSourceBuilt=$meta_source_built' -X 'main.MetaToolchainOS=$meta_toolchain_os' -X 'main.MetaToolchainArch=$meta_toolchain_arch' -X 'main.MetaToolchainCompiler=$meta_toolchain_compiler' -X 'main.MetaToolchainVariant=$meta_toolchain_variant'"
go_ldflags="-s -w -extldflags=-Wl,-rpath,\$ORIGIN"
go build -trimpath -o {{.Dir}}/libcgo_{{.Name}}.so -buildvcs=false -buildmode=c-shared -ldflags="$cgo_ldflags" {{.Dir}}/libcgo_{{.Name}}
go build -trimpath -o {{.Dir}}/libgo_{{.Name}}.so -buildvcs=false -buildmode=plugin -ldflags="$go_ldflags" {{.Dir}}/libgo_{{.Name}}
cp -rf {{.Dir}}/libgo_{{.Name}}.so {{.Dir}}/libgo_{{.Name}}.so.$meta_source_built
cp -rf {{.Dir}}/libcgo_{{.Name}}.so {{.Dir}}/libcgo_{{.Name}}.so.$meta_source_built
`
