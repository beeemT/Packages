package filehandler

import (
	"archive/zip"
	"bytes"
	"os"
)

//WriteZip takes a path and zips all subelements. returns the zip
func WriteZip(path string) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)

	zipWriter := zip.NewWriter(buffer)
	defer zipWriter.Close()

	rootSubElems, err := getFiles(path)
	if err != nil {
		return nil, err
	}

	for _, fileInf := range rootSubElems {
		fileInfHeader, err := zip.FileInfoHeader(fileInf)
		if err != nil {
			return nil, err
		}
		
		zipFileWriter, err := zipWriter.CreateHeader(fileInfHeader)
		if err != nil {
			return nil, err
		}

		file, err := os.Open(fileInfHeader.Name)
		if err != nil {
			return nil, err
		}

		var bytes []byte
		_, err = file.Read(bytes)
		if err != nil {
			return nil, err
		}

		_, err = zipFileWriter.Write(bytes)
		if err != nil {
			return nil, err
		}
	}

	return buffer, nil
}
