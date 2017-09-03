package sanitize

import (
	"errors"
	"os"
	"path/filepath"
)

//PathToAbsDir takes a string and returns the absolute Dir of that string.
func PathToAbsDir(path string) (string, error) {
	if !IsDir(path) {
		path = filepath.Dir(path)
	}
	return PathToAbs(path)
}

//PathToAbsFile takes a string and returns the absolute file of that string.
//returns an error if path does not point to a file
func PathToAbsFile(path string) (string, error) {
	if IsFile(path) {
		return PathToAbs(path)
	}
	return "", errors.New("path is not a file: " + path)
}

//PathToAbs takes a string and returns the absolute Path of that string
func PathToAbs(path string) (string, error) {
	return filepath.Abs(path)
}

//IsDir takes a string and returns whether path is a Directory
func IsDir(path string) bool {
	path, _ = PathToAbs(path)
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

//IsFile takes a string and returns whether path is a File
//calls PathToAbs on path beforehand
func IsFile(path string) bool {
	path, _ = PathToAbs(path)
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}
