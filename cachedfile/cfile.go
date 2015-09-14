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

func cacheWorker(key string, value interface{}) {
	file.Write(key, fmt.Sprint(value), false)
}

func Read(fn string) (string, error) {
	if cacheInit == false {
		cacheInit = true
		fileCache = cache.New2(15*time.Minute, 1*time.Minute, cacheWorker)
	}
	var err error
	var data string
	fn, err = file.GetAbsolutePath(fn)
	if err != nil {
		return "", err
	}
	if xdata := fileCache.Get(fn); xdata == nil {
		if data, err = file.Read(fn); err != nil {
			return "", err
		}
		fileCache.Set(fn, data)
	} else {
		data = fmt.Sprint(xdata)
	}
	return data, nil
}

func Write(fn, str string, append bool) error {
	if cacheInit == false {
		cacheInit = true
		fileCache = cache.New2(15*time.Minute, 1*time.Minute, cacheWorker)
	}
	var err error
	var data string
	fn, err = file.GetAbsolutePath(fn)
	if err != nil {
		return err
	}
	if append {
		data, err = Read(fn)
		if err != nil {
			return err
		}
		fileCache.Set(fn, data+str)
	} else {
		fileCache.Set(fn, str)
	}
	return nil
}

func Size(fn string) (int64, error) {
	str, err := Read(fn)

	if err != nil {
		return -1, err
	}

	return int64(len(str)), nil
}

func Clean(fn string) error {
	return Write(fn, "", false)
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
