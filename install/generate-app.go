package install

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/MartinSahlen/installer/brew"
	"github.com/pkg/errors"
)

func ArchiveInstallApp(inFilePath, installScript string, writer io.Writer) error {
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

func GenerateInstallScript(db *brew.DB, configID string) (string, error) {
	type deps struct {
		BrewDeps     []brew.Dependency
		BrewCaskDeps []brew.Dependency
	}

	d, err := db.GetDependenciesForID(configID)

	if err != nil {
		return "", nil
	}

	brewDeps := []brew.Dependency{}
	brewCaskDeps := []brew.Dependency{}

	for _, dep := range d {
		if dep.Type == brew.Brew {
			brewDeps = append(brewDeps, dep)
		}
		if dep.Type == brew.BrewCask {
			brewCaskDeps = append(brewCaskDeps, dep)
		}
	}

	const tmpl = `function install_brew {
  echo \"Installing Homebrew...\"
	/usr/local/bin/brew > /dev/null 2>&1;
  if [ $? -ne 0 ]; then
    ruby \
    -e \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)\" \
    </dev/null
    /usr/local/bin/brew doctor
  else
    echo \"Brew was already installed, upgrading\"
    /usr/local/bin/brew update
    /usr/local/bin/brew upgrade
    /usr/local/bin/brew prune
  fi
}

function install_brew_cask {
  echo \"Installing Homebrew Cask...\"
  /usr/local/bin/brew cask > /dev/null 2>&1;
  if [ $? -ne 0 ]; then
    /usr/local/bin/brew install caskroom/cask/brew-cask
    /usr/local/bin/brew cask doctor
  else
    echo \"Brew cask was already installed, upgrading\"
    /usr/local/bin/brew update
    /usr/local/bin/brew upgrade
    /usr/local/bin/brew prune
  fi
}

function setup_brew {
  echo \"Setting up brew...\"
  install_brew
  install_brew_cask
}

function install_brew_deps {
  echo \"Installing brew dependencies...\"
  {{range .BrewDeps}}/usr/local/bin/brew install {{.Name}}
  {{end}}
  /usr/local/bin/brew cleanup
  /usr/local/bin/brew doctor
}

function install_brew_cask_deps {
  echo \"Installing brew cask dependencies...\"
  {{range .BrewCaskDeps}}/usr/local/bin/brew cask install {{.Name}}
  {{end}}
  /usr/local/bin/brew cleanup
  /usr/local/bin/brew doctor
}

setup_brew
install_brew_deps
install_brew_cask_deps`

	t := template.Must(template.New("sh").Parse(tmpl))

	var byt bytes.Buffer

	err = t.Execute(&byt, deps{
		BrewDeps:     brewDeps,
		BrewCaskDeps: brewCaskDeps,
	})

	if err != nil {
		return "", errors.Wrap(err, "Could not execute template")
	}

	return byt.String(), nil
}
