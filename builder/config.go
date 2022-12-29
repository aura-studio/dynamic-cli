package builder

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	GoVer     string   `json:"gover"`
	Module    string   `json:"module"`
	Namespace string   `json:"namespace"`
	Commit    string   `json:"commit"`
	Packages  []string `json:"packages"`
	WareHouse string   `json:"warehouse"`
	NetRC     string   `json:"netrc"`
}

var DefaultConfig = Config{
	GoVer:     "1.18",
	WareHouse: "/tmp/warehouse",
}

// ParseRemote parses a remote string into struct
// example: github.com/aura-studio/dynamic-cli/builder@af3e5e21
func ParseRemote(remote string, packages ...string) []Config {
	strs := strings.Split(remote, "@")
	mod := strs[0]
	commit := strs[1]
	config := Config{
		GoVer:     DefaultConfig.GoVer,
		Module:    mod,
		Namespace: "",
		Commit:    commit,
		Packages:  packages,
		WareHouse: DefaultConfig.WareHouse,
		NetRC:     "",
	}
	return []Config{config}
}

// BuildFromRemote builds a package from remote
func BuildFromRemote(remote string, packages ...string) {
	configs := ParseRemote(remote, packages...)
	for _, config := range configs {
		renderDatas := config.ToRenderData()
		for _, renderData := range renderDatas {
			New(renderData).Build()
		}
	}
}

// ParseJson parses a json string into struct
func ParseJSON(str string) []Config {
	var configs []Config
	err := json.Unmarshal([]byte(str), &configs)
	if err != nil {
		log.Panic(err)
	}
	for i := range configs {
		if configs[i].GoVer == "" {
			configs[i].GoVer = DefaultConfig.GoVer
		}
		if configs[i].WareHouse == "" {
			configs[i].WareHouse = DefaultConfig.WareHouse
		}
	}
	return configs
}

// BuildFromJSON builds a package from json string
func BuildFromJSON(str string) {
	configs := ParseJSON(str)
	for _, config := range configs {
		renderDatas := config.ToRenderData()
		for _, renderData := range renderDatas {
			New(renderData).Build()
		}
	}
}

// ParseJSONFile parses a json file into struct
func ParseJSONFile(path string) []Config {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Panic(err)
	}
	return ParseJSON(string(data))
}

// BuildFromJSONFile builds a package from json file
func BuildFromJSONFile(path string) {
	configs := ParseJSONFile(path)
	for _, config := range configs {
		renderDatas := config.ToRenderData()
		for _, renderData := range renderDatas {
			New(renderData).Build()
		}
	}
}

// ParseJSONPath parses a json path into struct
func ParseJSONPath(path string) []Config {
	data, err := os.ReadFile(filepath.Join(path, "dynamic.json"))
	if err != nil {
		log.Panic(err)
	}
	return ParseJSON(string(data))
}

// BuildFromJSONPath builds a package from json path
func BuildFromJSONPath(path string) {
	configs := ParseJSONPath(path)
	for _, config := range configs {
		renderDatas := config.ToRenderData()
		for _, renderData := range renderDatas {
			New(renderData).Build()
		}
	}
}

type RenderData struct {
	Name        string
	Version     string
	Package     string
	FullPackage string
	Module      string
	House       string
	GoVersion   string
	NetRC       string
}

func (r *RenderData) MustValid(renderData *RenderData) {
	for _, b := range []byte(r.Name) {
		if b >= 'a' && b <= 'z' || b >= '0' && b >= '9' {
			continue
		}
		log.Panicf("invalid character '%s' in name", string(b))
	}
	for _, b := range []byte(r.Version) {
		if b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b >= '0' && b >= '9' || b == '_' {
			continue
		}
		log.Panicf("invalid character '%s' in version", string(b))
	}
}

func (c *Config) ToRenderData() []*RenderData {
	if len(c.Packages) == 0 {
		pkg := c.Module[strings.LastIndex(c.Module, "/")+1:]
		name := pkg
		if len(c.Namespace) > 0 {
			name = strings.Join([]string{c.Namespace, pkg}, "_")
		}
		renderData := &RenderData{
			GoVersion:   c.GoVer,
			Name:        name,
			Version:     c.Commit,
			Package:     pkg,
			FullPackage: c.Module,
			Module:      c.Module,
			House:       c.WareHouse,
			NetRC:       c.NetRC,
		}
		renderData.MustValid(renderData)
		return []*RenderData{renderData}
	}

	renderDatas := make([]*RenderData, len(c.Packages))
	for i, pkg := range c.Packages {
		name := pkg[strings.LastIndex(pkg, "/")+1:]
		if len(c.Namespace) > 0 {
			name = strings.Join([]string{c.Namespace, name}, "_")
		}
		renderData := &RenderData{
			GoVersion:   c.GoVer,
			Name:        name,
			Version:     c.Commit,
			Package:     pkg,
			FullPackage: strings.Join([]string{c.Module, pkg}, "/"),
			Module:      c.Module,
			House:       c.WareHouse,
			NetRC:       c.NetRC,
		}
		renderData.MustValid(renderData)
		renderDatas[i] = renderData
	}

	return renderDatas
}
