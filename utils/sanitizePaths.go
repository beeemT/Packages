package utils

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

//PathToAbsDir takes a string and returns the absolute Dir of that string.
//Calls PathToAbs on path.
func PathToAbsDir(path string) (string, error) {
	if !IsDir(path) {
		path = filepath.Dir(path)
	}
	return PathToAbs(path)
}

//PathToAbsFile takes a string and returns the absolute file of that string.
//Returns an error if path does not point to a file.
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

//IsDirSubdirOf checks if subdir is a subdirectory of dir. Does also check if subdir exists.
//Accepts relative paths. Calls PathToAbs beforehand. Returns false if one of the dirs is not a dir.
func IsDirSubdirOf(subdir, dir string) (bool, error) {
	if !(IsDir(subdir) && IsDir(dir)) {
		return false, nil
	}
	subdir, err := PathToAbs(subdir)
	dir, err = PathToAbs(dir)
	if err != nil {
		return false, err
	}
	if strings.HasPrefix(subdir, dir) {
		return true, nil
	}
	return false, nil
}

//Exists checks if a element at path exists.
//Calls PathToAbs on path beforehand
func Exists(path string) (bool, error) {
	path, err := PathToAbs(path)
	_, err = os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
