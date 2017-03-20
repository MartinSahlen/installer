package main

import (
	"bytes"
	"log"
	"net/http"
	"text/template"

	"github.com/MartinSahlen/installer/brew"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func ShHandler(db *brew.DB) http.HandlerFunc {
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
		_, err = w.Write([]byte(sh))
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
