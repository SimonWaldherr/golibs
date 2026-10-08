package cachedfile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func newTestCache(t *testing.T, expiration, cleanup time.Duration) string {
	t.Helper()
	dir := t.TempDir()
	Init(expiration, cleanup)
	t.Cleanup(func() {
		if err := Flush(); err != nil {
			t.Errorf("cleanup Flush: %v", err)
		}
		Stop()
		// Remove only this test's entries if a failure prevented retry, so a
		// failed filesystem assertion cannot contaminate subsequent tests.
		cacheDir, err := absolutePath(dir)
		if err == nil {
			state.Lock()
			for filename := range state.entries {
				if strings.HasPrefix(filename, cacheDir+string(filepath.Separator)) {
					delete(state.entries, filename)
				}
			}
			state.Unlock()
		}
	})
	return dir
}

func expireFile(t *testing.T, filename string) {
	t.Helper()
	filename, err := absolutePath(filename)
	if err != nil {
		t.Fatal(err)
	}
	state.Lock()
	defer state.Unlock()
	state.entries[filename].expiration = time.Now().Add(-time.Second)
}

func readDisk(t *testing.T, filename string) string {
	t.Helper()
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}

func TestFlushCreatesFilesAndSupportsAppend(t *testing.T) {
	dir := newTestCache(t, time.Minute, 0)
	for _, append := range []bool{false, true} {
		filename := filepath.Join(dir, fmt.Sprint(append))
		if err := Write(filename, "first", append); err != nil {
			t.Fatal(err)
		}
		if err := Write(filename, "second", true); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(filename); !os.IsNotExist(err) {
			t.Fatalf("buffered file was already persisted: %v", err)
		}
		if err := Flush(); err != nil {
			t.Fatal(err)
		}
		if got := readDisk(t, filename); got != "firstsecond" {
			t.Fatalf("persisted content = %q", got)
		}
	}
}

func TestFlushRetainsFailedWritesForRetry(t *testing.T) {
	dir := newTestCache(t, time.Minute, 0)
	filename := filepath.Join(dir, "missing-parent", "file.txt")
	if err := Write(filename, "pending", false); err != nil {
		t.Fatal(err)
	}
	expireFile(t, filename)
	if err := Flush(); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Flush error = %v, want missing directory", err)
	}
	Stop()
	Reset()
	Init(time.Minute, 0)
	if got, err := Read(filename); err != nil || got != "pending" {
		t.Fatalf("pending read after lifecycle changes = %q, %v", got, err)
	}
	if err := os.Mkdir(filepath.Dir(filename), 0700); err != nil {
		t.Fatal(err)
	}
	if err := Flush(); err != nil {
		t.Fatal(err)
	}
	if got := readDisk(t, filename); got != "pending" {
		t.Fatalf("retried content = %q", got)
	}
}

func TestFlushAttemptsSuccessfulFilesAfterFailure(t *testing.T) {
	dir := newTestCache(t, time.Minute, 0)
	bad := filepath.Join(dir, "a-missing", "file")
	good := filepath.Join(dir, "z-good")
	for _, filename := range []string{bad, good} {
		if err := Write(filename, "content", false); err != nil {
			t.Fatal(err)
		}
	}
	if err := Flush(); err == nil {
		t.Fatal("Flush did not report failed write")
	}
	if got := readDisk(t, good); got != "content" {
		t.Fatalf("successful file = %q", got)
	}
	if err := os.Mkdir(filepath.Dir(bad), 0700); err != nil {
		t.Fatal(err)
	}
	if err := Flush(); err != nil {
		t.Fatal(err)
	}
}

