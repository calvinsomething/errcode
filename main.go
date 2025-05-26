package main

import (
	"log"
	"os"
)

func main() {
	log.SetFlags(0)

	handleOptions()

	if initMissingPackage() {
		os.Exit(0)
	}

	loadInitialExportedIdents()

	loadTargetModule()

	updateCalls()

	generateDecls()

	writeFiles()
}
