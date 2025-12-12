package config

type Build struct {
	Toolchain struct {
		OS       string `yaml:"os"`
		Arch     string `yaml:"arch"`
		Compiler string `yaml:"compiler"`
		Variant  string `yaml:"variant"`
	} `yaml:"Toolchain"`
	Warehouse struct {
		Local  string   `yaml:"local"`
		Remote []string `yaml:"remote"`
	} `yaml:"warehouse"`
	Source struct {
		Repo    string `yaml:"repo"`
		Module  string `yaml:"module"`
		Version string `yaml:"version"`
	} `yaml:"source"`
	Target struct {
		Namespace string `yaml:"namespace"`
		Package   string `yaml:"package"`
		Version   string `yaml:"version"`
	} `yaml:"target"`
}

// BuildForProcedure constructs a Build object from Config by procedure name.
// It looks up the procedure, then uses its environment to fill toolchain and warehouse.
// Panics if procedure or environment is not found or any required field is empty.
func BuildForProcedure(c Config, procedureName string) Build {
	if procedureName == "" {
		panic("build: procedure name must not be empty")
	}
	// find procedure
	var pIdx = -1
	for i, p := range c.Procedures {
		if p.Name == procedureName {
			pIdx = i
			break
		}
	}
	if pIdx < 0 {
		panic("build: procedure '" + procedureName + "' not found")
	}
	p := c.Procedures[pIdx]

	// find environment by name
	var eIdx = -1
	for i, e := range c.Environments {
		if e.Name == p.Environment {
			eIdx = i
			break
		}
	}
	if eIdx < 0 {
		panic("build: environment '" + p.Environment + "' not found for procedure '" + procedureName + "'")
	}
	e := c.Environments[eIdx]

	// compose Build
	var b Build
	b.Toolchain.OS = e.Toolchain.OS
	b.Toolchain.Arch = e.Toolchain.Arch
	b.Toolchain.Compiler = e.Toolchain.Compiler
	b.Toolchain.Variant = e.Toolchain.Variant

	b.Warehouse.Local = e.Warehouse.Local
	b.Warehouse.Remote = append(b.Warehouse.Remote, e.Warehouse.Remote...)

	b.Source.Repo = p.Source.Repo
	b.Source.Module = p.Source.Module
	b.Source.Version = p.Source.Version

	b.Target.Namespace = p.Target.Namespace
	b.Target.Package = p.Target.Package
	b.Target.Version = p.Target.Version

	return b
}
