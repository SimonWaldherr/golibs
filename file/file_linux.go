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
	atime := time.Unix(int64(stat.Atim.Sec), int64(stat.Atim.Nsec))
	ctime := time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec))
	return atime, mtime, ctime, nil
}
