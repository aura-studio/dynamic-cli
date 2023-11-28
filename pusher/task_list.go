package pusher

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/aura-studio/dynamic-cli/config"
)

type TaskList struct {
	Tasks map[string][]Pair
}

type Pair struct {
	RemoteFilePath string
	LocalFilePath  string
}

func (f *TaskList) Add(remote string, remoteFilePath string, localFilePath string) {
	f.Tasks[remote] = append(f.Tasks[remote], Pair{
		RemoteFilePath: remoteFilePath,
		LocalFilePath:  localFilePath,
	})
}

func NewTaskList(c config.Config) *TaskList {
	fileList := &TaskList{
		Tasks: make(map[string][]Pair),
	}

	if len(c.Packages) == 0 {
		pkg := c.Module[strings.LastIndex(c.Module, "/")+1:]
		name := pkg
		if len(c.Namespace) > 0 {
			name = strings.Join([]string{c.Namespace, pkg}, "_")
		}
		libcgo := fmt.Sprintf("%s_%s/libcgo_%s_%s.so", name, c.Commit, name, c.Commit)
		libgo := fmt.Sprintf("%s_%s/libgo_%s_%s.so", name, c.Commit, name, c.Commit)

		for _, remote := range c.Remotes {
			if err := filepath.WalkDir(c.WareHouse, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					return nil
				}
				if strings.Contains(path, libcgo) {
					fileList.Add(remote, libcgo, path)
				}
				if strings.Contains(path, libgo) {
					fileList.Add(remote, libgo, path)
				}
				return nil
			}); err != nil {
				panic(err)
			}
		}
		return fileList
	}

	for _, pkg := range c.Packages {
		name := pkg[strings.LastIndex(pkg, "/")+1:]
		if len(c.Namespace) > 0 {
			name = strings.Join([]string{c.Namespace, name}, "_")
		}
		libcgo := fmt.Sprintf("%s_%s/libcgo_%s_%s.so", name, c.Commit, name, c.Commit)
		libgo := fmt.Sprintf("%s_%s/libgo_%s_%s.so", name, c.Commit, name, c.Commit)
		for _, remote := range c.Remotes {
			fileList.Add(remote, libcgo, filepath.Join(c.WareHouse, libcgo))
			fileList.Add(remote, libgo, filepath.Join(c.WareHouse, libgo))
		}
	}

	return fileList
}
