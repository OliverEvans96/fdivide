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

func getDirnames(inputDir string) []string {
	entryNames := ls(inputDir)
	var dirnames []string
	for _, entryName := range entryNames {
		entryPath := path.Join(inputDir, entryName)
		// Use Stat instead of Lstat to follow symlinks
		entryInfo, err := os.Stat(entryPath)
		if err != nil {
			panic(err)
		}
		if entryInfo.IsDir() {
			dirnames = append(dirnames, entryInfo.Name())
		}
	}
	return dirnames
}

func combine(inputDir string, outputDir string, verbose bool) {
	inputDirTruePath, err := filepath.EvalSymlinks(inputDir)
	if err != nil {
		panic(err)
	}
	inputDirAbsPath, err := filepath.Abs(inputDirTruePath)
	if err != nil {
		panic(err)
	}
	subdirnames := getDirnames(inputDirAbsPath)
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
			trueOldpath, err := filepath.EvalSymlinks(oldpath)
			if err != nil {
				panic(err)
			}
			newpath := path.Join(outputDir, filename)
			if verbose {
				fmt.Printf("%s -> %s\n", trueOldpath, newpath)
			}
			os.Symlink(trueOldpath, newpath)
		}
	}
}
