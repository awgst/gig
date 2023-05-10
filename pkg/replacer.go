package pkg

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

var oldString string
var newString string

// Visit all file while filepath.Walk and replace oldString with the newString
func visit(path string, fi os.FileInfo, err error) error {

	if err != nil {
		return err
	}

	if fi.IsDir() {
		return nil
	}

	matched, err := filepath.Match("*", fi.Name())

	if err != nil {
		return err
	}

	if matched {
		read, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		newContents := strings.Replace(string(read), oldString, newString, -1)

		err = os.WriteFile(path, []byte(newContents), 0)
		if err != nil {
			return err
		}

	}

	return nil
}

func Replace(path, old, new string) {
	oldString = old
	newString = new
	err := filepath.Walk(path, visit)
	if err != nil {
		log.Fatal(err)
	}
}
