package file

import (
	"os"
	"syscall"
	"time"
)

func Time(fn string) (time.Time, time.Time, time.Time, error) {
	file, err := os.Stat(fn)
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, err
	}
	mtime := file.ModTime()
	atime := time.Unix(0, file.Sys().(*syscall.Win32FileAttributeData).LastAccessTime.Nanoseconds())
	ctime := time.Unix(0, file.Sys().(*syscall.Win32FileAttributeData).LastWriteTime.Nanoseconds())
	return atime, mtime, ctime, nil
}
