package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/troylelandshields/govis/fGraph"
	"github.com/troylelandshields/govis/parser"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	//Get current directory
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	fG := fGraph.NewFunctionGraph()
	dones := []chan bool{}

	//walk the filesystem.
	walkFunc := func(path string, info os.FileInfo, err error) error {
		//Skip .git directory.
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		//Parse if file is .go file.
		extension := filepath.Ext(path)
		if strings.ToLower(extension) == ".go" {

			var temp []chan bool
			fG, temp = parser.ParseFile(path, nil, fG)

			dones = append(dones, temp...)
		}

		return nil
	}

	filepath.Walk(dir, walkFunc)

	for _, d := range dones {
		d <- true
	}
	
	d, err := fG.ToJSON()
	
	if err != nil {
		log.Printf("Error: [%s]\n", err)
	}
	
	os.Stdout.Write(d)
}
