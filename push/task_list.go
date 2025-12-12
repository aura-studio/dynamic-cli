package push

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"runtime"
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
		pkg := c.Path[strings.LastIndex(c.Path, "/")+1:]
		name := pkg
		if len(c.Namespace) > 0 {
			name = strings.Join([]string{c.Namespace, pkg}, "_")
		}
		libcgo := fmt.Sprintf("%s_%s/libcgo_%s_%s.so", name, c.Ref, name, c.Ref)
		libgo := fmt.Sprintf("%s_%s/libgo_%s_%s.so", name, c.Ref, name, c.Ref)

		for _, remote := range c.Remotes {
			if err := filepath.WalkDir(c.WareHouse, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					return nil
				}
				if strings.Contains(path, libcgo) {
					index := strings.Index(path, libcgo)
					fileList.Add(remote, filepath.Join(runtime.Version(), path[index:]), path)
				}
				if strings.Contains(path, libgo) {
					index := strings.Index(path, libgo)
					fileList.Add(remote, filepath.Join(runtime.Version(), path[index:]), path)
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
		libcgo := fmt.Sprintf("%s_%s/libcgo_%s_%s.so", name, c.Ref, name, c.Ref)
		libgo := fmt.Sprintf("%s_%s/libgo_%s_%s.so", name, c.Ref, name, c.Ref)
		for _, remote := range c.Remotes {
			if err := filepath.WalkDir(c.WareHouse, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					return nil
				}
				if strings.Contains(path, libcgo) {
					index := strings.Index(path, libcgo)
					fileList.Add(remote, filepath.Join(runtime.Version(), path[index:]), path)
				}
				if strings.Contains(path, libgo) {
					index := strings.Index(path, libgo)
					fileList.Add(remote, filepath.Join(runtime.Version(), path[index:]), path)
				}
				return nil
			}); err != nil {
				panic(err)
			}
		}
	}

	return fileList
}
