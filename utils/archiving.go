package utils

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
)

//CreateZipOfFiles takes a path and zips all subelements. returns the zip
func CreateZipOfFiles(path string) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)

	zipWriter := zip.NewWriter(buffer)

	rootSubElems, err := GetFileInfos(path)
	if err != nil {
		return nil, err
	}

	for _, fileInf := range rootSubElems {
		fileInfHeader, err := zip.FileInfoHeader(*fileInf)
		if err != nil {
			return nil, err
		}

		zipFileWriter, err := zipWriter.CreateHeader(fileInfHeader)
		if err != nil {
			return nil, err
		}

		bytes, err := ioutil.ReadFile(fileInfHeader.Name)
		if err != nil {
			return nil, err
		}

		_, err = zipFileWriter.Write(bytes)
		if err != nil {
			return nil, err
		}
	}

	err = zipWriter.Close()
	if err != nil {
		return nil, err
	}
	return buffer, nil
}
