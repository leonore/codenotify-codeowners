package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func walkDirectoryForCodenotify(ctx *cli.Context) error {
	fmt.Println("works")
	return nil
}

func main() {
	app := &cli.App{
		Name:   "codenotify -> codeowners",
		Usage:  "convert codenotify files to single codeowners file",
		Action: walkDirectoryForCodenotify,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
