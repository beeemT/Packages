package jsonutil

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/grekhor/Packages/utils"
)

//JSONError is either a Load Or a Store Error.
//it capsulates the underlying error.
type JSONError struct {
	name string
	v    interface{}
	err  error
}

func (s *JSONError) Error() string {
	return fmt.Sprintf("%s - %s: %s\n", s.name, s.v, s.err.Error())
}

//LoadDataFromJSON uses ioutil.ReadAll to read the byte slice from dataReader and passes the slice to json.Unmarshal.
//It does not close the dataReader.
//Only returns errors of type jsonutil.JSONError.
func LoadDataFromJSON(dataReader io.Reader, v interface{}) error {
	data, err := ioutil.ReadAll(dataReader)
	err = json.Unmarshal(data, &v)
	if err != nil {
		return &JSONError{"LoadError: Couldn't load JSON from data", nil, err}
	}
	return nil
}

//LoadDataFromCompressedJSON converts the reader into a gzip.Reader and passes it to LoadDataFromJSON.
//Only closes the gzipReader, not the dataReader.
//Only returns errors of type jsonutil.JSONError.
func LoadDataFromCompressedJSON(dataReader io.Reader, v interface{}) error {
	gzipReader, err := gzip.NewReader(dataReader)
	if err != nil {
		return &JSONError{"LoadError: Failed to decompress Reader", nil, err}
	}
	defer gzipReader.Close()

	err = LoadDataFromJSON(gzipReader, v)
	return err
}

//Checks if path exists and is a file. Opens the file and returns it.
//File has to be closed by the caller.
//Only returns jsonutil.JSONErrors.
func createReaderFromFile(path string) (*os.File, error) {
	path, err := utils.PathToAbs(path)
	if err != nil {
		return nil, &JSONError{"PathError: Failed to make path absolute", path, err}
	}

	if !utils.IsFile(path) {
		return nil, &JSONError{"PathError: Path does either not exist or is not a file", path, nil}
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, &JSONError{"LoadError: Couldnt open file", path, err}
	}

	return file, nil
}

//LoadDataFromJSONFile loads the content of a file specified by path.
//Makes path absolute
//Content is loaded into v.
//Returns a JSONError on failure.
func LoadDataFromJSONFile(path string, v interface{}) error {
	file, err := createReaderFromFile(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = LoadDataFromJSON(file, v)
	return err
}

//LoadDataFromCompressedJSONFile loads the content of a file specified by path and passes the file to LoadDataFromCompressedJSON.
//Makes path absolute.
//Closes all Readers. Only returns errors of type jsonutil.JSONError.
func LoadDataFromCompressedJSONFile(path string, v interface{}) error {
	file, err := createReaderFromFile(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return LoadDataFromCompressedJSON(file, v)
}

//StoreDataToJSON is the same as json.Marshal but only returns errors of type jsonutil.JSONError.
func StoreDataToJSON(v interface{}) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, &JSONError{"StoreError: Couldn't store data", nil, err}
	}
	return b, nil
}

//StoreDataToCompressedJSON passes v to StoreDataToJSON.
//The JSON is then compressed into a buffer and buf.Bytes() is returned.
//Only returns errors of type jsonutil.JSONError.
func StoreDataToCompressedJSON(v interface{}, compressionLevel int) ([]byte, error) {
	b, err := StoreDataToJSON(v)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(make([]byte, 0, cap(b)))
	gzipWriter, err := gzip.NewWriterLevel(buf, compressionLevel)
	if err != nil {
		return nil, &JSONError{"StoreError: Failed to create compressor for data", nil, err}
	}
	defer gzipWriter.Close()

	_, err = gzipWriter.Write(b)
	if err != nil {
		return nil, &JSONError{"StoreError: Failed to compress data", nil, err}
	}

	return buf.Bytes(), nil
}

func prepareJSONWrite(path string, overwrite bool) error {
	path, err := utils.PathToAbs(path)
	if err != nil {
		return &JSONError{"PathError: Failed to make path absolute", path, err}
	}

	exist, err := utils.Exists(path)
	if err != nil {
		return &JSONError{"PathError: Failed to check if file exists", path, err}
	}

	if exist && !overwrite {
		return &JSONError{"StoreError: File already exists and overwrite flag is not set", path, os.ErrExist}
	}
	return nil
}

//StoreDataToJSONFile converts the content of v into JSON if possible and stores the content into a file at path.
//Makes path absolute.
//If path already exists and the overwrite flag is not set, an error is returned.
func StoreDataToJSONFile(path string, v interface{}, overwrite bool) error {
	err := prepareJSONWrite(path, overwrite)
	if err != nil {
		return err
	}

	b, err := StoreDataToJSON(v)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		return &JSONError{"StoreError: Failed to write data to file", path, err}
	}
	return nil
}

//StoreDataToCompressedJSONFile converts the content of v into JSON if possible,
//compresses the byte slice using compressionLevel (gzip.constants) and stores the content into a file at path.
//Makes path absolute.
//If path already exists and the overwrite flag is not set, an error is returned.
func StoreDataToCompressedJSONFile(path string, v interface{}, overwrite bool, compressionLevel int) error {
	err := prepareJSONWrite(path, overwrite)
	if err != nil {
		return err
	}

	b, err := StoreDataToCompressedJSON(v, compressionLevel)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		return &JSONError{"StoreError: Failed to write data to file", path, err}
	}
	return nil
}
