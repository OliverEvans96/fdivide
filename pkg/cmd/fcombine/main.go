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
    --include-files      Also transfer non-directory files in <input-parent-dir> to <output-dir>
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

type FcombineOptions struct {
	InputDir     string
	OutputDir    string
	Method       Method
	Follow       bool
	IncludeFiles bool
	HiddenDirs   bool
	DryRun       bool
	Verbose      bool
}

func main() {
	options := getOptions()
	combine(options)
}

func getOptions() FcombineOptions {
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

	includeFilesFlag, err := opts.Bool("--include-files")
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

	followFlag, err := opts.Bool("--follow")
	if err != nil {
		panic(err)
	}
	hiddenDirsFlag, err := opts.Bool("--hidden-dirs")
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

	return FcombineOptions{
		InputDir:     inputDir,
		OutputDir:    outputDir,
		Method:       method,
		Follow:       followFlag,
		IncludeFiles: includeFilesFlag,
		HiddenDirs:   hiddenDirsFlag,
		DryRun:       dryRunFlag,
		Verbose:      verbose,
	}
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

func combine(options FcombineOptions) {
	if options.DryRun {
		fmt.Println("DRY RUN")
	}
	err := os.MkdirAll(options.OutputDir, 0755)
	if err != nil {
		panic(err)
	}
	inputDirTruePath, err := filepath.EvalSymlinks(options.InputDir)
	if err != nil {
		panic(err)
	}
	inputDirAbsPath, err := filepath.Abs(inputDirTruePath)
	if err != nil {
		panic(err)
	}
	outputDirTruePath, err := filepath.EvalSymlinks(options.OutputDir)
	if err != nil {
		panic(err)
	}
	outputDirAbsPath, err := filepath.Abs(outputDirTruePath)
	if err != nil {
		panic(err)
	}
	subdirnames := getDirnames(inputDirAbsPath, options.HiddenDirs)
	if err != nil {
		panic(err)
	}

	var allLinkSpecs []LinkSpec
	// Get dirs to combine
	for _, subdirname := range subdirnames {
		subdirPath := path.Join(inputDirAbsPath, subdirname)
		// Ignore options.OutputDir (only relevent if it's a subdirectory of options.InputDir)
		if subdirPath == outputDirAbsPath {
			continue
		}
		filenames := getAllFilenames(subdirPath)
		for _, filename := range filenames {
			oldpath := path.Join(subdirPath, filename)
			if options.Follow {
				oldpath, err = filepath.EvalSymlinks(oldpath)
				if err != nil {
					panic(err)
				}
			}
			newpath := path.Join(options.OutputDir, filename)
			linkSpec := LinkSpec{oldpath, newpath}
			allLinkSpecs = append(allLinkSpecs, linkSpec)
		}
	}
	// Get files if desired
	if options.IncludeFiles {
		filenames := getAllFilenames(inputDirAbsPath)
		for _, filename := range filenames {
			oldpath := path.Join(inputDirAbsPath, filename)
			if options.Follow {
				oldpath, err = filepath.EvalSymlinks(oldpath)
				if err != nil {
					panic(err)
				}
			}
			newpath := path.Join(options.OutputDir, filename)
			linkSpec := LinkSpec{oldpath, newpath}
			allLinkSpecs = append(allLinkSpecs, linkSpec)
		}
	}
	// Perform file transfer
	for _, linkSpec := range allLinkSpecs {
		switch options.Method {
		case Move:
			if options.Verbose {
				fmt.Printf("%s ~> %s\n", linkSpec.Oldpath, linkSpec.Newpath)
			}
			if !options.DryRun {
				os.Rename(linkSpec.Oldpath, linkSpec.Newpath)
			}
		case Copy:
			if options.Verbose {
				fmt.Printf("%s => %s\n", linkSpec.Oldpath, linkSpec.Newpath)
			}
			if !options.DryRun {
				copyFile(linkSpec.Oldpath, linkSpec.Newpath)
			}
		case Link:
			if options.Verbose {
				fmt.Printf("%s -> %s\n", linkSpec.Oldpath, linkSpec.Newpath)
			}
			if !options.DryRun {
				os.Symlink(linkSpec.Oldpath, linkSpec.Newpath)
			}
		default:
			err = fmt.Errorf("Unknown method '%s'", options.Method)
			panic(err)
		}
	}
}
