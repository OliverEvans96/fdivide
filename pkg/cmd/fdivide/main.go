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
Divide regular files from a directory into subdirectories by number of files using symlinks.

Usage:
    fdivide --size <dir-size> <input-dir> <output-dir> [--verbose]
    fdivide --into <num-dirs> <input-dir> <output-dir> [--verbose]
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
	verbose, err := opts.Bool("--verbose")
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
		divideBySize(dirSize, inputDir, outputDir, verbose)

	} else {
		numDirs, err := opts.Int("<num-dirs>")
		if err != nil {
			panic(err)
		}
		divideInto(numDirs, inputDir, outputDir, verbose)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func ls(dirPath string) []string {
	var entryNames []string
	entries, err := ioutil.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		entryNames = append(entryNames, entry.Name())
	}
	return entryNames
}

func getAllFilenames(inputDir string) []string {
	allEntrynames := ls(inputDir)
	var allFilenames []string
	for _, entryName := range allEntrynames {
		entryPath := path.Join(inputDir, entryName)
		// Use Stat instead of Lstat to follow symlinks
		entryInfo, err := os.Stat(entryPath)
		if err != nil {
			panic(err)
		}
		if !entryInfo.IsDir() {
			allFilenames = append(allFilenames, entryInfo.Name())
		}
	}
	return allFilenames
}

func divideBySize(dirSize int, inputDir string, outputDir string, verbose bool) {
	filenames := getAllFilenames(inputDir)
	numFiles := len(filenames)
	numDirs := int(math.Ceil(float64(numFiles) / float64(dirSize)))
	divide(numDirs, dirSize, numFiles, filenames, inputDir, outputDir, verbose)
}

func divideInto(numDirs int, inputDir string, outputDir string, verbose bool) {
	filenames := getAllFilenames(inputDir)
	numFiles := len(filenames)
	dirSize := int(math.Ceil(float64(numFiles) / float64(numDirs)))
	divide(numDirs, dirSize, numFiles, filenames, inputDir, outputDir, verbose)
}

func getDirNameTemplate(numDirs int) string {
	dirnameDigits := int(math.Ceil(math.Log10(float64(numDirs))))
	// Create the template string, e.g. %10d
	return fmt.Sprintf("%%0%dd", dirnameDigits)
}

func divide(numDirs int, dirSize int, numFiles int, filenames []string, inputDir string, outputDir string, verbose bool) {
	dirnameTemplate := getDirNameTemplate(numDirs)
	inputDirAbsPath, err := filepath.Abs(inputDir)
	if err != nil {
		panic(err)
	}
	for dirNum := 0; dirNum < numDirs; dirNum++ {
		subdirname := fmt.Sprintf(dirnameTemplate, dirNum)
		subdirPath := path.Join(outputDir, subdirname)
		err := os.MkdirAll(subdirPath, 0755)
		if err != nil {
			panic(err)
		}
		maxFileNumPlusOne := int(min(numFiles, dirSize*(dirNum+1)))
		for fileNum := dirSize * dirNum; fileNum < maxFileNumPlusOne; fileNum++ {
			filename := filenames[fileNum]
			oldpath := path.Join(inputDirAbsPath, filename)
			trueOldpath, err := filepath.EvalSymlinks(oldpath)
			if err != nil {
				panic(err)
			}
			newpath := path.Join(subdirPath, filename)
			if verbose {
				fmt.Printf("%s -> %s\n", trueOldpath, newpath)
			}
			os.Symlink(trueOldpath, newpath)
		}
	}
}
