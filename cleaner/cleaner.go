package cleaner

import (
	"log"
	"os"
	"path/filepath"
)

type Cleaner struct {
	*PathList
}

func New(pathList *PathList) *Cleaner {
	return &Cleaner{
		PathList: pathList,
	}
}

func (c *Cleaner) Clean(all bool) {
	if all {
		c.cleanDirs()
	} else {
		c.cleanFiles()
	}
}

func (c *Cleaner) cleanDirs() {
	for _, dir := range c.Dirs {
		if _, err := os.Stat(dir); err != nil {
			if os.IsNotExist(err) {
				continue
			} else {
				log.Panic(err)
			}
		} else {
			log.Printf("clean remove %s", dir)
			if err := os.Remove(dir); err != nil {
				log.Panic(err)
			}
		}
	}
}

func (c *Cleaner) cleanFiles() {
	for _, dir := range c.Dirs {
		if _, err := os.Stat(dir); err != nil {
			if os.IsNotExist(err) {
				continue
			} else {
				log.Panic(err)
			}
		} else {
			if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				var preserve bool
				for _, file := range c.Files {
					if path == file {
						preserve = true
					}
				}
				if !preserve {
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
}
