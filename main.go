// Copyright (c) 2020, Alliance for Sustainable Energy, LLC.
// All rights reserved.

// Package main runs the formatter
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	write       = flag.Bool("w", false, "overwrite the file(s)")
	versionFlag = flag.Bool("v", false, "display tool version")
	lineLength  = flag.Int("line-length", -1, "how many characters allowed per line; -1 means no max")
	// build information added by goreleaser
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: modelicafmt [path ...]")
	flag.PrintDefaults()
}

func isModelicaFile(f os.FileInfo) bool {
	name := f.Name()
	return !f.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".mo")
}

func processAndWriteFile(filename string) {
	var b bytes.Buffer
	err := processFile(filename, bufio.NewWriter(&b), *lineLength)
	if err != nil {
		panic(err)
	}
	if *write {
		err := ioutil.WriteFile(filename, b.Bytes(), 777)
		if err != nil {
			panic(err)
		}
	} else {
		b.WriteTo(os.Stdout)
	}
}

func visitFile(filename string, f os.FileInfo, err error) error {
	if err != nil && !os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	if isModelicaFile(f) {
		processAndWriteFile(filename)
	}

	return nil
}

func walkDir(path string) {
	filepath.Walk(path, visitFile)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if *versionFlag {
		fmt.Printf("modelicafmt v%s (SHA %s)\nBuilt %s by %s\n", version, commit, date, builtBy)
		return
	}
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "error: must provide at least one file or directory")
		usage()
		os.Exit(2)
	}

	for i := 0; i < flag.NArg(); i++ {
		path := flag.Arg(i)
		switch dir, err := os.Stat(path); {
		case err != nil:
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(2)
		case dir.IsDir():
			walkDir(path)
		default:
			processAndWriteFile(path)
		}
	}
}
