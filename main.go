package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

var pathToWalk string
var dirPrefix string

func walkDirectoryForCodenotify(ctx *cli.Context) error {
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

	ownersFile, err := os.Create("CODEOWNERS")
	if err != nil {
		return err
	}
	defer ownersFile.Close()
	w := bufio.NewWriter(ownersFile)

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
			reformatted := filepath.Dir(path)[len(dirPrefix):] + "/" + p.line + "\n"
			if _, err := w.WriteString(reformatted); err != nil {
				return err
			}
		}
		w.WriteString("\n")
		if err := scanner.Err(); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	w.Flush()
	return nil
}

func main() {
	app := &cli.App{
		Name:  "codenotify -> codeowners",
		Usage: "convert codenotify files to single codeowners file.\nWARNING: this needs to be called from the root of a repo.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "dir-prefix",
				Value:       "",
				Usage:       "directory prefix to exclude from regenerated owner path",
				Destination: &dirPrefix,
			},
			&cli.StringFlag{
				Name:        "path",
				Value:       "",
				Usage:       "directory path to walk",
				Destination: &pathToWalk,
			},
		},
		Action: walkDirectoryForCodenotify,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
