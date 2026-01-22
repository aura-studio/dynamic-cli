package clean

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type CleanType int

const (
	CleanTypeCache CleanType = iota
	CleanTypePackage
	CleanTypeAll
	CleanTypeUseless
)

type Cleaner struct {
	*PathList
}

func New(pathList *PathList) *Cleaner {
	return &Cleaner{
		PathList: pathList,
	}
}

func (c *Cleaner) Clean(cleanType CleanType) {
	switch cleanType {
	case CleanTypeCache:
		c.cleanCache()
	case CleanTypePackage:
		c.cleanPackage()
	case CleanTypeAll:
		c.cleanAll()
	case CleanTypeUseless:
		c.cleanUseless()
	default:
		log.Panicf("clean type not found: %d", cleanType)
	}
}

func (c *Cleaner) cleanAll() {
	if _, err := os.Stat(c.WareHouse); err != nil {
		if os.IsNotExist(err) {
			return
		} else {
			log.Panic(err)
		}
	} else {
		entries, err := os.ReadDir(c.WareHouse)
		if err != nil {
			log.Panic(err)
		}
		for _, file := range entries {
			if file.IsDir() {
				log.Printf("clean remove %s", filepath.Join(c.WareHouse, file.Name()))
				if err := os.RemoveAll(filepath.Join(c.WareHouse, file.Name())); err != nil {
					log.Panic(err)
				}
			} else {
				log.Printf("clean remove %s", filepath.Join(c.WareHouse, file.Name()))
				if err := os.Remove(filepath.Join(c.WareHouse, file.Name())); err != nil {
					log.Panic(err)
				}
			}
		}
	}
}

func (c *Cleaner) cleanPackage() {
	// clean package: 删除所有的 so 和 .json 及其带日期的后缀文件
	if _, err := os.Stat(c.WareHouse); err == nil {
		if err := filepath.Walk(c.WareHouse, func(path string, info os.FileInfo, err error) error {
			if path == c.WareHouse || info.IsDir() {
				return nil
			}
			name := strings.ToLower(info.Name())
			if strings.Contains(name, ".so") || strings.Contains(name, ".json") {
				log.Printf("clean remove artifact %s", path)
				if err := os.Remove(path); err != nil {
					log.Panic(err)
				}
			}
			return nil
		}); err != nil {
			log.Panic(err)
		}
	}
}

func (c *Cleaner) cleanCache() {
	// clean cache: 只保留所有的 so 和 .json 及其带日期的后缀文件
	if _, err := os.Stat(c.WareHouse); err != nil {
		if os.IsNotExist(err) {
			return
		} else {
			log.Panic(err)
		}
	} else {
		if err := filepath.Walk(c.WareHouse, func(path string, info os.FileInfo, err error) error {
			if path == c.WareHouse || info.IsDir() {
				return nil
			}
			name := strings.ToLower(info.Name())
			if strings.Contains(name, ".so") || strings.Contains(name, ".json") {
				return nil
			}
			log.Printf("clean remove non-artifact %s", path)
			if err := os.Remove(path); err != nil {
				log.Panic(err)
			}
			return nil
		}); err != nil {
			log.Panic(err)
		}
	}
}

func (c *Cleaner) cleanUseless() {
	// clean useless: 仅保留无日期后缀的 so 和 json 文件，其他文件都删除
	if _, err := os.Stat(c.WareHouse); err != nil {
		if os.IsNotExist(err) {
			return
		} else {
			log.Panic(err)
		}
	} else {
		if err := filepath.Walk(c.WareHouse, func(path string, info os.FileInfo, err error) error {
			if path == c.WareHouse || info.IsDir() {
				return nil
			}
			name := strings.ToLower(info.Name())
			ext := filepath.Ext(name)
			// 检查是否是 .so 或 .json 结尾，且不包含日期后缀
			if ext == ".so" || ext == ".json" {
				return nil
			}
			log.Printf("clean remove %s", path)
			if err := os.Remove(path); err != nil {
				log.Panic(err)
			}
			return nil
		}); err != nil {
			log.Panic(err)
		}
	}
}
