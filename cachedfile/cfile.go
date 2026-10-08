// Package cachedfile buffers file reads and writes in memory. Buffered writes
// are persisted by Flush, Stop, or automatic cleanup after their TTL expires.
package cachedfile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"simonwaldherr.de/go/golibs/file"
)

const (
	defaultExpiration = 15 * time.Minute
	defaultCleanup    = time.Minute
)

// ErrConflict reports disk state that changed since a cached read or pending
// write. Pending content remains available and is not discarded.
var ErrConflict = errors.New("cachedfile: file changed on disk")

type diskSnapshot struct {
	exists   bool
	modified time.Time
	size     int64
	info     os.FileInfo
}

func snapshotDisk(filename string) (diskSnapshot, error) {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return diskSnapshot{}, nil
	}
	if err != nil {
		return diskSnapshot{}, err
	}
	return diskSnapshot{exists: true, modified: info.ModTime(), size: info.Size(), info: info}, nil
}

func (snapshot diskSnapshot) matches(other diskSnapshot) bool {
	if snapshot.exists != other.exists {
		return false
	}
	return !snapshot.exists || (snapshot.modified.Equal(other.modified) &&
		snapshot.size == other.size && os.SameFile(snapshot.info, other.info))
}

type entry struct {
	content    string
	expiration time.Time
	disk       diskSnapshot
	dirty      bool
}

type fileState struct {
	sync.Mutex
	entries     map[string]*entry
	expiration  time.Duration
	cleanup     time.Duration
	initialized bool
	stop        chan struct{}
	stopped     chan struct{}
}

var state = fileState{
	entries:    make(map[string]*entry),
	expiration: defaultExpiration,
	cleanup:    defaultCleanup,
}

func absolutePath(filename string) (string, error) {
	if filename == "" {
		return "", os.ErrInvalid
	}
	if filename == "~" || strings.HasPrefix(filename, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		filename = filepath.Join(home, strings.TrimPrefix(filename, "~"))
	}
	filename, err := filepath.Abs(filename)
	if err != nil {
		return "", err
	}
	// The file helper resolves existing symlinks while permitting a new file.
	return file.GetAbsolutePath(filename)
}

func (s *fileState) deadline() time.Time {
	if s.expiration <= 0 {
		return time.Time{}
	}
	return time.Now().Add(s.expiration)
}

func (e *entry) expired(now time.Time) bool {
	return !e.expiration.IsZero() && !e.expiration.After(now)
}

func (s *fileState) startLocked() {
	if s.initialized {
		return
	}
	s.initialized = true
	if s.cleanup <= 0 {
		return
	}
	s.stop = make(chan struct{})
	s.stopped = make(chan struct{})
	go s.cleaner(s.cleanup, s.stop, s.stopped)
}

func (s *fileState) stopLocked() <-chan struct{} {
	stopped := s.stopped
	if s.stop != nil {
		close(s.stop)
	}
	s.stop = nil
	s.stopped = nil
	s.initialized = false
	return stopped
}

func (s *fileState) cleaner(interval time.Duration, stop <-chan struct{}, stopped chan<- struct{}) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer close(stopped)
	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			s.Lock()
			select {
			case <-stop:
				s.Unlock()
				return
			default:
			}
			// Failed writes stay buffered for a later cleanup or explicit Flush.
			_ = s.flushLocked(true)
			s.Unlock()
		}
	}
}

func flushEntry(filename string, e *entry) error {
	current, err := snapshotDisk(filename)
	if err != nil {
		return err
	}
	if !e.disk.matches(current) {
		return ErrConflict
	}
	return file.Write(filename, e.content, false)
}

func readSnapshot(filename string) (string, diskSnapshot, error) {
	before, err := snapshotDisk(filename)
	if err != nil {
		return "", diskSnapshot{}, err
	}
	content, err := file.Read(filename)
	if err != nil {
		return "", diskSnapshot{}, err
	}
	after, err := snapshotDisk(filename)
	if err != nil {
		return "", diskSnapshot{}, err
	}
	if !before.matches(after) {
		return "", diskSnapshot{}, fmt.Errorf("cachedfile: read %q: %w", filename, ErrConflict)
	}
	return content, after, nil
}

// flushLocked attempts every selected entry and returns the first error in path
// order. Successful entries are removed; failed writes remain available to Read.
func (s *fileState) flushLocked(expiredOnly bool) error {
	now := time.Now()
	paths := make([]string, 0, len(s.entries))
	for filename, e := range s.entries {
		if !expiredOnly || e.expired(now) {
			paths = append(paths, filename)
		}
	}
	sort.Strings(paths)
	var firstErr error
	for _, filename := range paths {
		e := s.entries[filename]
		if e.dirty {
			if err := flushEntry(filename, e); err != nil {
				if firstErr == nil {
					firstErr = fmt.Errorf("cachedfile: flush %q: %w", filename, err)
				}
				continue
			}
		}
		delete(s.entries, filename)
	}
	return firstErr
}

