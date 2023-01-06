package cleaner

import "github.com/aura-studio/dynamic-cli/config"

func CleanFromRepo(cleanType CleanType, repo string, packages ...string) {
	configs := config.ParseRepo(repo, packages...)
	for _, config := range configs {
		pathList := NewPathList(config)
		New(pathList).Clean(cleanType)
	}
}

func CleanFromJSON(cleanType CleanType, str string) {
	configs := config.ParseJSON(str)
	for _, config := range configs {
		pathList := NewPathList(config)
		New(pathList).Clean(cleanType)
	}
}

func CleanFromJSONFile(cleanType CleanType, path string) {
	configs := config.ParseJSONFile(path)
	for _, config := range configs {
		pathList := NewPathList(config)
		New(pathList).Clean(cleanType)
	}
}

func CleanFromJSONDir(cleanType CleanType, path string) {
	configs := config.ParseJSONPath(path)
	for _, config := range configs {
		pathList := NewPathList(config)
		New(pathList).Clean(cleanType)
	}
}
