package runner

import (
	"fmt"
	"github.com/eray-can/cycle-detect/utils"
	"github.com/olekukonko/tablewriter"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	//table val
	path     = "Path"
	imports  = "Imports"
	tableVal = " %s => %s "

	//response val
	response = "Response Time: %s, Total file: %d, Go file: %d, Imports: %d | [%d-%02d-%02d %02d:%02d:%02d]"

	//err msg val
	importErrMsg = "Error extracting imports from %s: %v"
)

type Detector struct {
	wg            *sync.WaitGroup
	mutex         *sync.Mutex
	allPackage    map[string][]string
	detectPackage [][]string
	dir           string
	moduleName    string
	stats         Stats
	Logger        ILogger
}

type Stats struct {
	responseTime     time.Time
	scannedGoFile    int
	scannedImports   int
	scannedTotalFile int
}

func (d *Detector) Run(isWatcher bool) error {
	if isWatcher {
		d.clear()
	}

	d.stats.responseTime = time.Now()
	//Browses all folders and subfiles
	err := filepath.Walk(d.dir, d.processFile)
	if err != nil {
		return err
	}

	d.wg.Wait()
	d.newDetect()

	return nil
}

func (d *Detector) processFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	//increase the number of scanned files
	d.stats.scannedTotalFile++

	//just check the go files
	if utils.IsDetectFile(info.Name()) {

		//increase the number of scanned go files by one
		d.stats.scannedGoFile++
		d.wg.Add(1)

		go func(filePath string) {
			defer d.wg.Done()

			//import Extract
			errImp := d.extractImports(filePath)
			if errImp != nil {
				d.Logger.ErrorLogf(importErrMsg, filePath, errImp)
			}
		}(path)

	}

	return nil
}

func (d *Detector) newDetect() {
	//compares i with all j until i runs out
	for i, iImports := range d.allPackage {
		for j, jImports := range d.allPackage {

			//no self-import cycle
			if i == j {
				continue
			}

			d.wg.Add(1)

			go func(iImports, jImports []string, i string, j string) {
				defer d.wg.Done()

				//i is the format in which other files should import i
				iPath := utils.FileSplitDir(d.dir, i)

				//j is the format in which other files should import j
				jPath := utils.FileSplitDir(d.dir, j)

				//Compare i imports with all j imports
				for iImportIdx := range iImports {
					for jImportIdx := range jImports {

						//Check if the imported path and file import match
						if iImports[iImportIdx] == jPath && iPath == jImports[jImportIdx] {
							if d.checkDetect(utils.MinimizePath(d.dir, i), utils.MinimizePath(d.dir, j)) {
								//Write all data to show in a table
								d.mutex.Lock()
								d.detectPackage = append(d.detectPackage, []string{
									fmt.Sprintf(
										tableVal, utils.MinimizePath(d.dir, i),
										utils.MinimizePath(d.dir, j),
									),
									fmt.Sprintf(tableVal,
										utils.BuildFullPath(d.moduleName, iImports[iImportIdx]),
										utils.BuildFullPath(d.moduleName, jImports[jImportIdx]),
									),
								})
								d.mutex.Unlock()
							}
						}
					}
				}
			}(iImports, jImports, i, j)
		}
	}

	//created table
	d.buildTable()
}

// It was done to avoid adding the same data in different rows
// exam: /test/blabla.go => /runner/engine.go | /runner/engine.go => /test/blabla.go
// no need to show the same thing in reverse
func (d *Detector) checkDetect(i string, j string) bool {
	for idx := range d.detectPackage {
		parts := strings.Split(d.detectPackage[idx][0], " => ")

		if i == strings.TrimSpace(parts[1]) && j == strings.TrimSpace(parts[0]) {
			return false
		}

	}

	return true
}

// imports only the project's own imports
// does not take external packages within the project
func (d *Detector) extractImports(filename string) error {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	//check the imports in the file
	for i := range f.Imports {
		//increase the number of scanned imports by one
		d.stats.scannedImports++

		//Extract imports in go file
		detectImportName, found := utils.DetectProjectImport(f.Imports[i].Path.Value, d.moduleName)
		if found {
			//ilgili importları kaydet
			d.mutex.Lock()
			key := utils.ImportKeyGenerate(filename, f.Name.Name)
			d.allPackage[key] = append(d.allPackage[key], detectImportName)
			d.mutex.Unlock()
		}

	}

	return nil
}

// Set the table and the required data to be displayed on the terminal
func (d *Detector) buildTable() {
	t := time.Now()
	d.Logger.Printf(response, time.Since(d.stats.responseTime),
		d.stats.scannedTotalFile, d.stats.scannedGoFile, d.stats.scannedImports,
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{path, imports})

	for _, result := range d.detectPackage {
		table.Append(result)
	}

	table.Render()
}

// We delete old data because Watcher repeats
func (d *Detector) clear() {
	d.stats.responseTime = time.Time{}
	d.stats.scannedGoFile = 0
	d.stats.scannedImports = 0
	d.stats.scannedTotalFile = 0
	d.detectPackage = nil
}

// NewCycleDetector New Cycle
func NewCycleDetector(dir string, logger ILogger) *Detector {

	// searches for the go mod file in the given path and gets the modul name
	// Terminates the application if not found
	moduleName, err := utils.GetModuleName(dir)
	if err != nil {
		logger.ErrorLog(err.Error())
		os.Exit(1)
	}

	return &Detector{
		allPackage: make(map[string][]string),
		moduleName: moduleName,
		dir:        strings.ToLower(dir),
		wg:         &sync.WaitGroup{},
		mutex:      &sync.Mutex{},
		stats:      Stats{},
		Logger:     logger,
	}
}
