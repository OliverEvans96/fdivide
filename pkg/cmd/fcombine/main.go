package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/docopt/docopt-go"
)

const usage string = `fcombine
Combine files from sibiling subdirectories into a single output directory using symlinks.

Usage:
    fcombine [--link | --copy | --move] <input-parent-dir> <output-dir> [options]

Options:
    --link               Link files (default)
    --copy               Copy files
    --move               Move files
    --follow             Follow symlinks to source
    --verbose, -v        Verbose logging
`

type Method string

const (
	Link Method = "Link"
	Copy Method = "Copy"
	Move Method = "Move"
)

func main() {
	opts, err := docopt.ParseDoc(usage)
	if err != nil {
		panic(err)
	}

	copyFlag, err := opts.Bool("--copy")
	if err != nil {
		panic(err)
	}
	moveFlag, err := opts.Bool("--move")
	if err != nil {
		panic(err)
	}

	var method Method
	if copyFlag {
		method = Copy
	} else if moveFlag {
		method = Move
	} else {
		method = Link
	}

	inputDir, err := opts.String("<input-parent-dir>")
	if err != nil {
		panic(err)
	}
	outputDir, err := opts.String("<output-dir>")
	if err != nil {
		panic(err)
	}

	followFlag, err := opts.Bool("--follow")
	if err != nil {
		panic(err)
	}

	verbose, err := opts.Bool("--verbose")
	if err != nil {
		panic(err)
	}

	combine(inputDir, outputDir, method, followFlag, verbose)
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

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
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

func combine(inputDir string, outputDir string, method Method, followFlag bool, verbose bool) {
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
			if followFlag {
				oldpath, err = filepath.EvalSymlinks(oldpath)
				if err != nil {
					panic(err)
				}
			}
			newpath := path.Join(outputDir, filename)
			switch method {
			case Move:
				if verbose {
					fmt.Printf("%s ~> %s\n", oldpath, newpath)
				}
				os.Rename(oldpath, newpath)
			case Copy:
				if verbose {
					fmt.Printf("%s => %s\n", oldpath, newpath)
				}
				copy(oldpath, newpath)
			case Link:
				if verbose {
					fmt.Printf("%s -> %s\n", oldpath, newpath)
				}
				os.Symlink(oldpath, newpath)
			default:
				err = fmt.Errorf("Unknown method '%s'", method)
				panic(err)
			}
		}
	}
}
