package utils

import (
	"os"
	"path/filepath"
)

//GetFileInfos returns a slice of Fileinfos containing the fileinfos of all subelements of path.
func GetFileInfos(path string) ([]*os.FileInfo, error) {
	headers := make([]*os.FileInfo, 0)

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
