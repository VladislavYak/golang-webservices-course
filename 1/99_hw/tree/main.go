package main

import (
	"io"
	"log"
	"os"
	pathLib "path"
)

func main() {
	out := os.Stdout

	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}

	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(writer io.Writer, path string, isPrintFiles bool) error {
	dirs, err := os.ReadDir(path)

	if err != nil {
		return err
	}
	// processDirectory(path)

	baseDir := Directory{}
	baseDir.name = pathLib.Dir(path)
	for _, val := range dirs {
		if val.IsDir() {
			dir, err := processDirectory(pathLib.Join(path, val.Name()))
			if err != nil {
				panic("panic")
			}
			baseDir.dirs = append(baseDir.dirs, dir)
		} else {
			baseDir.files = append(baseDir.files, val.Name())
		}
	}

	log.Printf("baseDir: %+v\n", baseDir)

	// src := strings.NewReader(path)w
	// if _, err := io.Copy(writer, src); err != nil {
	// 	return err
	// }
	return nil
}

func processDirectory(RouterPath string) (Directory, error) {
	log.Println("PROCESSING FOLLOWING PATH:", RouterPath)
	dirs, err := os.ReadDir(RouterPath)
	log.Println("FOUNEDED FILES&FOLDERS:", dirs)

	if err != nil {
		return Directory{}, err
	}

	d := Directory{}
	d.name = pathLib.Base(RouterPath)
	for _, val := range dirs {
		if val.IsDir() {
			d.dirs = append(d.dirs, Directory{name: val.Name()})
			log.Println("---")
			dir, _ := processDirectory(pathLib.Join(RouterPath, val.Name()))

			d.dirs = append(d.dirs, dir)
		} else {
			d.files = append(d.files, val.Name())
		}
	}

	return d, nil
}

type Directory struct {
	name  string
	files []string
	dirs  []Directory
}
