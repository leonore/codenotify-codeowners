package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

var pathToWalk string

func walkDirectoryForCodenotify(ctx *cli.Context) error {
	if pathToWalk == "" {
		return errors.New("path is required")
	}

	ownersFile, err := os.Create("CODEOWNERS")
	if err != nil {
		return err
	}
	defer ownersFile.Close()
	w := bufio.NewWriter(ownersFile)

	if err := filepath.WalkDir(pathToWalk, func(path string, d os.DirEntry, err error) error {
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
			reformatted := filepath.Dir(path)[len(pathToWalk):] + "/" + p.line + "\n"
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
		Usage: "convert codenotify files to single codeowners file.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "path",
				Value:       "",
				Usage:       "directory path to walk",
				Destination: &pathToWalk,
				Required:    true, // required because we need to remove it as prefix
			},
		},
		Action: walkDirectoryForCodenotify,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
