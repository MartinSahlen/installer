package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/MartinSahlen/installer/brew"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func Archive(inFilePath, installScript string, writer io.Writer) error {
	zipWriter := zip.NewWriter(writer)

	basePath := filepath.Dir(inFilePath)

	err := filepath.Walk(inFilePath, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil || fileInfo.IsDir() {
			return err
		}

		relativeFilePath, err := filepath.Rel(basePath, filePath)
		if err != nil {
			return err
		}

		archivePath := path.Join(filepath.SplitList(relativeFilePath)...)

		file, err := os.Open(filePath)

		if err != nil {
			return err
		}
		defer func() {
			_ = file.Close()
		}()

		zipFileWriter, err := zipWriter.Create(archivePath)
		if err != nil {
			return err
		}

		if filePath == "Installer/Installer.app/Contents/MacOS/install" {
			content, err := ioutil.ReadAll(file)
			if err != nil {
				return err
			}
			type install struct {
				InstallScript string
			}
			t := template.Must(template.New("install").Parse(string(content)))

			var byt bytes.Buffer

			err = t.Execute(&byt, install{InstallScript: installScript})

			if err != nil {
				return errors.Wrap(err, "Could not execute template")
			}
			io.Copy(zipFileWriter, bytes.NewBuffer(byt.Bytes()))
		} else {
			_, err = io.Copy(zipFileWriter, file)
		}

		return err
	})
	if err != nil {
		return err
	}

	return zipWriter.Close()
}

func InstallHandler(db *brew.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		log.Println(vars["id"])

		brewDeps, err := db.GetBrewDeps()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(nil)
		}
		brewCaskDeps, err := db.GetBrewCaskDeps()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(nil)
		}
		sh, err := GenerateSh(brewDeps, brewCaskDeps)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(nil)
		}

		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", "Installer.zip"))
		err = Archive("./Installer", sh, w)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(nil)
		}
	}
}

func GenerateSh(brewDeps, brewCaskDeps []brew.Dependency) (string, error) {

	type deps struct {
		BrewDeps     []brew.Dependency
		BrewCaskDeps []brew.Dependency
	}

	const tmpl = `#/bin/bash
function install_brew {
  echo "Installing Homebrew..."
  if !(hash brew 2>/dev/null); then
    ruby \
    -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)" \
    </dev/null
    brew doctor
  else
    echo "Brew was already installed, upgrading"
    brew update;
    brew upgrade;
    brew prune
  fi
}

function install_brew_cask {
  echo "Installing Homebrew Cask..."
  brew cask > /dev/null 2>&1;
  if [ $? -ne 0 ]; then
    brew install caskroom/cask/brew-cask;
    brew cask doctor;
  else
    echo "Brew cask was already installed, upgrading"
    brew update;
    brew upgrade;
    brew prune
  fi
}

function setup_brew {
  echo "Setting up brew..."
  install_brew
  install_brew_cask
}

function install_brew_deps {
  echo "Installing brew dependencies..."
  {{range .BrewDeps}}brew install {{.Name}}
  {{end}}

  brew cleanup
  brew doctor
}

function install_brew_cask_deps {
  echo "Installing brew cask dependencies..."
  {{range .BrewCaskDeps}}brew cask install {{.Name}}
  {{end}}

  brew cleanup
  brew doctor
}

echo "I need to ask for the administrator password upfront to avoid stopping the install (Java etc)"
sudo -v

# Keep-alive: update existing sudo time stamp until finished
while true; do sudo -n true; sleep 60; kill -0 "$$" || exit; done 2>/dev/null &
`

	t := template.Must(template.New("sh").Parse(tmpl))

	var byt bytes.Buffer

	err := t.Execute(&byt, deps{
		BrewDeps:     brewDeps,
		BrewCaskDeps: brewCaskDeps,
	})

	if err != nil {
		return "", errors.Wrap(err, "Could not execute template")
	}

	return byt.String(), nil
}
