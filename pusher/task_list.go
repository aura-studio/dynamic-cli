package pusher

import (
	"fmt"
	"strings"

	"github.com/aura-studio/dynamic-cli/config"
)

type TaskList struct {
	tasks map[string][]string
}

func (f *TaskList) Add(remote string, file string) {
	f.tasks[remote] = append(f.tasks[remote], file)
}

func NewTaskList(remotes []string, c config.Config) *TaskList {
	fileList := &TaskList{}

	if len(c.Packages) == 0 {
		pkg := c.Module[strings.LastIndex(c.Module, "/")+1:]
		name := pkg
		if len(c.Namespace) > 0 {
			name = strings.Join([]string{c.Namespace, pkg}, "_")
		}
		libcgo := fmt.Sprintf("%s/%s_%s/libcgo_%s_%s.so", c.WareHouse, name, c.Commit, name, c.Commit)
		libgo := fmt.Sprintf("%s/%s_%s/libgo_%s_%s.so", c.WareHouse, name, c.Commit, name, c.Commit)

		for _, remote := range remotes {
			fileList.Add(remote, libcgo)
			fileList.Add(remote, libgo)
		}
		return fileList
	}

	for _, pkg := range c.Packages {
		name := pkg[strings.LastIndex(pkg, "/")+1:]
		if len(c.Namespace) > 0 {
			name = strings.Join([]string{c.Namespace, name}, "_")
		}
		libcgo := fmt.Sprintf("%s/%s_%s/libcgo_%s_%s.so", c.WareHouse, name, c.Commit, name, c.Commit)
		libgo := fmt.Sprintf("%s/%s_%s/libgo_%s_%s.so", c.WareHouse, name, c.Commit, name, c.Commit)
		for _, remote := range remotes {
			fileList.Add(remote, libcgo)
			fileList.Add(remote, libgo)
		}
	}

	return fileList
}
