package builder

func init() {
	templateMap["{{.House}}/{{.Name}}/go.mod"] = `module dynamicbuilder

go {{.GoVersion}}

require (
	{{.Module}} {{.Version}}
)
`
}
