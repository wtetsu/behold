/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package gazer

import (
	"os/exec"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/wtetsu/gaze/pkg/config"
	"github.com/wtetsu/gaze/pkg/fs"
	"github.com/wtetsu/gaze/pkg/logger"
	"github.com/wtetsu/gaze/pkg/notify"
	"github.com/wtetsu/gaze/pkg/time"
)

// Gazer gazes filesystem.
type Gazer struct {
	patterns []string
	notify   *notify.Notify
	isClosed bool
	counter  uint64
}

// New returns a new Gazer.
func New(patterns []string) *Gazer {
	cleanPatterns := make([]string, len(patterns))
	for i, p := range patterns {
		cleanPatterns[i] = filepath.Clean(p)
	}

	notify, _ := notify.New(cleanPatterns)
	return &Gazer{
		patterns: cleanPatterns,
		notify:   notify,
		isClosed: false,
	}
}

// Close disposes internal resources.
func (g *Gazer) Close() {
	if g.isClosed {
		return
	}
	g.notify.Close()
	g.isClosed = true
}

// Run starts to gaze.
func (g *Gazer) Run(configs *config.Config, timeout int, restart bool) error {
	err := g.repeatRunAndWait(configs, timeout, restart)
	return err
}

func (g *Gazer) repeatRunAndWait(commandConfigs *config.Config, timeout int, restart bool) error {
	var lastExecutionTime int64

	sigInt := sigIntChannel()

	var ignorePeriod int64 = 10 * 1000000

	var ongoingCommand *exec.Cmd

	isDisposed := false
	for {
		if isDisposed {
			break
		}
		select {
		case event, ok := <-g.notify.Events:
			logger.Debug("Receive: %s", event.Name)
			flag := fsnotify.Write | fsnotify.Rename
			if ok && event.Op|flag == 0 {
				continue
			}
			if !matchAny(g.patterns, event.Name) {
				continue
			}
			modifiedTime := time.GetFileModifiedTime(event.Name)
			if (modifiedTime - lastExecutionTime) < ignorePeriod {
				continue
			}

			g.counter++
			commandString := getAppropriateCommand(event.Name, commandConfigs)
			if commandString == "" {
				logger.Debug("Command not found: %s", event.Name)
				continue
			}

			logger.NoticeWithBlank("[%s]", commandString)

			if ongoingCommand != nil {
				kill(ongoingCommand, "Restart")
				ongoingCommand = nil
			}

			cmd := createCommand(commandString)
			lastExecutionTime = time.Now()
			if !restart {
				err := executeCommandOrTimeout(cmd, timeout)
				if err != nil {
					logger.NoticeObject(err)
				}
			} else {
				// restartable
				ongoingCommand = cmd
				go func() {
					err := executeCommandOrTimeout(cmd, timeout)
					if err != nil {
						logger.NoticeObject(err)
					}
					ongoingCommand = nil
				}()
			}

		case <-sigInt:
			isDisposed = true
			return nil
		}
	}
	return nil
}

func matchAny(watchFiles []string, s string) bool {
	result := false
	for _, f := range watchFiles {
		if fs.GlobMatch(f, s) {
			result = true
			break
		}
	}
	return result
}

func getAppropriateCommand(filePath string, commandConfigs *config.Config) string {
	var result string
	for _, c := range commandConfigs.Commands {
		if c.Run == "" || c.Ext == "" && c.Re == "" {
			continue
		}
		if c.Match(filePath) {
			command := render(c.Run, filePath)
			result = command
			break
		}
	}
	return result
}

// Counter returns the current execution counter
func (g *Gazer) Counter() uint64 {
	return g.counter
}
