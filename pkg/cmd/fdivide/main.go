package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"path/filepath"

	"github.com/docopt/docopt-go"
)

const usage string = `fdivide
Divide a regular files from a directory into subdirectories by number of files using symlinks.

Usage:
    fdivide --size <dir-size> <input-dir> <output-dir>
    fdivide --into <num-dirs> <input-dir> <output-dir>
`

func main() {
	opts, err := docopt.ParseDoc(usage)
	if err != nil {
		panic(err)
	}
	bySize, err := opts.Bool("--size")
	if err != nil {
		panic(err)
	}
	inputDir, err := opts.String("<input-dir>")
	if err != nil {
		panic(err)
	}
	outputDir, err := opts.String("<output-dir>")
	if err != nil {
		panic(err)
	}
	if bySize {
		dirSize, err := opts.Int("<dir-size>")
		if err != nil {
			panic(err)
		}
		divideBySize(dirSize, inputDir, outputDir)

	} else {
		numDirs, err := opts.Int("<num-dirs>")
		if err != nil {
			panic(err)
		}
		divideInto(numDirs, inputDir, outputDir)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getRegularFilenames(inputDir string) []string {
	allFiles, err := ioutil.ReadDir(inputDir)
	if err != nil {
		panic(err)
	}
	var regularFilenames []string
	for _, file := range allFiles {
		if file.Mode().IsRegular() {
			regularFilenames = append(regularFilenames, file.Name())
		}
	}
	return regularFilenames
}

func symlink(oldpath string, newpath string) {
	// Create symlink, deleting newpath if it already exists
	err := os.Symlink(oldpath, newpath)
	if err != nil {
		if _, err := os.Lstat(newpath); err == nil {
			os.Remove(newpath)
			err = os.Symlink(oldpath, newpath)
			if err != nil {
				panic(err)
			}
		}
	}

}

func divideBySize(dirSize int, inputDir string, outputDir string) {
	filenames := getRegularFilenames(inputDir)
	numFiles := len(filenames)
	numDirs := int(math.Ceil(float64(numFiles) / float64(dirSize)))
	doDivide(numDirs, dirSize, numFiles, filenames, inputDir, outputDir)
}

func divideInto(numDirs int, inputDir string, outputDir string) {
	filenames := getRegularFilenames(inputDir)
	numFiles := len(filenames)
	dirSize := int(math.Ceil(float64(numFiles) / float64(numDirs)))
	doDivide(numDirs, dirSize, numFiles, filenames, inputDir, outputDir)
}

func getDirNameTemplate(numDirs int) string {
	dirNameDigits := int(math.Ceil(math.Log10(float64(numDirs))))
	// Create the template string, e.g. %10d
	return fmt.Sprintf("%%0%dd", dirNameDigits)
}

func doDivide(numDirs int, dirSize int, numFiles int, filenames []string, inputDir string, outputDir string) {
	dirNameTemplate := getDirNameTemplate(numDirs)
	inputDirAbsPath, err := filepath.Abs(inputDir)
	if err != nil {
		panic(err)
	}
	for dirNum := 0; dirNum < numDirs; dirNum++ {
		subdirName := fmt.Sprintf(dirNameTemplate, dirNum)
		subdirPath := path.Join(outputDir, subdirName)
		err := os.MkdirAll(subdirPath, 0755)
		if err != nil {
			panic(err)
		}
		maxFileNumPlusOne := int(min(numFiles, dirSize*(dirNum+1)))
		for fileNum := dirSize * dirNum; fileNum < maxFileNumPlusOne; fileNum++ {
			filename := filenames[fileNum]
			oldpath := path.Join(inputDirAbsPath, filename)
			newpath := path.Join(subdirPath, filename)
			os.Symlink(oldpath, newpath)
		}
	}
}
