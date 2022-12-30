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
	Module    string   `json:"module"`
	Debug     bool     `json:"debug"`
	Namespace string   `json:"namespace"`
	Commit    string   `json:"commit"`
	Packages  []string `json:"packages"`
	WareHouse string   `json:"warehouse"`
	NetRC     string   `json:"netrc"`
	Remotes   []string `json:"remotes"`
}

var DefaultConfig = Config{
	GoVer:     "1.18",
	WareHouse: "/tmp/warehouse",
	Debug:     false,
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

// ParseRepo parses a remote string into struct
// example: github.com/aura-studio/dynamic-cli/builder@af3e5e21
func ParseRepo(remote string, packages ...string) []Config {
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
