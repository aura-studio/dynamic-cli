package builder

import "github.com/aura-studio/dynamic-cli/config"

// BuildFromRepo builds a package from remote
func BuildFromRepo(repo string, packages ...string) {
	configs := config.ParseRepo(repo, packages...)
	for _, config := range configs {
		renderDatas := ConfigToRenderData(config)
		for _, renderData := range renderDatas {
			New(renderData).Build()
		}
	}
}

// BuildFromJSON builds a package from json string
func BuildFromJSON(str string) {
	configs := config.ParseJSON(str)
	for _, config := range configs {
		renderDatas := ConfigToRenderData(config)
		for _, renderData := range renderDatas {
			New(renderData).Build()
		}
	}
}

// BuildFromJSONFile builds a package from json file
func BuildFromJSONFile(path string) {
	configs := config.ParseJSONFile(path)
	for _, config := range configs {
		renderDatas := ConfigToRenderData(config)
		for _, renderData := range renderDatas {
			New(renderData).Build()
		}
	}
}

// BuildFromJSONPath builds a package from json path
func BuildFromJSONPath(path string) {
	configs := config.ParseJSONPath(path)
	for _, config := range configs {
		renderDatas := ConfigToRenderData(config)
		for _, renderData := range renderDatas {
			New(renderData).Build()
		}
	}
}
