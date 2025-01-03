package main

import (
	"io"
	"os"
	pathLib "path"
	"strconv"
	"strings"
)

var res string

func main() {
	// v := "xuipuzda"
	// v2, _ := AddBytesInfoToFilename(v, "/Users/vi/personal_proj/golang_web_services_2024-04-26/1/99_hw/tree/testdata/zline/lorem/dolor.txt")
	// fmt.Println("v", v)
	// fmt.Println("v2", v2)

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
	DirPrefix := ""

	_, err := processDirectory(path, isPrintFiles, DirPrefix)
	if err != nil {
		return err
	}
	src := strings.NewReader(res)
	if _, err := io.Copy(writer, src); err != nil {
		return err
	}
	return nil
}

func processDirectory(RouterPath string, isPrintFiles bool, DirPrefix string) (Directory, error) {
	// log.Println("PROCESSING FOLLOWING PATH:", RouterPath)
	dirs, err := os.ReadDir(RouterPath)
	var new_dirs []os.DirEntry
	for _, dir := range dirs {
		if !dir.IsDir() && !isPrintFiles {
			continue
		} else {
			new_dirs = append(new_dirs, dir)
		}
	}
	dirs = new_dirs
	// log.Println("FOUNEDED FILES&FOLDERS:", dirs)

	if err != nil {
		return Directory{}, err
	}

	d := Directory{name: pathLib.Base(RouterPath)}
	for i, val := range dirs {
		if val.IsDir() {

			if i == len(dirs)-1 {
				res += DirPrefix + "└───" + val.Name() + "\n"
			} else {
				res += DirPrefix + "├───" + val.Name() + "\n"
			}

			postFix := ""
			if i == len(dirs)-1 {
				postFix = ""
			} else {
				postFix = "│"
			}
			dir, err := processDirectory(pathLib.Join(RouterPath, val.Name()), isPrintFiles, DirPrefix+postFix+"\t")
			if err != nil {
				return Directory{}, nil
			}
			d.dirs = append(d.dirs, dir)

		} else if isPrintFiles {
			filename, err := AddBytesInfoToFilename(val.Name(), pathLib.Join(RouterPath, val.Name()))
			if err != nil {
				return Directory{}, err
			}

			if i == len(dirs)-1 {
				res += DirPrefix + "└───" + filename + "\n"
			} else {
				res += DirPrefix + "├───" + filename + "\n"
			}
			d.files = append(d.files, filename)
		}
	}

	return d, nil
}

func AddBytesInfoToFilename(filename, path string) (string, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	// get the size
	size := fi.Size()
	if size == 0 {
		filename = filename + " (empty)"
	} else {
		filename = filename + " (" + strconv.Itoa(int(size)) + "b)"
	}
	return filename, nil
}

type Directory struct {
	name  string
	files []string
	dirs  []Directory
}
