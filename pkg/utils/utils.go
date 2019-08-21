package utils

import (
	"io/ioutil"
	"os"
)

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func GetRegularFilenames(inputDir string) []string {
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

func GetAllFilenames(inputDir string) []string {
	allFiles, err := ioutil.ReadDir(inputDir)
	if err != nil {
		panic(err)
	}
	var allFilenames []string
	for _, file := range allFiles {
		allFilenames = append(allFilenames, file.Name())
	}
	return allFilenames
}

func GetDirnames(inputDir string) []string {
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

func Symlink(oldpath string, newpath string) {
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
