package pusher

import (
	"github.com/aura-studio/dynamic-cli/config"
)

// PushFromRepo pushes a package from repo
func PushFromRepo(repo string, packages ...string) {
	configs := config.ParseRepo(repo, packages...)
	for _, config := range configs {
		fileList := NewTaskList(config)
		New(fileList).Push()
	}
}

// PushFromJSON pushes a package from json string
func PushFromJSON(str string) {
	configs := config.ParseJSON(str)
	for _, config := range configs {
		fileList := NewTaskList(config)
		New(fileList).Push()
	}
}

// BuildFromJSONFile pushes a package from json file
func PushFromJSONFile(path string) {
	configs := config.ParseJSONFile(path)
	for _, config := range configs {
		fileList := NewTaskList(config)
		New(fileList).Push()
	}
}

// BuildFromJSONPath pushes a package from json path
func PushFromJSONPath(path string) {
	configs := config.ParseJSONPath(path)
	for _, config := range configs {
		fileList := NewTaskList(config)
		New(fileList).Push()
	}
}
