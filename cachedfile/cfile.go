// simplifies file access and adds a simple caching method
package cachedfile

import (
	"fmt"
	"simonwaldherr.de/go/golibs/cache"
	"simonwaldherr.de/go/golibs/file"
	"time"
)

var fileCache *cache.Cache
var cacheInit bool

func cacheWorker(filename string, value interface{}) {
	_, mtime, _, err := file.Time(filename)
	modify := mtime.UnixNano()
	if err == nil && modify < fileCache.Time(filename).UnixNano() {
		file.Write(filename, fmt.Sprint(value), false)
	}
}

func Read(filename string) (string, error) {
	if cacheInit == false {
		cacheInit = true
		fileCache = cache.New2(15*time.Minute, 1*time.Minute, cacheWorker)
	}
	var err error
	var data string
	filename, err = file.GetAbsolutePath(filename)
	if err != nil {
		return "", err
	}
	if xdata := fileCache.Get(filename); xdata == nil {
		if data, err = file.Read(filename); err != nil {
			return "", err
		}
		fileCache.Set(filename, data)
	} else {
		data = fmt.Sprint(xdata)
	}
	return data, nil
}

func Write(filename, str string, append bool) error {
	if cacheInit == false {
		cacheInit = true
		fileCache = cache.New2(15*time.Minute, 1*time.Minute, cacheWorker)
	}
	var err error
	var data string
	filename, err = file.GetAbsolutePath(filename)
	if err != nil {
		return err
	}
	if append {
		data, err = Read(filename)
		if err != nil {
			return err
		}
		fileCache.Set(filename, data+str)
	} else {
		fileCache.Set(filename, str)
	}
	return nil
}

func Size(filename string) (int64, error) {
	str, err := Read(filename)

	if err != nil {
		return -1, err
	}

	return int64(len(str)), nil
}

func Clean(filename string) error {
	return Write(filename, "", false)
}

func Stop() {
	if cacheInit == true {
		fileCache.DeleteAllWithFunc(cacheWorker)
	}
}

func Reset() {
	fileCache = cache.New2(15*time.Minute, 1*time.Minute, cacheWorker)
	cacheInit = false
}
