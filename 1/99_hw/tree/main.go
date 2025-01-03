package main

import (
	"fmt"
	"io"
	"os"
	pathLib "path"
	"strings"
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
	FilePrefix := ""
	DirPrefix := ""

	_, err := processDirectory(path, isPrintFiles, FilePrefix, DirPrefix)
	if err != nil {
		return err
	}
	src := strings.NewReader("xui pizda ebana")
	if _, err := io.Copy(writer, src); err != nil {
		return err
	}
	return nil
}

func processDirectory(RouterPath string, isPrintFiles bool, FilePrefix string, DirPrefix string) (Directory, error) {
	// log.Println("PROCESSING FOLLOWING PATH:", RouterPath)
	dirs, err := os.ReadDir(RouterPath)
	// log.Println("FOUNEDED FILES&FOLDERS:", dirs)

	if err != nil {
		return Directory{}, err
	}

	d := Directory{name: pathLib.Base(RouterPath)}
	for i, val := range dirs {
		if val.IsDir() {

			if i == len(dirs)-1 {
				fmt.Println(DirPrefix + "└───" + val.Name())
			} else {
				fmt.Println(DirPrefix + "├───" + val.Name())
			}

			postFix := ""
			if i == len(dirs)-1 {
				postFix = ""
			} else {
				postFix = "|"
			}
			dir, err := processDirectory(pathLib.Join(RouterPath, val.Name()), isPrintFiles, FilePrefix, DirPrefix+postFix+"\t")
			if err != nil {
				return Directory{}, nil
			}
			d.dirs = append(d.dirs, dir)

		} else if isPrintFiles {

			if i == len(dirs)-1 {
				fmt.Println(DirPrefix + "└───" + val.Name())
			} else {
				fmt.Println(DirPrefix + "├───" + val.Name())
			}
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
