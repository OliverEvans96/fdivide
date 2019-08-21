package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/docopt/docopt-go"
)

const usage string = `fcombine
Combine files from sibiling subdirectories into a single output directory using symlinks.

Usage:
    fcombine <input-parent-dir> <output-dir> [--verbose]
`

func main() {
	opts, err := docopt.ParseDoc(usage)
	if err != nil {
		panic(err)
	}
	verbose, err := opts.Bool("--verbose")
	if err != nil {
		panic(err)
	}
	inputDir, err := opts.String("<input-parent-dir>")
	if err != nil {
		panic(err)
	}
	outputDir, err := opts.String("<output-dir>")
	if err != nil {
		panic(err)
	}

	combine(inputDir, outputDir, verbose)
}

func getAllFilenames(inputDir string) []string {
	allFiles, err := ioutil.ReadDir(inputDir)
	if err != nil {
		panic(err)
	}
	var allFilenames []string
	for _, file := range allFiles {
		if !file.IsDir() {
			allFilenames = append(allFilenames, file.Name())
		}
	}
	return allFilenames
}

func getDirnames(inputDir string) []string {
	allFiles, err := ioutil.ReadDir(inputDir)
	if err != nil {
		panic(err)
	}
	var dirnames []string
	for _, file := range allFiles {
		if file.IsDir() {
			dirnames = append(dirnames, file.Name())
		}
	}
	return dirnames
}

func combine(inputDir string, outputDir string, verbose bool) {
	subdirnames := getDirnames(inputDir)
	inputDirAbsPath, err := filepath.Abs(inputDir)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		panic(err)
	}

	for _, subdirname := range subdirnames {
		subdirPath := path.Join(inputDirAbsPath, subdirname)
		filenames := getAllFilenames(subdirPath)
		for _, filename := range filenames {
			oldpath := path.Join(subdirPath, filename)
			newpath := path.Join(outputDir, filename)
			if verbose {
				fmt.Printf("%s -> %s\n", oldpath, newpath)
			}
			os.Symlink(oldpath, newpath)
		}
	}
}