// Init flushes the current cache, then configures its TTL and cleanup interval.
// Nonpositive TTLs disable expiration; nonpositive intervals disable background
// cleanup. Pending writes that cannot be flushed are retained. Call Flush to
// observe or retry those errors.
func Init(expirationTime, cleanupInterval time.Duration) {
	state.Lock()
	stopped := state.stopLocked()
	_ = state.flushLocked(false)
	state.expiration = expirationTime
	state.cleanup = cleanupInterval
	state.startLocked()
	state.Unlock()
	if stopped != nil {
		<-stopped
	}
}

func (s *fileState) readLocked(filename string) (string, error) {
	if e, ok := s.entries[filename]; ok {
		// Dirty data is authoritative until successfully persisted, even after
		// expiration. Reading must never replace it with stale disk content.
		if e.dirty || !e.expired(time.Now()) {
			return e.content, nil
		}
		delete(s.entries, filename)
	}
	content, disk, err := readSnapshot(filename)
	if err != nil {
		return "", err
	}
	s.entries[filename] = &entry{content: content, expiration: s.deadline(), disk: disk}
	return content, nil
}

// Read returns buffered content, or reads and caches the file from disk.
func Read(filename string) (string, error) {
	filename, err := absolutePath(filename)
	if err != nil {
		return "", err
	}
	state.Lock()
	defer state.Unlock()
	state.startLocked()
	return state.readLocked(filename)
}

// Write buffers content. Append reads and updates the buffer atomically, and
// creates a new file when the target does not exist. Appending to a clean cache
// entry rejects disk changes since its read; pending appends retain the first
// write's disk snapshot so Flush can report conflicts. A replacement (append is
// false) explicitly establishes a fresh snapshot against the current disk state.
// Flush reports disk errors and compares file identity, modification time, and
// size with that snapshot before writing.
func Write(filename, str string, append bool) error {
	filename, err := absolutePath(filename)
	if err != nil {
		return err
	}
	state.Lock()
	defer state.Unlock()
	state.startLocked()
	var disk diskSnapshot
	if append {
		content, err := state.readLocked(filename)
		if err != nil {
			if !os.IsNotExist(err) {
				return err
			}
			disk, err = snapshotDisk(filename)
			if err != nil {
				return err
			}
			// A file created after the failed read must not be overwritten by
			// an append built from empty content.
			if disk.exists {
				return fmt.Errorf("cachedfile: append %q: %w", filename, ErrConflict)
			}
		} else {
			e := state.entries[filename]
			disk = e.disk
			if !e.dirty {
				current, err := snapshotDisk(filename)
				if err != nil {
					return err
				}
				if !disk.matches(current) {
					return fmt.Errorf("cachedfile: append %q: %w", filename, ErrConflict)
				}
			}
		}
		str = content + str
	} else {
		disk, err = snapshotDisk(filename)
		if err != nil {
			return err
		}
	}
	state.entries[filename] = &entry{
		content: str, expiration: state.deadline(), disk: disk, dirty: true,
	}
	return nil
}

// Size returns the byte length of the buffered or on-disk file content.
func Size(filename string) (int64, error) {
	str, err := Read(filename)
	if err != nil {
		return -1, err
	}
	return int64(len(str)), nil
}

// Clean buffers an empty file.
func Clean(filename string) error {
	return Write(filename, "", false)
}

// Flush persists all pending writes and invalidates successfully flushed cache
// entries. It attempts every file and returns the first error in path order.
// Failed writes remain buffered, including after Stop, Reset, or Init.
func Flush() error {
	state.Lock()
	defer state.Unlock()
	return state.flushLocked(false)
}

// Stop flushes pending writes and stops background cleanup. It is idempotent.
// Failed writes remain buffered; use Flush to observe and retry errors. A later
// Read or Write restarts cleanup with the most recently configured settings.
func Stop() {
	state.Lock()
	stopped := state.stopLocked()
	_ = state.flushLocked(false)
	state.Unlock()
	if stopped != nil {
		<-stopped
	}
}

// Reset flushes pending writes and restores the default TTL and cleanup interval.
// Failed writes are retained and can be retried with Flush.
func Reset() {
	Init(defaultExpiration, defaultCleanup)
}