func TestExpiredDirtyReadKeepsPendingContent(t *testing.T) {
	dir := newTestCache(t, time.Minute, 0)
	filename := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(filename, []byte("disk"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := Write(filename, "pending", false); err != nil {
		t.Fatal(err)
	}
	expireFile(t, filename)
	if got, err := Read(filename); err != nil || got != "pending" {
		t.Fatalf("expired dirty read = %q, %v", got, err)
	}
	if err := Write(filename, "+append", true); err != nil {
		t.Fatal(err)
	}
	if got := readDisk(t, filename); got != "disk" {
		t.Fatalf("disk changed before flush: %q", got)
	}
	if err := Flush(); err != nil {
		t.Fatal(err)
	}
	if got := readDisk(t, filename); got != "pending+append" {
		t.Fatalf("flushed content = %q", got)
	}
}

func TestReadUsesConfiguredExpiration(t *testing.T) {
	for _, expiration := range []time.Duration{0, -time.Minute, time.Minute} {
		t.Run(expiration.String(), func(t *testing.T) {
			dir := newTestCache(t, expiration, 0)
			filename := filepath.Join(dir, "file.txt")
			if err := os.WriteFile(filename, []byte("old"), 0600); err != nil {
				t.Fatal(err)
			}
			before := time.Now()
			if got, err := Read(filename); err != nil || got != "old" {
				t.Fatalf("initial read = %q, %v", got, err)
			}
			cachePath, err := absolutePath(filename)
			if err != nil {
				t.Fatal(err)
			}
			deadline := func() time.Time {
				state.Lock()
				defer state.Unlock()
				return state.entries[cachePath].expiration
			}()
			if expiration <= 0 {
				if !deadline.IsZero() {
					t.Fatalf("nonpositive TTL has deadline %v", deadline)
				}
			} else if deadline.Before(before.Add(expiration)) || deadline.After(time.Now().Add(expiration)) {
				t.Fatalf("deadline %v does not use configured TTL %v", deadline, expiration)
			}
			if err := os.WriteFile(filename, []byte("new"), 0600); err != nil {
				t.Fatal(err)
			}
			if got, err := Read(filename); err != nil || got != "old" {
				t.Fatalf("cached read = %q, %v", got, err)
			}
			if expiration > 0 {
				expireFile(t, filename)
				if got, err := Read(filename); err != nil || got != "new" {
					t.Fatalf("expired clean read = %q, %v", got, err)
				}
			}
		})
	}
}

func TestInitResetAndStopPersistPendingWrites(t *testing.T) {
	dir := newTestCache(t, time.Minute, 0)
	for name, action := range map[string]func(){
		"Init":  func() { Init(time.Hour, 0) },
		"Reset": Reset,
		"Stop":  Stop,
	} {
		t.Run(name, func(t *testing.T) {
			filename := filepath.Join(dir, name)
			if err := Write(filename, "pending", false); err != nil {
				t.Fatal(err)
			}
			action()
			if got := readDisk(t, filename); got != "pending" {
				t.Fatalf("content after %s = %q", name, got)
			}
		})
	}
	Stop()
	Stop()
	state.Lock()
	initialized, stop := state.initialized, state.stop
	state.Unlock()
	if initialized || stop != nil {
		t.Fatal("Stop left background cleanup active")
	}
}

func TestConcurrentAppendsAndLifecycle(t *testing.T) {
	dir := newTestCache(t, time.Minute, 0)
	filename := filepath.Join(dir, "appends.txt")
	const writers = 64
	var wg sync.WaitGroup
	errors := make(chan error, writers)
	for i := 0; i < writers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := Write(filename, fmt.Sprintf("%d\n", i), true); err != nil {
				errors <- err
			}
		}(i)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 20; i++ {
			Init(time.Minute, 0)
			Reset()
			Stop()
		}
	}()
	wg.Wait()
	close(errors)
	for err := range errors {
		t.Error(err)
	}
	if err := Flush(); err != nil {
		t.Fatal(err)
	}
	lines := strings.Fields(readDisk(t, filename))
	if len(lines) != writers {
		t.Fatalf("persisted %d appends, want %d", len(lines), writers)
	}
	seen := make(map[string]bool)
	for _, line := range lines {
		if seen[line] {
			t.Fatalf("duplicate append %q", line)
		}
		seen[line] = true
	}
	for i := 0; i < writers; i++ {
		if !seen[fmt.Sprint(i)] {
			t.Errorf("missing append %d", i)
		}
	}
}

func TestCleanerFlushesExpiredWritesAndStops(t *testing.T) {
	dir := newTestCache(t, time.Hour, time.Millisecond)
	filename := filepath.Join(dir, "file.txt")
	if err := Write(filename, "automatic", false); err != nil {
		t.Fatal(err)
	}
	expireFile(t, filename)
	deadline := time.NewTimer(2 * time.Second)
	defer deadline.Stop()
	poll := time.NewTicker(time.Millisecond)
	defer poll.Stop()
	for {
		content, err := os.ReadFile(filename)
		if err == nil && string(content) == "automatic" {
			break
		}
		select {
		case <-deadline.C:
			t.Fatalf("cleaner did not persist pending write: %q, %v", content, err)
		case <-poll.C:
		}
	}
	state.Lock()
	stopped := state.stopped
	state.Unlock()
	Stop()
	select {
	case <-stopped:
	default:
		t.Fatal("Stop returned before cleaner exited")
	}
}

