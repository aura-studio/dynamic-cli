package clean

import (
	"log"
	"os"
	"path/filepath"
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
	for _, dir := range c.Dirs {
		if _, err := os.Stat(dir); err != nil {
			if os.IsNotExist(err) {
				continue
			} else {
				log.Panic(err)
			}
		} else {
			log.Printf("clean remove %s", dir)
			if err := os.RemoveAll(dir); err != nil {
				log.Panic(err)
			}
		}
	}
}

func (c *Cleaner) cleanCache() {
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
				for _, file := range c.Files {
					if path == file {
						return nil
					}
				}
				log.Printf("clean remove %s", path)
				if err := os.RemoveAll(path); err != nil {
					log.Panic(err)
				}
				return nil
			}); err != nil {
				log.Panic(err)
			}
		}
	}
}
