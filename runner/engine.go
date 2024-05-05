package runner

import (
	"github.com/fsnotify/fsnotify"
	"time"
)

type Engine struct {
	logger        ILogger
	cnf           *Config
	cycleDetector *Detector
	fsWatcher     *fsnotify.Watcher
}

type IEngine interface {
	Run()
	Close() error
}

func NewEngine() IEngine {
	cnf, err := ReadFileConfig()
	if err != nil {
		panic(err)
	}
	lgr := NewLogger(cnf)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		lgr.ErrorLog(err.Error())
	}
	return &Engine{
		logger:        lgr,
		cnf:           cnf,
		cycleDetector: NewCycleDetector(*cnf.Build.BaseProjectPath, lgr),
		fsWatcher:     watcher,
	}
}

func (e *Engine) Close() error {
	err := e.fsWatcher.Close()
	if err != nil {
		e.logger.ErrorLog(err.Error())
		return err
	}

	return nil
}

func (e *Engine) Run() {
	//first run
	err := e.cycleDetector.Run(false)
	if err != nil {
		e.logger.ErrorLog(err.Error())
	}
	e.watcher()
}

func (e *Engine) watcher() {

	projectPath := *e.cnf.Build.BaseProjectPath
	if pathErr := e.fsWatcher.Add(projectPath); pathErr != nil {
		e.logger.ErrorLog(pathErr.Error())
		return
	}

	defer func(logger ILogger) {
		if r := recover(); r != nil {
			logger.ErrorLogf("[WATCHER][PASSED]", r)
		}
	}(e.logger)

	fileModifications := make(map[string]time.Time)

	// Start listening for events.
	for {
		select {
		case event, ok := <-e.fsWatcher.Events:
			if !ok {
				return
			}
			// Handle only write events
			if eventData, check := fileModifications[event.Name]; check {
				if time.Since(eventData) < time.Second {
					continue
				}
			} else {
				fileModifications[event.Name] = time.Now()
			}

			// Watcher detect run
			if watcherErr := e.cycleDetector.Run(true); watcherErr != nil {
				e.logger.ErrorLog(watcherErr.Error())
			}

		case wErr, ok := <-e.fsWatcher.Errors:
			if !ok {
				e.logger.ErrorLog(wErr.Error())
				return
			}
		}
	}

}
