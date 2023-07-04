package file

import (
	"io/ioutil"
	"strings"
)

// GetNamesFromFile reads the input names file and returns a slice of names
func GetNamesFromFile(namesFile string) ([]string, error) {
	namesData, err := ioutil.ReadFile(namesFile)
	if err != nil {
		return nil, err
	}

	names := strings.Split(string(namesData), "\n")
	return names, nil
}
