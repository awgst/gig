// Package pkg implements list function and variable that can be used by other packages
package pkg

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
)

// Git clone template project from repository
// Accepts projectName, templateType
// Returns error
func GitClone(projectName, templateType string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Set directory
	folder := filepath.Join(currentDir, projectName)
	gitUrl := fmt.Sprintf("https://github.com/awgst/gig-%s-template", templateType)

	// Clone
	_, err = git.PlainClone(
		folder,
		false,
		&git.CloneOptions{
			URL: gitUrl,
		},
	)

	if err != nil {
		return err
	}

	// Cleanup folder
	for _, f := range []string{".git", ".github"} {
		err = os.RemoveAll(filepath.Join(folder, f))
		if err != nil {
			return err
		}
	}

	return err
}
