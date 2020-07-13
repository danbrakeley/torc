package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// CompareOpts are optional flags that can be passed to CompareTorrentPathsToDisk
type CompareOpts struct {
	// TorrentRootOverride, if not empty, will override the base path name that is normally set in the torrent file
	TorrentRootOverride string
}

// Report is the result of a Compare
type Report struct {
	IsSingleFile  bool
	OnlyOnDisk    []string
	OnlyInTorrent []string
}

// Compare opens the given torrent file and compiles a report of any differences with the given path on disk
func CompareTorrentPathsToDisk(torrentFile, path string, opts *CompareOpts) (*Report, error) {
	t, err := OpenAndParseTorrent(torrentFile)
	if err != nil {
		return nil, err
	}

	r := Report{}

	// a single file torrent only needs one on-disk check
	if len(t.Info.Files) < 1 {
		r.IsSingleFile = true
		nativePath := filepath.Join(path, t.Info.Name)
		if _, err := os.Stat(nativePath); os.IsNotExist(err) {
			r.OnlyInTorrent = append(r.OnlyInTorrent, nativePath)
		} else if err != nil {
			return nil, fmt.Errorf(`error in os.Stat("%s"): %v`, nativePath, err)
		}
		return &r, nil
	}

	r.IsSingleFile = false

	torrentRootPath := t.Info.Name
	if opts != nil && len(opts.TorrentRootOverride) > 0 {
		torrentRootPath = opts.TorrentRootOverride
	}

	rootPath := filepath.Join(path, torrentRootPath)
	info, err := os.Stat(rootPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf(`"%s" does not exist`, rootPath)
	}
	if err != nil {
		return nil, fmt.Errorf(`error trying to stat "%s": %v`, rootPath, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf(`"%s" is not a folder`, rootPath)
	}

	// build sorted list of files on disk
	filesOnDisk := make([]string, 0, 4096)
	err = filepath.Walk(rootPath,
		func(file string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				relPath, err := filepath.Rel(path, file)
				if err != nil {
					return err
				}
				filesOnDisk = append(filesOnDisk, relPath)
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	// Note: Walk is supposed to operate in lexical order, but because we are generating a list of strings that
	// include a path separator character, filesOnDisk may not end up properly strictly sorted. So always sort.
	// The specific example that was sorted wrong was: "multi.torrent" and "multi\*".
	sort.Slice(filesOnDisk, func(i, j int) bool { return filesOnDisk[i] < filesOnDisk[j] })

	// build sorted list of files in the torrent
	filesInTorrent := make([]string, len(t.Info.Files))

	for i, entry := range t.Info.Files {
		filesInTorrent[i] = filepath.Join(torrentRootPath, filepath.Join(entry.Path...))
	}
	sort.Slice(filesInTorrent, func(i, j int) bool { return filesInTorrent[i] < filesInTorrent[j] })

	// find the differences in the two sorted slices
	r.OnlyInTorrent, r.OnlyOnDisk = DiffSortedSlices(filesInTorrent, filesOnDisk)

	return &r, nil
}
