package build

func init() {
	templateMap["{{.Dir}}/go.mod"] = templateGoMod
}

const templateGoMod = `module dynamicbuilder

go 1.21

toolchain {{.Compiler}}

require (
	{{.Module}} {{.Version}}
)
`
