package runner

import (
	"os"
	"testing"
	"time"
)

// TODO Don't forget to read the comments left for the dir and filename functions
// give the path to any GOLANG project in your local
func dir() string {
	return "C:/Users/ceray/OneDrive/Masaüstü/sites/kubernetes"
}

// give a golang file path to any GOLANG project in your local area
func fileName() string {
	return "C:/Users/ceray/OneDrive/Masaüstü/sites/kubernetes/staging/src/k8s.io/api/core/v1/lifecycle.go"
}

func mockDetector() *Detector {
	return NewCycleDetector(dir(), mockLogger())
}

func TestDetector_Run(t *testing.T) {

	detector := mockDetector()
	if err := detector.Run(false); err != nil {
		t.Errorf("Run() returned an error: %v", err)
	}

	// Test with watcher
	detector.clear()
	if err := detector.Run(true); err != nil {
		t.Errorf("Run() with watcher returned an error: %v", err)
	}
}

func TestDetector_processFile(t *testing.T) {
	detector := mockDetector()
	file, fileErr := os.Open(fileName())
	if fileErr != nil {
		t.Errorf("file problem %s", fileErr.Error())
	}
	fileStat, statErr := file.Stat()
	if statErr != nil {
		t.Errorf("file problem %s", statErr.Error())
	}

	err := detector.processFile(fileName(), fileStat, nil)
	if err != nil {
		t.Errorf("processFile() returned an error for Go file: %v", err)
	}

}

func TestDetector_newDetect(t *testing.T) {
	detector := mockDetector()
	if detector == nil {
		t.Error("NewCycleDetector() returned nil for a valid directory")
	}
}

func TestDetector_extractImports(t *testing.T) {
	detector := mockDetector()

	err := detector.extractImports(fileName())
	if err != nil {
		t.Errorf("extractImports() returned an error for a valid Go file: %v", err)
	}

	err = detector.extractImports("/path/to/file.txt")
	if err == nil {
		t.Errorf("extractImports() did not return an error for a non-Go file")
	}
}

func TestDetector_buildTable(t *testing.T) {
	detector := &Detector{
		detectPackage: [][]string{
			{"file1", "file2"},
			{"import1", "import2"},
		},
		stats: Stats{
			responseTime:     time.Now(),
			scannedGoFile:    10,
			scannedImports:   20,
			scannedTotalFile: 30,
		},
		Logger: mockLogger(),
	}

	detector.buildTable()

}

func TestDetector_clear(t *testing.T) {
	detector := mockDetector()

	// Populate some data for testing
	detector.stats.responseTime = time.Now()
	detector.stats.scannedGoFile = 10
	detector.stats.scannedImports = 20
	detector.stats.scannedTotalFile = 30
	detector.detectPackage = [][]string{{"file1", "file2"}}
	detector.clear()

	// Check if the data is cleared
	if !detector.stats.responseTime.IsZero() ||
		detector.stats.scannedGoFile != 0 ||
		detector.stats.scannedImports != 0 ||
		detector.stats.scannedTotalFile != 0 ||
		len(detector.detectPackage) != 0 {
		t.Error("clear() did not clear the data properly")
	}
}
