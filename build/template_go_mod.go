package build

func init() {
	templateMap["{{.Dir}}/go.mod"] = templateGoMod
}

const templateGoMod = `module dynamicbuilder

go 1.18

require (
	{{.Module}} {{.Version}}
)
`
