//go:generate sh -c "go build -ldflags=\"-X main.version=$(git branch --show-current)\""
package main

import "log"

var version = "undefined"

func main() {
	log.SetFlags(0)

	handleOptions()

	initMissingPackage()

	loadTargetModule()

	updateCalls()

	generateDecls()

	writeFiles()
}
