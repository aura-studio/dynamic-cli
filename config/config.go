package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	GoVer     string   `json:"gover"`
	Path      string   `json:"path"`
	Ref       string   `json:"ref"`
	NetRC     string   `json:"netrc"`
	Debug     bool     `json:"debug"`
	Namespace string   `json:"namespace"`
	Packages  []string `json:"packages"`
	WareHouse string   `json:"warehouse"`
	Remotes   []string `json:"remotes"`
}

var DefaultConfig = Config{
	GoVer:     "1.18",
	WareHouse: "/tmp/warehouse",
	Debug:     false,
	Remotes:   []string{},
	NetRC:     os.Getenv("DYNAMIC_CLI_NETRC"),
}

func SetDefaultGoVer(ver string) {
	DefaultConfig.GoVer = ver
}

func SetDefaultWareHouse(warehouse string) {
	DefaultConfig.WareHouse = warehouse
}

func SetDefaultDebug(debug bool) {
	DefaultConfig.Debug = debug
}

func SetDefaultRemotes(remotes []string) {
	DefaultConfig.Remotes = remotes
}

// ParseRepo parses a repo string into struct
// example: github.com/aura-studio/dynamic-cli/builder@af3e5e21
func ParseRepo(repo string, packages ...string) []Config {
	strs := strings.Split(repo, "@")
	mod := strs[0]
	ref := strs[1]
	config := Config{
		GoVer:     DefaultConfig.GoVer,
		Path:      mod,
		Namespace: "",
		Ref:       ref,
		Packages:  packages,
		WareHouse: DefaultConfig.WareHouse,
		NetRC:     DefaultConfig.NetRC,
		Remotes:   DefaultConfig.Remotes,
	}
	return []Config{config}
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
		if configs[i].Debug {
			configs[i].Debug = DefaultConfig.Debug
		}
		if len(configs[i].Remotes) == 0 {
			configs[i].Remotes = DefaultConfig.Remotes
		}
		if configs[i].NetRC == "" {
			configs[i].NetRC = DefaultConfig.NetRC
		}
	}
	return configs
}

// ParseJSONFile parses a json file into struct
func ParseJSONFile(path string) []Config {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Panic(err)
	}
	return ParseJSON(string(data))
}

// ParseJSONPath parses a json path into struct
func ParseJSONPath(path string) []Config {
	data, err := os.ReadFile(filepath.Join(path, "dynamic.json"))
	if err != nil {
		log.Panic(err)
	}
	return ParseJSON(string(data))
}
