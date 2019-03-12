package config

import (
	"bbexgo/help"
	"bbexgo/log"
	"github.com/kylelemons/go-gypsy/yaml"
	"sync"
)

var lock *sync.Mutex = &sync.Mutex{}
var confile = make(map[string]*yaml.File)

func Get(key string, file ...string) string {
	var fn string
	if len(file) == 0 {
		fn = "conf"
	} else {
		fn = file[0]
		if len(fn) > 5 && string([]rune(fn)[len(fn)-5:]) == ".yaml" {
			fn = string([]rune(fn)[:len(fn)-5])
		}
	}
	currPath, err := help.GetCurrentPath()
	if err != nil {
		log.Fatal(err)
		return ""
	}

	if _, ok := confile[fn]; !ok {
		lock.Lock()
		defer lock.Unlock()
		if _, ok := confile[fn]; !ok {
			fi, err := yaml.ReadFile(currPath + "configs/" + fn + ".yaml")
			if err != nil {
				log.Fatal(err)
			}
			confile[fn] = fi
		}
	}
	res, err := confile[fn].Get(key)
	if err != nil {
		log.Warning(err)
		return ""
	}
	return res
}
