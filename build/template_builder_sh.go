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

# Meta variables
meta_module="{{.Module}}"
meta_version=$(go list -m -f '{{"{{"}}.Version{{"}}"}}' "$meta_module" 2>/dev/null || echo "unknown")
meta_built=$(TZ='Asia/Shanghai' date "+%Y-%m-%d_%H:%M:%S_CST%z")
meta_os="{{.OS}}"
meta_arch="{{.Arch}}"
meta_compiler="{{.Compiler}}"
meta_variant="{{.Variant}}"

# Build shared library with CGO enabled.
cgo_ldflags="-s -w -extldflags=-Wl,-rpath,\$ORIGIN -X 'main.MetaModule=$meta_module' -X 'main.MetaVersion=$meta_version' -X 'main.MetaBuilt=$meta_built' -X 'main.MetaOS=$meta_os' -X 'main.MetaArch=$meta_arch' -X 'main.MetaCompiler=$meta_compiler' -X 'main.MetaVariant=$meta_variant'"
go_ldflags="-s -w -extldflags=-Wl,-rpath,\$ORIGIN"
go build -trimpath -o {{.Dir}}/libcgo_{{.Name}}.so -buildvcs=false -buildmode=c-shared -ldflags="$cgo_ldflags" {{.Dir}}/libcgo_{{.Name}}
go build -trimpath -o {{.Dir}}/libgo_{{.Name}}.so -buildvcs=false -buildmode=plugin -ldflags="$go_ldflags" {{.Dir}}/libgo_{{.Name}}
cp -rf {{.Dir}}/libgo_{{.Name}}.so {{.Dir}}/libgo_{{.Name}}.so.$meta_built
cp -rf {{.Dir}}/libcgo_{{.Name}}.so {{.Dir}}/libcgo_{{.Name}}.so.$meta_built

# Generate meta file for this libcgo/libgo pair (for backward compatibility).
meta={{.Dir}}/meta_{{.Name}}.json
cat > "$meta" <<EOF
{
  "module": "$meta_module",
  "version": "$meta_version",
  "built": "$meta_built",
  "os": "$meta_os",
  "arch": "$meta_arch",
  "compiler": "$meta_compiler",
  "variant": "$meta_variant"
}
EOF

# Backup meta file with the same timestamp suffix.
cp -rf "$meta" "$meta.$meta_built"
`
