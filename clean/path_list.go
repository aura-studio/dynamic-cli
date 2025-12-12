package clean

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/aura-studio/dynamic-cli/config"
)

type PathList struct {
	WareHouse string
	Dirs      []string
	Files     []string
}

func (f *PathList) AddDir(dir string) {
	f.Dirs = append(f.Dirs, dir)
}

func (f *PathList) AddFile(file string) {
	f.Files = append(f.Files, file)
}

func NewPathList(c config.Config) *PathList {
	pathList := &PathList{
		WareHouse: c.WareHouse,
		Dirs:      make([]string, 0),
		Files:     make([]string, 0),
	}

	if len(c.Packages) == 0 {
		pkg := c.Path[strings.LastIndex(c.Path, "/")+1:]
		name := pkg
		if len(c.Namespace) > 0 {
			name = strings.Join([]string{c.Namespace, pkg}, "_")
		}

		dir := fmt.Sprintf("%s_%s", name, c.Ref)
		libcgo := fmt.Sprintf("%s_%s/libcgo_%s_%s.so", name, c.Ref, name, c.Ref)
		libgo := fmt.Sprintf("%s_%s/libgo_%s_%s.so", name, c.Ref, name, c.Ref)

		pathList.AddDir(filepath.Join(c.WareHouse, runtime.Version(), dir))
		pathList.AddFile(filepath.Join(c.WareHouse, runtime.Version(), libcgo))
		pathList.AddFile(filepath.Join(c.WareHouse, runtime.Version(), libgo))

		return pathList
	}

	for _, pkg := range c.Packages {
		name := pkg[strings.LastIndex(pkg, "/")+1:]
		if len(c.Namespace) > 0 {
			name = strings.Join([]string{c.Namespace, name}, "_")
		}

		dir := fmt.Sprintf("%s_%s", name, c.Ref)
		libcgo := fmt.Sprintf("%s_%s/libcgo_%s_%s.so", name, c.Ref, name, c.Ref)
		libgo := fmt.Sprintf("%s_%s/libgo_%s_%s.so", name, c.Ref, name, c.Ref)

		pathList.AddDir(filepath.Join(c.WareHouse, runtime.Version(), dir))
		pathList.AddFile(filepath.Join(c.WareHouse, runtime.Version(), libcgo))
		pathList.AddFile(filepath.Join(c.WareHouse, runtime.Version(), libgo))
	}

	return pathList
}
