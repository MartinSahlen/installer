package main

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

const (
	brewFile     string = "gen/brew.json"
	brewCaskFile string = "gen/brew-cask.txt"
)

func GetBrewDeps() ([]Dependency, error) {
	deps := []Dependency{}

	byt, err := ioutil.ReadFile(brewFile)

	if err != nil {
		return nil, errors.Wrap(err, "Could not read "+brewFile)
	}

	err = json.Unmarshal(byt, &deps)

	if err != nil {
		return nil, errors.Wrap(err, "Could parse json for "+brewFile)
	}

	for _, dep := range deps {
		dep.Type = Brew
	}

	return deps, nil
}

func GetBrewCaskDeps() ([]Dependency, error) {
	byt, err := ioutil.ReadFile(brewCaskFile)

	if err != nil {
		errors.Wrap(err, "Could not read "+brewCaskFile)
	}

	casks := strings.Split(string(byt), "\nSPLIT_HERE_PLEASE\n")

	deps := []Dependency{}

	for _, cask := range casks {
		dep := Dependency{}
		didSetName := false
		for _, line := range strings.Split(cask, "\n") {
			pair := strings.SplitN(strings.TrimSpace(line), " ", 2)
			if len(pair) >= 2 {
				key := pair[0]
				val := strings.Replace(pair[1], "\"", "", -1)
				val = strings.Replace(val, "'", "", -1)
				if key == "homepage" {
					dep.HomePage = val
				}
				if key == "name" && !didSetName {
					dep.FullName = val
					didSetName = true
				}
				if key == "cask" {
					dep.Name = strings.Split(val, " ")[0]
				}
			}
			dep.Type = BrewCask
		}
		deps = append(deps, dep)
	}
	return deps, nil
}
