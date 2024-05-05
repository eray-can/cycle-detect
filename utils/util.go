package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	moduleNotFound     = "Module not found"
	verticalSep        = "|"
	fSlashSep          = "/"
	pointSep           = "."
	doubleBackSlashSep = "\\"
	doubleQueSep       = "\""
	//
	filePath = "go.mod"
	module   = "module "
	empty    = ""
	test     = "_test"
	goFile   = "go"
	//
	onlyMessage = "Only controls golang projects.Please make sure that [go.mod] is found in the file path you provided"

	//format
	steepF   = "%s|%s"
	defaultF = "%s%s"
)

// Get only golang files
// Does not receive test files
func IsDetectFile(fileName string) bool {
	fileNames := strings.Split(strings.ToLower(fileName), pointSep)
	if len(fileNames) > 1 {
		return !strings.HasSuffix(fileNames[len(fileNames)-2], test) && fileNames[len(fileNames)-1] == goFile
	}

	return false
}

// get go mod module name
func GetModuleName(dir string) (string, error) {
	//check go.mod file
	goMod := filepath.Join(dir, filePath)
	moduleName, err := findModuleName(goMod)
	if err != nil {
		return empty, fmt.Errorf(onlyMessage)
	}

	return moduleName, nil
}

func findModuleName(modFilePath string) (string, error) {
	file, err := os.Open(modFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, module) {
			moduleName := strings.TrimSpace(strings.TrimPrefix(line, module))
			return moduleName, nil
		}
	}

	if scanErr := scanner.Err(); scanErr != nil {
		return empty, scanErr
	}

	return "", fmt.Errorf(moduleNotFound)
}

func DetectProjectImport(packageName string, moduleName string) (string, bool) {
	//because packages come with double quotes
	packageName = strings.Replace(packageName, doubleQueSep, empty, -1)

	//check original project import, not package import todo burayı configden yonet
	if strings.HasPrefix(packageName, moduleName) {
		return packageName[len(moduleName):], true
	}

	return "", false
}

// key get filepath
func GetFilePath(text string) string {
	return strings.Split(text, verticalSep)[0]
}

// key get packageName
func GetPackageName(text string) string {
	return strings.Split(text, verticalSep)[1]
}

func GetBaseUri(FullPath string) string {
	dir, _ := filepath.Split(GetFilePath(FullPath))
	return dir
}

// I keep imports as file name and package
func ImportKeyGenerate(importName string, detectPackageName string) string {
	return strings.ToLower(fmt.Sprintf(steepF, strings.Replace(importName, doubleBackSlashSep, fSlashSep, -1), detectPackageName))
}

// build full path
func BuildFullPath(moduleName string, importName string) string {
	return fmt.Sprintf(defaultF, moduleName, importName)
}

func FileSplitDir(dir string, path string) string {
	text := strings.Replace(GetBaseUri(path), dir, empty, -1)

	return text[:len(text)-1]
}

// I shortened the path because long paths look bad in the table
func MinimizePath(dir string, path string) string {
	return strings.Replace(GetFilePath(path), dir, empty, -1)
}
