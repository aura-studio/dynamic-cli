package cleaner

import "github.com/aura-studio/dynamic-cli/config"

func CleanFromRepo(all bool, repo string, packages ...string) {
	configs := config.ParseRepo(repo, packages...)
	for _, config := range configs {
		pathList := NewPathList(config)
		New(pathList).Clean(all)
	}
}

func CleanFromJSON(all bool, str string) {
	configs := config.ParseJSON(str)
	for _, config := range configs {
		pathList := NewPathList(config)
		New(pathList).Clean(all)
	}
}

func CleanFromJSONFile(all bool, path string) {
	configs := config.ParseJSONFile(path)
	for _, config := range configs {
		pathList := NewPathList(config)
		New(pathList).Clean(all)
	}
}

func CleanFromJSONDir(all bool, path string) {
	configs := config.ParseJSONPath(path)
	for _, config := range configs {
		pathList := NewPathList(config)
		New(pathList).Clean(all)
	}
}
