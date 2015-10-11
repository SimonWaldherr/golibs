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
	stat := file.Sys().(*syscall.Stat_t)
	atime := time.Unix(int64(stat.Atimespec.Sec), int64(stat.Atimespec.Nsec))
	ctime := time.Unix(int64(stat.Ctimespec.Sec), int64(stat.Ctimespec.Nsec))
	return atime, mtime, ctime, nil
}
