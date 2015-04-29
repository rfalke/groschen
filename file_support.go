package groschen

import (
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
)

const (
	dirMode  = 0700
	fileMode = 0600
)

const (
	enoent    = iota
	directory = iota
	regular   = iota
)

func file_info_of(path string) int {
	fi, err := os.Stat(path)
	if err != nil {
		return enoent
	}
	if fi.IsDir() {
		return directory
	}
	if fi.Mode().IsRegular() {
		return regular
	}
	panic("Unknown mode")
}

func adjust_file_name(fname string) string {
	if strings.HasSuffix(fname, "/") {
		fname += "index.html"
		if file_info_of(fname) == directory {
			fname += ".alternative"
		}
	} else {
		if file_info_of(fname) == directory {
			fname += "/index.html"
		}
	}
	return fname
}

func ensure_dir(dir string) {
	info := file_info_of(dir)
	if info == enoent {
		parent := path.Dir(dir)
		ensure_dir(parent)
		os.Mkdir(dir, dirMode)
	} else if info == directory {
		return
	} else if info == regular {
		os.Rename(dir, dir+".tmp")
		os.Mkdir(dir, dirMode)
		os.Rename(dir+".tmp", dir+"/index.html")
	}
}

var fileLock sync.Mutex

func WriteResponseToFile(basePath string, content []byte, theUrl string) string {
	fileLock.Lock()
	defer fileLock.Unlock()

	u, err := url.Parse(theUrl)
	if err != nil {
		log.Fatal(err)
	}
	fname := basePath + "/" + u.Host + u.Path + u.RawQuery
	fname2 := adjust_file_name(fname)
	dir := path.Dir(fname2)
	ensure_dir(dir)
	name := path.Base(fname2)
	if len(name) > 100 {
		name = name[:100]
	}
	result := dir + "/" + name
	ioutil.WriteFile(result, content, fileMode)
	return result
}
