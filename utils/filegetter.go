package utils

import (
	"os"
	"path/filepath"
)

//GetFileInfos returns a slice of fileinfos.
//These are the fileinfos of all sublements of the rootpath path
func GetFileInfos(path string) ([]*os.FileInfo, error) {
	var headers []*os.FileInfo

	visit := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		headers = append(headers, &info)
		return nil
	}

	err := filepath.Walk(path, visit)
	if err != nil {
		return nil, err
	}

	return headers, nil
}
