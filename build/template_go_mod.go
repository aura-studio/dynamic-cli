package build

func init() {
	templateMap["{{.Dir}}/go.mod"] = `module dynamicbuilder

go 1.18

require (
	{{.Module}} {{.Version}
}
