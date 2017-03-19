package main

import "log"

func main() {
	//deps, err := GetAllDeps()

	//parseError(err)

	brewDeps, _ := GetBrewDeps()
	brewCaskDeps, _ := GetBrewCaskDeps()

	for _, brewDep := range brewDeps {
		for _, brewCaskDep := range brewCaskDeps {
			if brewDep.Name == brewCaskDep.Name {
				log.Println("brew: " + brewDep.Name + " " + brewDep.HomePage)

				log.Println("cask: " + brewCaskDep.Name + " " + brewCaskDep.HomePage)
			}
		}
	}

}
