package main

import (
	"fmt"
	"os"

	"github.com/waldoapp/waldo-go-lib"
)

func main() {
	args := os.Args[1:]

	var cd string

	if len(args) > 0 {
		cd = args[0]

		if err := os.Chdir(cd); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to change directory to %s: %s\n", cd, err)
			return
		}
	} else {
		var err error

		if cd, err = os.Getwd(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get current directory: %s\n", err)
			return
		}
	}

	gitInfo := waldo.InferGitInfo(0)

	fmt.Printf("Git information for %s:\n\n", cd)
	fmt.Printf("Access: %s\n", gitInfo.Access())
	fmt.Printf("Branch: %s\n", gitInfo.Branch())
	fmt.Printf("Commit: %s\n", gitInfo.Commit())
	fmt.Print("\n")
}
