package clean

import (
	"path/filepath"

	"github.com/aura-studio/dynamic-cli/config"
)

type PathList struct {
	WareHouse string
	Dirs      []string
	Files     []string
}

func (f *PathList) AddDir(dir string) {
	f.Dirs = append(f.Dirs, dir)
}

func (f *PathList) AddFile(file string) {
	f.Files = append(f.Files, file)
}

// NewPathListForProcedure constructs paths aligned with Procedure-based build outputs.
// Dir = Warehouse.Local / OS_Arch_Compiler_Variant / Namespace_Package_Version
// Files include libcgo_<name>.so and libgo_<name>.so under Dir.
func NewPathListForProcedure(proc config.Procedure) *PathList {
	name := proc.Target.Namespace + "_" + proc.Target.Package + "_" + proc.Target.Version
	env := proc.Toolchain.OS + "_" + proc.Toolchain.Arch + "_" + proc.Toolchain.Compiler + "_" + proc.Toolchain.Variant
	dir := filepath.Join(proc.Warehouse.Local, env, name)

	pl := &PathList{
		WareHouse: proc.Warehouse.Local,
		Dirs:      []string{dir},
		Files: []string{
			filepath.Join(dir, "libcgo_"+name+".so"),
			filepath.Join(dir, "libgo_"+name+".so"),
			filepath.Join(dir, "meta_"+name+".json"),
		},
	}
	return pl
}