func TestFlushReportsExternalModification(t *testing.T) {
	dir := newTestCache(t, time.Minute, 0)
	filename := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(filename, []byte("initial"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := Write(filename, "buffered", false); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filename, []byte("external"), 0600); err != nil {
		t.Fatal(err)
	}
	future := time.Now().Add(time.Hour)
	if err := os.Chtimes(filename, future, future); err != nil {
		t.Fatal(err)
	}
	if err := Flush(); !errors.Is(err, ErrConflict) {
		t.Fatalf("Flush error = %v, want ErrConflict", err)
	}
	if got := readDisk(t, filename); got != "external" {
		t.Fatalf("Flush overwrote external content: %q", got)
	}
	if got, err := Read(filename); err != nil || got != "buffered" {
		t.Fatalf("buffered conflict read = %q, %v", got, err)
	}
	// An explicit replacement establishes a fresh disk baseline and resolves
	// the conflict without relying on the file's timestamp being in the past.
	if err := Write(filename, "replacement", false); err != nil {
		t.Fatal(err)
	}
	if err := Flush(); err != nil {
		t.Fatal(err)
	}
	if got := readDisk(t, filename); got != "replacement" {
		t.Fatalf("explicit replacement = %q", got)
	}
}

func TestPathAliasesSharePendingWrites(t *testing.T) {
	dir := newTestCache(t, time.Minute, 0)
	filename := filepath.Join(dir, "file.txt")
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	relative, err := filepath.Rel(cwd, filename)
	if err != nil {
		t.Fatal(err)
	}
	if err := Write(relative, "pending", false); err != nil {
		t.Fatal(err)
	}
	if got, err := Read(filename); err != nil || got != "pending" {
		t.Fatalf("absolute alias read = %q, %v", got, err)
	}
}

func TestPendingAppendRetainsOriginalDiskSnapshot(t *testing.T) {
	dir := newTestCache(t, time.Minute, 0)
	filename := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(filename, []byte("initial"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := Write(filename, "+first", true); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filename, []byte("external modification"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := Write(filename, "+second", true); err != nil {
		t.Fatal(err)
	}
	if err := Flush(); !errors.Is(err, ErrConflict) {
		t.Fatalf("Flush after second append = %v, want ErrConflict", err)
	}
	if got := readDisk(t, filename); got != "external modification" {
		t.Fatalf("external content overwritten: %q", got)
	}
	if got, err := Read(filename); err != nil || got != "initial+first+second" {
		t.Fatalf("pending content = %q, %v", got, err)
	}
	if err := Write(filename, "resolved", false); err != nil {
		t.Fatal(err)
	}
	if err := Flush(); err != nil {
		t.Fatal(err)
	}
}

func TestCleanAppendRejectsChangedDiskSnapshot(t *testing.T) {
	dir := newTestCache(t, time.Minute, 0)
	filename := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(filename, []byte("initial"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := Read(filename); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filename, []byte("external modification"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := Write(filename, "+append", true); !errors.Is(err, ErrConflict) {
		t.Fatalf("append after clean read = %v, want ErrConflict", err)
	}
	if err := Flush(); err != nil {
		t.Fatal(err)
	}
	if got := readDisk(t, filename); got != "external modification" {
		t.Fatalf("external content overwritten: %q", got)
	}
	if err := Write(filename, "+append", true); err != nil {
		t.Fatal(err)
	}
	if err := Flush(); err != nil {
		t.Fatal(err)
	}
	if got := readDisk(t, filename); got != "external modification+append" {
		t.Fatalf("append after refreshed read = %q", got)
	}
}

func TestExistingFutureTimestampDoesNotConflict(t *testing.T) {
	dir := newTestCache(t, time.Minute, 0)
	for _, append := range []bool{false, true} {
		filename := filepath.Join(dir, fmt.Sprint(append))
		if err := os.WriteFile(filename, []byte("initial"), 0600); err != nil {
			t.Fatal(err)
		}
		future := time.Now().Add(time.Hour)
		if err := os.Chtimes(filename, future, future); err != nil {
			t.Fatal(err)
		}
		if err := Write(filename, "pending", append); err != nil {
			t.Fatal(err)
		}
		if err := Flush(); err != nil {
			t.Fatalf("unchanged future-mtime file conflicts: %v", err)
		}
		want := "pending"
		if append {
			want = "initialpending"
		}
		if got := readDisk(t, filename); got != want {
			t.Fatalf("future-mtime file = %q, want %q", got, want)
		}
	}
}

func TestCreatedOrReplacedFilesConflict(t *testing.T) {
	dir := newTestCache(t, time.Minute, 0)
	for _, existed := range []bool{false, true} {
		filename := filepath.Join(dir, fmt.Sprint(existed))
		var original os.FileInfo
		if existed {
			if err := os.WriteFile(filename, []byte("initial"), 0600); err != nil {
				t.Fatal(err)
			}
			var err error
			original, err = os.Stat(filename)
			if err != nil {
				t.Fatal(err)
			}
		}
		if err := Write(filename, "buffered", false); err != nil {
			t.Fatal(err)
		}
		replacement := filename + ".replacement"
		if err := os.WriteFile(replacement, []byte("newfile"), 0600); err != nil {
			t.Fatal(err)
		}
		if existed {
			// Same length and mtime still conflict when file identity changes.
			if err := os.Chtimes(replacement, original.ModTime(), original.ModTime()); err != nil {
				t.Fatal(err)
			}
		}
		if err := os.Rename(replacement, filename); err != nil {
			t.Fatal(err)
		}
		if err := Flush(); !errors.Is(err, ErrConflict) {
			t.Fatalf("created/replaced file Flush = %v, want ErrConflict", err)
		}
		if got := readDisk(t, filename); got != "newfile" {
			t.Fatalf("created/replaced file overwritten: %q", got)
		}
		if err := Write(filename, "resolved", false); err != nil {
			t.Fatal(err)
		}
		if err := Flush(); err != nil {
			t.Fatal(err)
		}
	}
}
