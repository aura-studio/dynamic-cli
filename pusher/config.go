package pusher

import "github.com/aura-studio/dynamic-cli/config"

func PushFromRepo(repo string, packages ...string) {
	configs := config.ParseRepo(repo, packages...)
	for _, config := range configs {
		
	}
}

// // BuildFromJSONFile builds a package from json file
func PushFromJSONFile(path string) {
	configs := config.ParseJSONFile(path)
	for _, config := range configs {
		// renderDatas := config.ToRenderData()
		// for _, renderData := range renderDatas {
		// 	New(renderData).Build()
		// }
	}
}

// // ParseJSONPath parses a json path into struct
// func ParseJSONPath(path string) []Config {
// 	data, err := os.ReadFile(filepath.Join(path, "dynamic.json"))
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	return ParseJSON(string(data))
// }
