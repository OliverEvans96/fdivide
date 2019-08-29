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
Symlinks are created by default, but files can also be moved or copied using the flags.

Usage:
    fcombine <input-parent-dir> <output-dir> [options]

Options:
    --copy               Copy files
    --move               Move files
    --follow             Follow symlinks to source
    --hidden-dirs        Include dirs beginning with "." - excluded by default
    --dry-run, -n        Don't actually perform any operations, only print statements
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
		if moveFlag {
			err := fmt.Errorf("at most one of [--copy,--move] may be given")
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
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
	hiddenFlag, err := opts.Bool("--hidden-dirs")
	if err != nil {
		panic(err)
	}
	dryRunFlag, err := opts.Bool("--dry-run")
	if err != nil {
		panic(err)
	}

	verbose, err := opts.Bool("--verbose")
	if err != nil {
		panic(err)
	}

	combine(inputDir, outputDir, method, followFlag, hiddenFlag, dryRunFlag, verbose)
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

func copyFile(src, dst string) (int64, error) {
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

func startswith(s string, r rune) bool {
	return rune(s[0]) == r
}

func getDirnames(inputDir string, hiddenFlag bool) []string {
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
			name := entryInfo.Name()
			if !startswith(name, '.') || hiddenFlag {
				dirnames = append(dirnames, name)
			}
		}
	}
	return dirnames
}

type LinkSpec struct {
	Oldpath string
	Newpath string
}

func combine(inputDir string, outputDir string, method Method, followFlag bool, hiddenFlag bool, dryRunFlag bool, verbose bool) {
	if dryRunFlag {
		fmt.Println("DRY RUN")
	}
	inputDirTruePath, err := filepath.EvalSymlinks(inputDir)
	if err != nil {
		panic(err)
	}
	inputDirAbsPath, err := filepath.Abs(inputDirTruePath)
	if err != nil {
		panic(err)
	}
	outputDirTruePath, err := filepath.EvalSymlinks(outputDir)
	if err != nil {
		panic(err)
	}
	outputDirAbsPath, err := filepath.Abs(outputDirTruePath)
	if err != nil {
		panic(err)
	}
	subdirnames := getDirnames(inputDirAbsPath, hiddenFlag)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		panic(err)
	}

	var allLinkSpecs []LinkSpec
	for _, subdirname := range subdirnames {
		subdirPath := path.Join(inputDirAbsPath, subdirname)
		// Ignore outputDir (only relevent if it's a subdirectory of inputDir)
		if subdirPath == outputDirAbsPath {
			continue
		}
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
			linkSpec := LinkSpec{oldpath, newpath}
			allLinkSpecs = append(allLinkSpecs, linkSpec)
		}
	}
	for _, linkSpec := range allLinkSpecs {
		switch method {
		case Move:
			if verbose {
				fmt.Printf("%s ~> %s\n", linkSpec.Oldpath, linkSpec.Newpath)
			}
			if !dryRunFlag {
				os.Rename(linkSpec.Oldpath, linkSpec.Newpath)
			}
		case Copy:
			if verbose {
				fmt.Printf("%s => %s\n", linkSpec.Oldpath, linkSpec.Newpath)
			}
			if !dryRunFlag {
				copyFile(linkSpec.Oldpath, linkSpec.Newpath)
			}
		case Link:
			if verbose {
				fmt.Printf("%s -> %s\n", linkSpec.Oldpath, linkSpec.Newpath)
			}
			if !dryRunFlag {
				os.Symlink(linkSpec.Oldpath, linkSpec.Newpath)
			}
		default:
			err = fmt.Errorf("Unknown method '%s'", method)
			panic(err)
		}
	}
}
