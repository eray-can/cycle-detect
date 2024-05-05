package runner

import (
	"github.com/fsnotify/fsnotify"
	"testing"
	"time"
)

func mockEngine(t *testing.T) *Engine {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatalf("error creating watcher: %v", err)
	}
	return &Engine{
		logger:        mockLogger(),
		cycleDetector: mockDetector(),
		cnf:           mockConfig(),
		fsWatcher:     watcher,
	}
}

func TestEngine_Run(t *testing.T) {
	engine := mockEngine(t)

	go engine.Run()
	time.Sleep(5 * time.Second)
}

func TestEngineClose(t *testing.T) {
	engine := mockEngine(t)

	err := engine.Close()
	if err != nil {
		t.Error(err)
	}

}

func TestWatcher(t *testing.T) {
	engine := mockEngine(t)

	go engine.watcher()

	engine.fsWatcher.Events <- fsnotify.Event{Name: fileName(), Op: fsnotify.Write}

	time.Sleep(1 * time.Second)

	// Mock write event 2
	engine.fsWatcher.Events <- fsnotify.Event{Name: fileName(), Op: fsnotify.Write}

	time.Sleep(10 * time.Second)

}
