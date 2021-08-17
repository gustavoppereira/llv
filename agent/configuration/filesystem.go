package configuration

import (
	"io/fs"
	"io/ioutil"
)

func readFile(filepath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func saveFile(content []byte, filepath string) error {
	err := ioutil.WriteFile(filepath, content, fs.ModeAppend)
	return err
}
