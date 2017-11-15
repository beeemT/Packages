package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

)

//SettingsError is either a Load Or a Store Error.
//it capsulates the underlying error.
type SettingsError struct {
	name string
	v    interface{}
	err  error
}

//LoadDataFromJSON currently the same as direct Unmarshal call.
//Might change in the future, thus this function exists
func LoadDataFromJSON(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return &SettingsError{"LoadError: Couldnt load JSON from data", nil, err}
	}
	return nil
}

//LoadDataFromJSONFile loads the content of a file specified by path.
//Content is loaded into v. Returns a SettingsError on failure.
func LoadDataFromJSONFile(path string, v interface{}) error {
	if !IsFile(path) {
		return &SettingsError{"PathError: Path does either not exist or is not a file", path, nil}
	}

	settings, err := ioutil.ReadFile(path)
	if err != nil {
		return &SettingsError{"LoadError: Couldnt load file", path, err}
	}

	err = LoadDataFromJSON(settings, v)
	if err != nil {
		return err //err is already a SettingsError
	}
	return nil
}

func StoreSettingsToJSON() {

}

func (s *SettingsError) Error() string {
	return fmt.Sprintf("%s - %s: %s\n", s.name, s.v, s.err.Error())
}
