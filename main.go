package main

import "fmt"
import "os"
import "os/exec"
import (
	"flag"
	"strings"
	"errors"
)

func handleError(err error) bool {
	if err == nil {
		return false
	}
	fmt.Println(err.Error())
	os.Exit(1)
	return true
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

func checkGitIsHead(lastCommitPtr string) (bool, error) {
	isHeadCmd := exec.Command("git", "rev-parse", "HEAD")

	if isHeadCmdOutput, err := isHeadCmd.Output(); err != nil {
		return false, errors.New("It appears that git cant be found in %PATH%")
	} else {
		return strings.HasPrefix(string(isHeadCmdOutput), lastCommitPtr), nil
	}
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

	// Define command line flags
	// @link https://gobyexample.com/command-line-flags
	firstCommitPtr := flag.String("first", "", "The git commit that we are to begin at.")
	lastCommitPtr := flag.String("last", "", "The git commit that we are to end at.")
	beVerbosePtr := flag.Bool("v", false, "Toggle verbose output.")

	flag.Parse()

	// Validate Flags
	if validateFlags(firstCommitPtr, lastCommitPtr) == false {
		os.Exit(1)
	}

	// Check the current HEAD is equal to last
	if isHead, err := checkGitIsHead(*lastCommitPtr); handleError(err) == false && isHead == false {
		handleError(errors.New("You need to checkout commit ["+ *lastCommitPtr +"] before continuing."))
	}

	fmt.Println("first:", *firstCommitPtr)
	fmt.Println("last:", *lastCommitPtr)
	fmt.Println("verbose:", *beVerbosePtr)
}
