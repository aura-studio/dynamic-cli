package config

import (
	"fmt"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Environments []struct {
		Name      string `yaml:"name"`
		Toolchain struct {
			OS       string `yaml:"os"`
			Arch     string `yaml:"arch"`
			Compiler string `yaml:"compiler"`
			Variant  string `yaml:"variant"`
		} `yaml:"toolchain"`
		Warehouse struct {
			Local  string   `yaml:"local"`
			Remote []string `yaml:"remote"`
		} `yaml:"warehouse"`
	} `yaml:"environments"`

	Procedures []struct {
		Name        string `yaml:"name"`
		Environment string `yaml:"environment"`
		Source      struct {
			Module  string `yaml:"module"`
			Package string `yaml:"package"`
			Version string `yaml:"version"`
		} `yaml:"source"`
		Target struct {
			Namespace string `yaml:"namespace"`
			Package   string `yaml:"package"`
			Version   string `yaml:"version"`
		} `yaml:"target"`
	} `yaml:"procedures"`
}

// Parse reads the given YAML file and returns Config
func Parse(file string) Config {
	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		panic(err)
	}
	fmt.Println(config)
	return config
}

// Validate checks Config fields and panics on failures.
// Rules:
// 1. All fields must be non-empty.
// 2. Each procedure.environment must exist in environments names.
// 3. Toolchain fields (os, arch, compiler, variant) must match ^[A-Za-z0-9.-]+$.
// 4. Target values (namespace, package, version) must match ^[A-Za-z0-9.-]+$.
func Validate(c Config) {
	// regex for allowed characters: letters, digits, dash, dot
	allowed := regexp.MustCompile(`^[A-Za-z0-9.-]+$`)

	// collect environment names for existence checks
	envNames := map[string]struct{}{}
	if len(c.Environments) == 0 {
		panic("config: environments must not be empty")
	}
	for i, env := range c.Environments {
		if env.Name == "" {
			panic("config: environments[" + itoa(i) + "].name must not be empty")
		}
		// toolchain
		if env.Toolchain.OS == "" || env.Toolchain.Arch == "" || env.Toolchain.Compiler == "" || env.Toolchain.Variant == "" {
			panic("config: environments[" + itoa(i) + "] toolchain fields must not be empty")
		}
		if !allowed.MatchString(env.Toolchain.OS) || !allowed.MatchString(env.Toolchain.Arch) || !allowed.MatchString(env.Toolchain.Compiler) || !allowed.MatchString(env.Toolchain.Variant) {
			panic("config: environments[" + itoa(i) + "].toolchain fields contain invalid characters (allowed: letters, digits, '.', '-')")
		}
		// warehouse
		if env.Warehouse.Local == "" {
			panic("config: environments[" + itoa(i) + "].warehouse.local must not be empty")
		}
		if len(env.Warehouse.Remote) == 0 {
			panic("config: environments[" + itoa(i) + "].warehouse.remote must not be empty")
		}
		for j, r := range env.Warehouse.Remote {
			if r == "" {
				panic("config: environments[" + itoa(i) + "].warehouse.remote[" + itoa(j) + "] must not be empty")
			}
		}
		envNames[env.Name] = struct{}{}
	}

	if len(c.Procedures) == 0 {
		panic("config: procedures must not be empty")
	}
	for i, p := range c.Procedures {
		if p.Name == "" {
			panic("config: procedures[" + itoa(i) + "].name must not be empty")
		}
		if p.Environment == "" {
			panic("config: procedures[" + itoa(i) + "].environment must not be empty")
		}
		if _, ok := envNames[p.Environment]; !ok {
			panic("config: procedures[" + itoa(i) + "].environment not found in environments")
		}
		// source
		if p.Source.Module == "" || p.Source.Package == "" || p.Source.Version == "" {
			panic("config: procedures[" + itoa(i) + "] source fields must not be empty")
		}
		// target validation: only allowed characters
		if p.Target.Namespace == "" || p.Target.Package == "" || p.Target.Version == "" {
			panic("config: procedures[" + itoa(i) + "] target fields must not be empty")
		}
		if !allowed.MatchString(p.Target.Namespace) {
			panic("config: procedures[" + itoa(i) + "].target.namespace contains invalid characters")
		}
		if !allowed.MatchString(p.Target.Package) {
			panic("config: procedures[" + itoa(i) + "].target.package contains invalid characters")
		}
		if !allowed.MatchString(p.Target.Version) {
			panic("config: procedures[" + itoa(i) + "].target.version contains invalid characters")
		}
	}
}

// itoa for small ints without importing strconv to keep minimal deps here.
func itoa(i int) string {
	// simple conversion for non-negative indices
	digits := []byte{}
	if i == 0 {
		return "0"
	}
	for i > 0 {
		d := byte('0' + (i % 10))
		digits = append([]byte{d}, digits...)
		i /= 10
	}
	return string(digits)
}
