package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

var pathToWalk string

func walkDirectoryForCodenotify(ctx *cli.Context) error {
	// todo handle pathToWalk
	var cwd string
	if pathToWalk != "" {
		cwd = pathToWalk
	} else {
		tmpCwd, err := os.Getwd()
		if err != nil {
			return err
		}
		cwd = tmpCwd
	}

	if err := filepath.WalkDir(cwd, func(path string, d os.DirEntry, err error) error {
		// if file is codenotify
		// go through each line in file
		// skip comments and empty
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, "CODENOTIFY") {
			return nil
		}
		codeNotifyFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer codeNotifyFile.Close()

		scanner := bufio.NewScanner(codeNotifyFile)
		p := new(parsing)

		for scanner.Scan() {
			p.nextLine(scanner.Text())
			if p.isBlank() {
				continue
			}
			fmt.Println(p.line)
		}

		if err := scanner.Err(); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func main() {
	app := &cli.App{
		Name:  "codenotify -> codeowners",
		Usage: "convert codenotify files to single codeowners file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "path",
				Value:       "",
				Usage:       "path to parse files for",
				Destination: &pathToWalk,
			},
		},
		Action: walkDirectoryForCodenotify,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
