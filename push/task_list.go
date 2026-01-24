package push

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

func NewTaskList(proc config.Procedure) *TaskList {
	fileList := &TaskList{Tasks: make(map[string][]Pair)}

	// Compose name/env/dir consistent with build
	name := proc.Target.Namespace + "_" + proc.Target.Package + "_" + proc.Target.Version
	env := proc.Toolchain.OS + "_" + proc.Toolchain.Arch + "_" + proc.Toolchain.Compiler + "_" + proc.Toolchain.Variant
	dir := filepath.Join(proc.Warehouse.Local, env, name)

	// Expected filenames under dir
	libcgoName := fmt.Sprintf("libcgo_%s.so", name)
	libgoName := fmt.Sprintf("libgo_%s.so", name)

	for _, remote := range proc.Warehouse.Remote {
		if err := filepath.WalkDir(proc.Warehouse.Local, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			// Only consider files inside the expected dir
			if !strings.HasPrefix(path, dir) {
				return nil
			}
			base := filepath.Base(path)
			// match primary and timestamped backups
			if base == libcgoName || strings.HasPrefix(base, libcgoName+".") ||
				base == libgoName || strings.HasPrefix(base, libgoName+".") {
				rel, relErr := filepath.Rel(proc.Warehouse.Local, path)
				if relErr != nil {
					return relErr
				}
				remotePath := filepath.ToSlash(rel)
				fileList.Add(remote, remotePath, path)
			}
			return nil
		}); err != nil {
			panic(err)
		}
	}

	return fileList
}
