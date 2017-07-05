package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/danbrakeley/torc/eff"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "torc"
	app.Version = "0.0.1"
	app.Copyright = "(c) 2017 Dan Brakeley"
	app.Usage = "compare local files to a torrent file"

	app.Commands = []cli.Command{
		{
			Name:    "compare",
			Aliases: []string{"cmp"},
			Usage:   "compares a torrent file with actual files on disk",
			Description: `Compares the files listed inside a torrent file with those present on disk. Any
differences are output to stdout. The optional --delete flag will go ahead and delete
the files on disk that are not present in the torrent.`,
			ArgsUsage: "<file.torrent> <root/path>",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "delete",
					Usage: "delete any files on disk that are not also in the torrent",
				},
				cli.StringFlag{
					Name:  "root",
					Usage: "change the root path to the given string",
				},
			},
			Action: func(c *cli.Context) error {
				args := c.Args()
				if err := CheckNumArgs(len(args), 2); err != nil {
					return err
				}
				shouldDelete := c.Bool("delete")
				opts := CompareOpts{
					TorrentRootOverride: c.String("root"),
				}
				return WrapInUsageError(DoCompareAction(args[0], args[1], &opts, shouldDelete))
			},
		},
		{
			Name:      "list",
			Aliases:   []string{"ls"},
			Usage:     "list all files for a given torrent",
			ArgsUsage: "<file.torrent>",
			Action: func(c *cli.Context) error {
				args := c.Args()
				if err := CheckNumArgs(len(args), 1); err != nil {
					return err
				}
				return WrapInUsageError(DoListAction(args[0]))
			},
		},
	}

	app.Run(os.Args)
}

type UsageError struct {
	error
}

func (e UsageError) ExitCode() int {
	return -1
}

func NewUsageError(format string, a ...interface{}) error {
	return UsageError{errors.New(fmt.Sprintf(format, a...))}
}

func WrapInUsageError(err error) error {
	if err == nil {
		return nil
	}
	return UsageError{err}
}

func CheckNumArgs(numArgs, expected int) error {
	if numArgs < expected {
		return eff.NewMsg("not enough arguments")
	}
	if numArgs > expected {
		return eff.NewMsg("too many arguments")
	}
	return nil
}

func DoCompareAction(torrentFile, rootPath string, opts *CompareOpts, shouldDelete bool) error {
	r, err := CompareTorrentPathsToDisk(torrentFile, rootPath, opts)
	if err != nil {
		return err
	}

	// print results
	if len(r.OnlyInTorrent) > 0 {
		fmt.Println("Files only in torrent (missing from disk):")
		for _, entry := range r.OnlyInTorrent {
			fmt.Println("- " + entry)
		}
	}
	if len(r.OnlyOnDisk) > 0 {
		fmt.Println("Files only on disk (missing in torrent):")
		for _, entry := range r.OnlyOnDisk {
			fmt.Println("+ " + entry)
		}
	}
	if len(r.OnlyInTorrent) == 0 && len(r.OnlyOnDisk) == 0 {
		fmt.Println("No differences found")
	}

	if shouldDelete && len(r.OnlyOnDisk) > 0 {
		yesToAll := false
		numSkipped := 0
		numDeleted := 0
		for i, entry := range r.OnlyOnDisk {
			deletePath := filepath.Join(rootPath, entry)
			yes := yesToAll

		tryagain:
			if !yes {
				fmt.Printf("[%d/%d] Delete \"%s\" (y/n/a)? ", i+1, len(r.OnlyOnDisk), deletePath)

				var rawReply string
				_, err := fmt.Scanf("%s\n", &rawReply)
				if err != nil {
					fmt.Println("")
					return eff.NewErr(err).WithMessage("problem reading response, aborting")
				}

				switch strings.ToLower(strings.TrimSpace(rawReply)) {
				case "y", "yes":
					yes = true
				case "a", "all":
					yes = true
					yesToAll = true
				case "n", "no":
				default:
					fmt.Println("Valid responses: [y]es, [n]o, yes to [a]ll")
					goto tryagain
				}
			}

			if yes {
				if err := os.Remove(deletePath); err != nil {
					return eff.NewErr(err).WithMessage("error deleting file, aborting").WithField("path", deletePath)
				}
				numDeleted++
			} else {
				numSkipped++
			}
		}

		if numDeleted > 0 || numSkipped > 0 {
			fmt.Printf("Deleted %d file%s, skipped %d file%s\n", numDeleted, pluralS(numDeleted), numSkipped, pluralS(numSkipped))
		}
	}

	return nil
}

func pluralS(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

func DoListAction(torrentFile string) error {
	t, err := OpenAndParseTorrent(torrentFile)
	if err != nil {
		return err
	}

	// single file torrent
	if len(t.Info.Files) < 1 {
		fmt.Println(t.Info.Name)
		return nil
	}

	for _, entry := range t.Info.Files {
		fmt.Println(filepath.Join(t.Info.Name, filepath.Join(entry.Path...)))
	}
	return nil
}
