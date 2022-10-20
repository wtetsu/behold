/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package notify

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/wtetsu/gaze/pkg/fs"
	"github.com/wtetsu/gaze/pkg/logger"
	"github.com/wtetsu/gaze/pkg/time"
	"github.com/wtetsu/gaze/pkg/uniq"
)

// Notify delives events to a channel when files are virtually updated.
// "create+rename" is regarded as "update".
type Notify struct {
	Events                  chan Event
	Errors                  chan error
	watcher                 *fsnotify.Watcher
	isClosed                bool
	times                   map[string]int64
	pendingPeriod           int64
	regardRenameAsModPeriod int64
	detectCreate            bool
}

// Event represents a single file system notification.
type Event struct {
	Name string
	Time int64
}

// Op describes a set of file operations.
type Op = fsnotify.Op

// Close disposes internal resources.
func (n *Notify) Close() {
	if n.isClosed {
		return
	}
	n.watcher.Close()
	n.isClosed = true
}

// New creates a Notify
func New(patterns []string, maxWatchDirs int) (*Notify, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.ErrorObject(err)
		return nil, err
	}

	watchDirs := findDirs(patterns, maxWatchDirs)

	if len(watchDirs) > maxWatchDirs {
		logger.Error(strings.Join(watchDirs[:maxWatchDirs], "\n") + "\n...")
		return nil, errors.New("too many watchDirs")
	}

	for _, t := range watchDirs {
		err = watcher.Add(t)
		if err != nil {
			logger.Error("%s: %v", t, err)
		} else {
			logger.Info("gazing at: %s", t)
		}
	}

	notify := &Notify{
		Events:                  make(chan Event),
		watcher:                 watcher,
		isClosed:                false,
		times:                   make(map[string]int64),
		pendingPeriod:           100,
		regardRenameAsModPeriod: 1000,
		detectCreate:            true,
	}

	go notify.wait()

	return notify, nil
}

func findDirs(patterns []string, maxWatchDirs int) []string {
	targets := uniq.New()

	for _, pattern := range patterns {
		patternDir := filepath.Dir(pattern)

		realDir := findRealDirectory(patternDir)
		if len(realDir) > 0 {
			targets.Add(realDir)
		}
		if targets.Len() > maxWatchDirs {
			return targets.List()
		}

		_, dirs1 := fs.Find(pattern)
		for _, d := range dirs1 {
			targets.Add(d)
		}
		if targets.Len() > maxWatchDirs {
			return targets.List()
		}

		_, dirs2 := fs.Find(patternDir)
		for _, d := range dirs2 {
			targets.Add(d)
		}
		if targets.Len() > maxWatchDirs {
			return targets.List()
		}
	}
	return targets.List()
}

func findRealDirectory(path string) string {
	entries := strings.Split(filepath.ToSlash(filepath.Clean(path)), "/")

	currentPath := ""
	for i := 0; i < len(entries); i++ {
		globIndex := strings.IndexAny(entries[i], "*?[{\\")
		if globIndex != -1 {
			break
		}

		currentPath += entries[i] + string(filepath.Separator)
	}
	currentPath = fs.TrimSuffix(currentPath, string(filepath.Separator))

	if fs.IsDir(currentPath) {
		return currentPath
	} else {
		return ""
	}
}

func (n *Notify) wait() {
	for {
		select {
		case event, ok := <-n.watcher.Events:

			normalizedName := filepath.Clean(event.Name)

			if event.Op == fsnotify.Create && fs.IsDir(normalizedName) {
				logger.Info("gazing at: %s", normalizedName)
				n.watcher.Add(normalizedName)
			}

			if !ok {
				continue
			}
			if !n.shouldExecute(normalizedName, event.Op) {
				continue
			}
			logger.Debug("notified: %s: %s", normalizedName, event.Op)
			now := time.Now()
			n.times[normalizedName] = now
			e := Event{
				Name: normalizedName,
				Time: now,
			}
			n.Events <- e
		case err, ok := <-n.watcher.Errors:
			if !ok {
				continue
			}
			n.Errors <- err
		}
	}
}

func (n *Notify) shouldExecute(filePath string, op Op) bool {
	const W = fsnotify.Write
	const R = fsnotify.Rename
	const C = fsnotify.Create

	if op != W && op != R && !(n.detectCreate && op == C) {
		logger.Debug("skipped: %s: %s (Op is not applicable)", filePath, op)
		return false
	}

	lastExecutionTime := n.times[filePath]

	if !fs.IsFile(filePath) {
		logger.Debug("skipped: %s: %s (not a file)", filePath, op)
		return false
	}

	modifiedTime := time.GetFileModifiedTime(filePath)

	if op == W || op == C {
		elapsed := modifiedTime - lastExecutionTime
		logger.Debug("lastExecutionTime(%s): %d, %d", op, lastExecutionTime, elapsed)
		if elapsed < n.pendingPeriod*1000000 {
			logger.Debug("skipped: %s: %s (too frequent)", filePath, op)
			return false
		}
	}
	if op == R {
		elapsed := time.Now() - modifiedTime
		logger.Debug("lastExecutionTime(%s): %d, %d", op, lastExecutionTime, elapsed)
		if elapsed > n.regardRenameAsModPeriod*1000000 {
			logger.Debug("skipped: %s: %s (unnatural rename)", filePath, op)
			return false
		}
	}

	return true
}

// PendingPeriod sets new pendingPeriod(ms).
func (n *Notify) PendingPeriod(p int64) {
	n.pendingPeriod = p
}

// Requeue requeue an event.
func (n *Notify) Requeue(event Event) {
	n.Events <- event
}
