package main

import (
	"os"
	"path"
	"path/filepath"

	"github.com/docopt/docopt-go"
	"gitlab.com/lavo-nutrition/fdivide/pkg/utils"
)

const usage string = `fcombine
Combine files from sibiling subdirectories into a single output directory using symlinks.

Usage:
    fcombine <input-parent-dir> <output-dir>
`

func main() {
	opts, err := docopt.ParseDoc(usage)
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

	combine(inputDir, outputDir)
}

func combine(inputDir string, outputDir string) {
	subdirnames := utils.GetDirnames(inputDir)
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
		filenames := utils.GetAllFilenames(subdirPath)
		for _, filename := range filenames {
			oldpath := path.Join(subdirPath, filename)
			newpath := path.Join(outputDir, filename)
			utils.Symlink(oldpath, newpath)
		}
	}
}
