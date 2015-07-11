package cachedfile

import (
	"../cache"
	"../file"
	"fmt"
	"time"
)

var fileCache *cache.Cache
var cacheInit bool = false

func cacheWorker(key string, value interface{}) {
	file.Write(key, fmt.Sprint(value), false)
}

func CachedRead(fn string) (string, error) {
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

func CachedWrite(fn, str string, append bool) error {
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
		data, err = CachedRead(fn)
		if err != nil {
			return err
		}
		fileCache.Set(fn, data+str)
	} else {
		fileCache.Set(fn, str)
	}
	return nil
}

func StopCache() {
	if cacheInit == true {
		fileCache.DeleteAllWithFunc(cacheWorker)
	}
}
