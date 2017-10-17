package main

import "fmt"
import "os"
import "os/exec"
import "flag"

func validateFlags(firstCommitPtr, lastCommitPtr *string) bool {
    valid := true

    if (len(*firstCommitPtr) == 0) {
        fmt.Println("Please provide a git commit to begin at via the -first flag.")
        valid = false
    }

    if (len(*firstCommitPtr) > 40) {
        fmt.Println("The hash for -first is not valid.")
        valid = false
    }

    if (len(*lastCommitPtr) == 0) {
        fmt.Println("Please provide a git commit to end at via the -last flag.")
        valid = false
    }

    if (len(*lastCommitPtr) > 40) {
        fmt.Println("The hash for -last is not valid.")
        valid = false
    }

    return valid
}

func main() {
    // Check that git is installed on this system
    checkGitCmd := exec.Command("git", "--version")
    if _, err := checkGitCmd.Output(); err != nil {
        fmt.Println("It appears that git cant be found in %PATH%.")
        os.Exit(1)
    }

    // Check if a .git repository exists
    if _, err := os.Stat("./.git"); err != nil {
        fmt.Println("No git repository could be found in current working directory.")
        os.Exit(1)
    }

    // Define command line flags
    // @link https://gobyexample.com/command-line-flags
    firstCommitPtr := flag.String("first", "", "The git commit that we are to begin at.")
    lastCommitPtr  := flag.String("last", "", "The git commit that we are to end at.")
    beVerbosePtr   := flag.Bool("v", false, "Toggle verbose output.")

    flag.Parse()

    // Validate Flags
    if (validateFlags(firstCommitPtr, lastCommitPtr) == false) {
        os.Exit(1)
    }

    // Check the current HEAD is equal to last
    // git rev-parse HEAD returns full 40 char hash, if the length of firstCommitPtr
    // is less than 40, check that it is equal from left to right.

    fmt.Println("first:", *firstCommitPtr)
    fmt.Println("last:", *lastCommitPtr)
    fmt.Println("verbose:", *beVerbosePtr)
}
