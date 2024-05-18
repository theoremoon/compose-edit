package main

import (
	"log"
	"os"

	composeedit "github.com/theoremoon/compose-edit"
)

func main() {
	app := composeedit.App()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
