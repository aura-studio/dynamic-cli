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
	// 1. 清理本次构建目录下的杂质
	for _, dir := range c.Dirs {
		if _, err := os.Stat(dir); err != nil {
			if os.IsNotExist(err) {
				continue
			} else {
				log.Panic(err)
			}
		} else {
			if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if path == dir {
					return nil
				}
				if info.IsDir() {
					log.Printf("clean remove %s", path)
					if err := os.RemoveAll(path); err != nil {
						log.Panic(err)
					}
					return filepath.SkipDir
				}
				isTarget := false
				for _, file := range c.Files {
					if path == file {
						isTarget = true
						break
					}
				}
				if !isTarget {
					log.Printf("clean remove %s", path)
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

	// 2. 清理 WareHouse 下所有历史版本的 so 和 json (不属于本次构建的)
	if _, err := os.Stat(c.WareHouse); err == nil {
		if err := filepath.Walk(c.WareHouse, func(path string, info os.FileInfo, err error) error {
			if path == c.WareHouse || info.IsDir() {
				return nil
			}
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".so" || ext == ".json" {
				isCurrent := false
				for _, file := range c.Files {
					if path == file {
						isCurrent = true
						break
					}
				}
				if !isCurrent {
					log.Printf("clean remove historical artifact %s", path)
					if err := os.Remove(path); err != nil {
						log.Panic(err)
					}
				}
			}
			return nil
		}); err != nil {
			log.Panic(err)
		}
	}
}

func (c *Cleaner) cleanCache() {
	if _, err := os.Stat(c.WareHouse); err != nil {
		if os.IsNotExist(err) {
			return
		} else {
			log.Panic(err)
		}
	} else {
		if err := filepath.Walk(c.WareHouse, func(path string, info os.FileInfo, err error) error {
			if path == c.WareHouse {
				return nil
			}
			if info.IsDir() {
				return nil
			}
			ext := strings.ToLower(filepath.Ext(path))
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
