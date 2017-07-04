package main

import (
	"io"
	"os"

	"github.com/zeebo/bencode"
)

// A torrent file is a bencoded (binary) file that comes in one of two flavors:
//  1. Single File
//    a. info/name: the file's name
//    b. info/length: the file's size
//  2. Multi File
//    a. info/name: the root path to hold all included files
//    b. info/files[i]/path: an array of strings, when joined is complete path and file name
//    c. info/files[i]/length: the file's size

// Torrent is a subset of the fields found in actual torrent files
type Torrent struct {
	Info TorrentInfo `bencode:"info"`
}

type TorrentInfo struct {
	Name   string             `bencode:"name"`
	Length int64              `bencode:"length"`
	Files  []TorrentInfoFiles `bencode:"files"`
}

type TorrentInfoFiles struct {
	Length int64    `bencode:"length"`
	Path   []string `bencode:"path"`
}

func ParseTorrent(in io.Reader) (*Torrent, error) {
	// parse file into decoded map
	dec := bencode.NewDecoder(in)
	var t Torrent
	if err := dec.Decode(&t); err != nil {
		return nil, err
	}

	return &t, nil
}

func OpenAndParseTorrent(filename string) (*Torrent, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ParseTorrent(file)
}
