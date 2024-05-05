package runner

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"text/template"
	"time"
)

type Logger struct {
	template *template.Template
	file     *os.File
	stdout   *os.File
}

type logTemplate struct {
	Level       string
	Message     string
	StackTracee string
	Time        string
}

var (
	logPath    = ".app.log"
	folderName = "logs"

	// Error types
	lvlInfo  = "Info"
	lvlError = "Error"
)

type ILogger interface {
	InfoLog(msg string)
	InfoLogf(msg string, data ...interface{})
	ErrorLog(msg string)
	ErrorLogf(msg string, data ...interface{})
	Print(msg string)
	Printf(msg string, data ...interface{})
}

func NewLogger(cfg *Config) ILogger {

	logger := &Logger{}

	logger.selectFile()

	logger.template = template.Must(
		template.New("logTemplate").Parse(
			" Time: [{{.Time}}], Level: [{{.Level}}], Message: [{{.Message}}], StackTracee: [{{.StackTracee}}]\n"),
	)

	return logger
}

func (l *Logger) InfoLogf(msg string, data ...interface{}) {
	l.log(fmt.Sprintf(msg, data...), lvlInfo)
}

func (l *Logger) InfoLog(msg string) {
	l.log(msg, lvlInfo)
}

func (l *Logger) ErrorLog(msg string) {
	l.log(msg, lvlError)
}

func (l *Logger) Print(msg string) {
	fmt.Println(msg)
}

func (l *Logger) Printf(msg string, data ...interface{}) {
	fmt.Println(fmt.Sprintf(msg, data...))
}

func (l *Logger) ErrorLogf(msg string, data ...interface{}) {
	l.log(fmt.Sprintf(msg, data...), lvlError)
}

func (l *Logger) log(msg string, logLevel string) {
	logEntry := buildTemplate(msg, logLevel)

	err := l.template.Execute(l.stdout, logEntry)
	if err != nil {
		panic(err)
	}

	err = l.template.Execute(l.file, logEntry)
	if err != nil {
		panic(err)
	}
}

func buildTemplate(msg string, logLevel string) logTemplate {
	return logTemplate{
		Message:     fmt.Sprint(msg),
		Level:       logLevel,
		StackTracee: strings.Replace(string(debug.Stack()), "\n", " ", -1),
		Time:        time.Now().Format(time.RFC3339),
	}
}

// if writing to a file, select the file
func (l *Logger) selectFile() {
	var err error

	cnfPath, err := os.Getwd()
	if err != nil {
		return
	}
	cnfPath = fmt.Sprintf("%s/%s/%s", cnfPath, folderName, logPath)

	l.file, err = openFile(cnfPath)
	if err != nil {
		cErr := createFolder()
		if !cErr {
			return
		}
		l.file, err = os.Create(cnfPath)
		if err != nil {
			return
		}
	}

	l.stdout = os.Stdout
}

func createFolder() bool {
	err := os.Mkdir(folderName, os.ModeDir|os.ModePerm)
	if err != nil {
		return false
	}

	return true
}

func openFile(logPath string) (*os.File, error) {
	file, fErr := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if fErr != nil {
		return nil, fErr
	}

	return file, nil
}
