package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	regularBranch = "├───"
	lastBranch    = "└───"
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

func walkTree(out io.Writer, dirpath string, printFiles bool, indent string) error {
	files, err := ioutil.ReadDir(dirpath)
	filteredFiles := filterFiles(files, printFiles)

	for i, file := range filteredFiles {
		fileName := file.Name()
		fullpath := filepath.Join(dirpath, fileName)
		isDir := file.IsDir()

		var isLast bool
		if i == len(filteredFiles)-1 {
			isLast = true
		}

		var output string

		if isLast == false {
			output = formatBranch(indent, regularBranch, formatFilename(fileName, file))
		} else {
			output = formatBranch(indent, lastBranch, formatFilename(fileName, file))
		}

		io.WriteString(out, output)
		if isDir == true {
			var newIndent string
			if isLast == true {
				newIndent = indent + "\t"
			} else {
				newIndent = indent + "│\t"
			}

			walkTree(out, fullpath, printFiles, newIndent)
		}
	}
	return err
}

func filterFiles(files []os.FileInfo, printFiles bool) (filteredFiles []os.FileInfo) {
	if printFiles == true {
		return files
	}
	for _, file := range files {
		if file.IsDir() == true {
			filteredFiles = append(filteredFiles, file)
		}
	}
	return filteredFiles
}

func dirTree(out io.Writer, dirpath string, printFiles bool) error {
	err := walkTree(out, dirpath, printFiles, "")
	return err
}

func formatBranch(indent string, branch string, fileName string) string {
	return fmt.Sprintf("%s%s%s\n", indent, branch, fileName)
}

func formatFilename(fileName string, info os.FileInfo) string {
	if info.IsDir() == true {
		return fileName
	}

	var size string
	if s := info.Size(); s == 0 {
		size = "empty"
	} else {
		size = fmt.Sprintf("%db", s)
	}

	return fmt.Sprintf("%s (%s)", fileName, size)
}
