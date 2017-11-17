package main

import "fmt"
import "os"
import "os/exec"
import (
	"flag"
	"strings"
	"errors"
	"archive/zip"
	"io"
	"path/filepath"
)

func handleError(err error) bool {
	if err == nil {
		return false
	}
	fmt.Println(err.Error())
	os.Exit(1)
	return true
}

//
// Returns the current HEAD hash, or an error if git can't be found.
//
func getGitHead() (string, error) {
	gitHead := exec.Command("git", "rev-parse", "HEAD")

	if gitHeadOutput, err := gitHead.Output(); err != nil {
		return "", errors.New("It appears that git cant be found in %PATH%")
	} else {
		return strings.TrimSpace(string(gitHeadOutput)), nil
	}
}

func validateFlags(firstCommitPtr, lastCommitPtr *string) bool {
	valid := true

	if len(*firstCommitPtr) == 0 {
		fmt.Println("Please provide a git commit to begin at via the -first flag.")
		valid = false
	}

	if len(*firstCommitPtr) > 40 {
		fmt.Println("The hash for -first is not valid.")
		valid = false
	}

	if len(*lastCommitPtr) == 0 {
		fmt.Println("Please provide a git commit to end at via the -last flag.")
		valid = false
	}

	if len(*lastCommitPtr) > 40 {
		fmt.Println("The hash for -last is not valid.")
		valid = false
	}

	return valid
}

//
// @link https://golangcode.com/create-zip-files-in-go/
//
func ZipFiles(filename string, baseDir string, files []string) error {
	newFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newFile.Close()

	zipWriter := zip.NewWriter(newFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		// Get the file information
		info, err := os.Stat(file)
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(file, baseDir)
		if strings.HasPrefix(header.Name, string(os.PathSeparator)) {
			header.Name = strings.TrimPrefix(header.Name, string(os.PathSeparator))
		}

		if info.IsDir() {
			header.Name += string(os.PathSeparator)
		} else {
			// Change to deflate to gain better compression
			// see http://golang.org/pkg/archive/zip/#pkg-constants
			header.Method = zip.Deflate
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			continue
		}

		zipfile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer zipfile.Close()

		_, err = io.Copy(writer, zipfile)
		if err != nil {
			return err
		}
	}
	return nil
}

func sliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func zipChanges(changes string, beVerbose bool) {
	slice := strings.Split(changes,"\n")
	if cwd, err := os.Getwd(); err != nil {
		handleError(err)
	} else {

		var files []string

		for _, v := range slice {
			if len(v) > 0 {
				file := filepath.Join(cwd, v)

				// Add the folder before adding the file, if the folder hasn't already been added
				fileDir := filepath.Dir(file);
				if (!sliceContains(files, fileDir) && fileDir != cwd){
					files = append(files, fileDir)
					if beVerbose {
						fmt.Println("Adding Directory: [" + fileDir + "]")
					}
				}

				if beVerbose {
					fmt.Println("Adding File: [" + file + "]")
				}
				files = append(files, file)
			}
		}

		if len(files) > 0 {
			err := ZipFiles("build.zip", cwd, files)
			if err != nil {
				handleError(err)
			}
		}
	}

	os.Exit(0)
}

func listChanges(changes string) {
	fmt.Print(changes)
	os.Exit(0)
}

func main() {
	// Check that git is installed on this system
	checkGitCmd := exec.Command("git", "--version")
	if _, err := checkGitCmd.Output(); err != nil {
		handleError(errors.New("It appears that git cant be found in %PATH%"))
	}

	// Check if a .git repository exists
	if _, err := os.Stat("./.git"); err != nil {
		handleError(errors.New("No git repository could be found in current working directory"))
	}

	defaultLastCommitHash, _ := getGitHead()

	// Define command line flags
	// @link https://gobyexample.com/command-line-flags
	firstCommitPtr := flag.String("first", "", "The git commit that we are to begin at.")
	lastCommitPtr := flag.String("last", defaultLastCommitHash, "The git commit that we are to end at.")
	beVerbosePtr := flag.Bool("v", false, "Toggle verbose output.")
	outputList := flag.Bool("list", false, "List files rather than write to zip.")

	flag.Parse()

	// Validate Flags
	if validateFlags(firstCommitPtr, lastCommitPtr) == false {
		os.Exit(1)
	}

	// Check the current HEAD is equal to last
	if ! strings.HasPrefix(defaultLastCommitHash, *lastCommitPtr) {
		handleError(errors.New("You need to checkout commit [" + *lastCommitPtr + "] before continuing."))
	}

	// @todo validate that first commit hash and last commit hash exist
	// @todo validate that first commit hash is chronologically before last commit hash

	gitCommand := fmt.Sprintf("git diff-tree -r --no-commit-id --name-only --diff-filter=ACMRT %s %s", strings.TrimSpace(*firstCommitPtr), strings.TrimSpace(*lastCommitPtr))

	if *beVerbosePtr == true {
		fmt.Println("Executing: [" + gitCommand + "]")
	}

	// strings.Fields splits the string into a slice by spaces
	// https://golang.org/pkg/strings/#Fields
	command := strings.Fields(gitCommand)
	if len(command) < 2 {
		handleError(errors.New("unexpected error with gitCommand"))
	}
	gitDiffTreeCommand := exec.Command(command[0], command[1:]...)
	if gitDiffTree, err := gitDiffTreeCommand.Output(); err != nil {
		handleError(err)
	} else {
		if *outputList == true {
			listChanges(string(gitDiffTree))
		} else {
			zipChanges(string(gitDiffTree), *beVerbosePtr)
		}
	}
}
