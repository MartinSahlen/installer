package brew

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

const (
	brewFile     string = "gen/brew.json"
	brewCaskFile string = "gen/brew-cask.txt"
)

type DB struct {
	brewDeps     []Dependency
	brewCaskDeps []Dependency
	store        map[string][]byte
}

func NewDB() (*DB, error) {
	brewDeps, err := getBrewDeps()
	if err != nil {
		return nil, errors.Wrap(err, "Could not get brew dependencies")
	}
	brewCaskDeps, err := getBrewCaskDeps()
	if err != nil {
		return nil, errors.Wrap(err, "Could not get brew cask dependencies")
	}
	return &DB{
		brewDeps:     brewDeps,
		brewCaskDeps: brewCaskDeps,
		store:        map[string][]byte{},
	}, nil
}

func (db *DB) GetBrewDeps() ([]Dependency, error) {
	return db.brewCaskDeps, nil
}

func (db *DB) GetBrewCaskDeps() ([]Dependency, error) {
	return db.brewCaskDeps, nil
}

func (db *DB) GetAllDeps() ([]Dependency, error) {
	return append(db.brewDeps, db.brewCaskDeps...), nil
}

func (db *DB) GetDepsByFullNames(fullNames []string) ([]Dependency, error) {
	deps, _ := db.GetAllDeps()
	found := []Dependency{}
	for _, fullName := range fullNames {
		for _, dep := range deps {
			if fullName == dep.FullName {
				found = append(found, dep)
			}
		}
	}
	return found, nil
}

func (db *DB) GetDepsByNames(names []string) ([]Dependency, error) {
	deps, _ := db.GetAllDeps()
	found := []Dependency{}
	for _, name := range names {
		for _, dep := range deps {
			if name == dep.Name {
				found = append(found, dep)
			}
		}
	}
	return found, nil
}

func (db *DB) SaveDependencies(deps []Dependency) (string, error) {

	byt, err := json.Marshal(deps)
	if err != nil {
		return "", errors.Wrap(err, "Could not save deps")
	}
	id := uuid.NewV4().String()
	db.store[id] = byt
	return id, nil
}

func (db *DB) GetDependenciesForID(id string) ([]Dependency, error) {

	byt, ok := db.store[id]
	if !ok {
		return nil, errors.New("Could not retrieve deps for id " + id)
	}
	deps := []Dependency{}
	err := json.Unmarshal(byt, &deps)
	if err != nil {
		return nil, errors.Wrap(err, "Could not get deps from DB")
	}
	return deps, nil
}

func getBrewDeps() ([]Dependency, error) {
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

func getBrewCaskDeps() ([]Dependency, error) {
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
