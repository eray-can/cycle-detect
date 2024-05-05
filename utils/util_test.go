package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	Key      = "c:/users/bleagle/filepath/path.go|blabla"
	FilePath = "c:/users/bleagle/filepath/path.go"
	Dir      = "C:/Users/bleagle/filepath/"
)

// give the path to any GOLANG project in your local
func dir() string {
	return "C:/Users/ceray/OneDrive/Masaüstü/sites/kubernetes"
}

func TestDetectProjectImport(t *testing.T) {

	detectImportName, found := DetectProjectImport("test/blabla", "test/")
	assert.Equal(t, "blabla", detectImportName)
	assert.Equal(t, true, found)
}

func TestGetFilePath(t *testing.T) {
	path := GetFilePath(Key)
	assert.Equal(t, "c:/users/bleagle/filepath/path.go", path)
}

func TestGetPackageName(t *testing.T) {
	packageName := GetPackageName(Key)
	assert.Equal(t, "blabla", packageName)
}

func TestBuildFullPath(t *testing.T) {
	fullPath := BuildFullPath("test/", "bleagle")
	assert.Equal(t, "test/bleagle", fullPath)
}

func TestKeyGenerate(t *testing.T) {
	fileKey := ImportKeyGenerate(FilePath, "blabla")
	assert.Equal(t, " c:/users/bleagle/filepath/path.go|blabla", fileKey)
}

func TestGetModuleName(t *testing.T) {
	moduleName, err := GetModuleName(dir())
	if err != nil {
		t.Error("err", err.Error())
	}
	// give the module name of the directory you provided
	assert.Equal(t, "k8s.io/kubernetes", moduleName)
}

func TestIsGoFile(t *testing.T) {
	isGo := IsDetectFile(FilePath)
	assert.Equal(t, true, isGo)
}

func TestFileSplitDir(t *testing.T) {

	response := FileSplitDir(Dir, Key)
	assert.Equal(t, "c:/users/bleagle/filepath", response)

}
