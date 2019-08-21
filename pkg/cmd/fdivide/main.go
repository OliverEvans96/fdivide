package main

import (
	"fmt"
	"math"
	"os"
	"path"
	"path/filepath"

	"github.com/docopt/docopt-go"
	"gitlab.com/lavo-nutrition/fdivide/pkg/utils"
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

func divideBySize(dirSize int, inputDir string, outputDir string) {
	filenames := utils.GetRegularFilenames(inputDir)
	numFiles := len(filenames)
	numDirs := int(math.Ceil(float64(numFiles) / float64(dirSize)))
	divide(numDirs, dirSize, numFiles, filenames, inputDir, outputDir)
}

func divideInto(numDirs int, inputDir string, outputDir string) {
	filenames := utils.GetRegularFilenames(inputDir)
	numFiles := len(filenames)
	dirSize := int(math.Ceil(float64(numFiles) / float64(numDirs)))
	divide(numDirs, dirSize, numFiles, filenames, inputDir, outputDir)
}

func getDirNameTemplate(numDirs int) string {
	dirnameDigits := int(math.Ceil(math.Log10(float64(numDirs))))
	// Create the template string, e.g. %10d
	return fmt.Sprintf("%%0%dd", dirnameDigits)
}

func divide(numDirs int, dirSize int, numFiles int, filenames []string, inputDir string, outputDir string) {
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
		maxFileNumPlusOne := int(utils.Min(numFiles, dirSize*(dirNum+1)))
		for fileNum := dirSize * dirNum; fileNum < maxFileNumPlusOne; fileNum++ {
			filename := filenames[fileNum]
			oldpath := path.Join(inputDirAbsPath, filename)
			newpath := path.Join(subdirPath, filename)
			utils.Symlink(oldpath, newpath)
		}
	}
}
