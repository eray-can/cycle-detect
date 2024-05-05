package runner

import (
	"testing"
)

func mockLogger() ILogger {
	return NewLogger(mockConfig())
}

func TestLogger(t *testing.T) {
	message := "Test log message"

	logger := mockLogger()

	logger.InfoLog(message)
	logger.Printf("print f message %s", "hello")
	logger.ErrorLogf("error message %s", "hello")
	logger.InfoLogf("info messagef %s", "hello")
	logger.Print(message)
	logger.Printf("%s", message)
	t.Log("PASSED")
}
