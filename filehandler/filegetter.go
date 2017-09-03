package filehandler

import (
	"os"
	"path/filepath"
)

//GetFiles returns a slice of fileinfos.
//These are the fileinfos of all sublements of the rootpath path
func GetFiles(path string) ([]*os.FileInfo, error) {
	var headers []*os.FileInfo

	visit := func(path string, info os.FileInfo, err error) error {
		fileInf, err := os.Stat(path)
		if err != nil {
			return err
		}
		headers = append(headers, &fileInf)
		return nil
	}

	err := filepath.Walk(path, visit)
	if err != nil {
		return nil, err
	}

	return headers, nil
}
