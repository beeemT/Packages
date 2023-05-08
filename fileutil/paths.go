package fileutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//FileTypeError is used whenever a path does not point to the expected file type.
type FileTypeError struct {
	Name string
	v    interface{}
}

func (err FileTypeError) Error() string {
	return fmt.Sprintf("%s: %s", err.Name, err.v)
}

//PathToAbsDir takes a string and returns the absolute Dir of that string.
//Calls PathToAbs on path.
func PathToAbsDir(path string) (string, error) {
	path, err := PathToAbs(path)
	if err != nil {
		return "", err
	}

	if IsDir(path) {
		return path, nil
	}
	return "", &FileTypeError{"FileTypeError: Path does not point to directory", path}
}

//PathToAbsFile returns the absolute filepath of that string.
//Returns an error if path does not point to a file.
func PathToAbsFile(path string) (string, error) {
	path, err := PathToAbs(path)
	if err != nil {
		return "", err
	}

	if IsFile(path) {
		return path, nil
	}
	return "", &FileTypeError{"FileTypeError: Path does not point to file", path}
}

//PathToAbs returns the absolute Path of that string.
func PathToAbs(path string) (string, error) {
	return filepath.Abs(path)
}

//IsDir returns whether path is a Directory.
func IsDir(path string) bool {
	path, _ = PathToAbs(path)
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

//IsFile takes a string and returns whether path is a File.
//Calls PathToAbs on path beforehand.
//If the file doesn't exist false is returned.
func IsFile(path string) bool {
	path, _ = PathToAbs(path)
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

//IsDirSubdirOf checks if subdir is a subdirectory of path. Does also check if subdir exists.
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
//Calls PathToAbs on path beforehand.
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

//Parent returns the path to the parent of path.
//Returns / if path is /.
func Parent(path string) (string, error) {
	dir, err := PathToAbs(path)
	if err != nil {
		return "", err
	}

	isBackslash := false
	if strings.Contains(dir, "\\") {
		isBackslash = true
		dir = strings.ReplaceAll(dir, "\\", "/")
	}

	index := strings.LastIndex(dir, "/")
	if index == -1 {
		return dir, nil
	}

	dir = dir[:index]
	if dir == "" {
		dir = "/"
	}

	if isBackslash {
		dir = strings.ReplaceAll(dir, "/", "\\")
	}

	return dir, nil
}
