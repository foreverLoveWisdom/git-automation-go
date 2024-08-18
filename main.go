package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	// Get the current branch
	originalBranch := getCurrentBranch()

	fmt.Println("Enter a command: fetch, update, merge, or cleanup")
	var command string
	_, err := fmt.Scanln(&command)

	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	// Execute function based on user input
	switch command {
	case "fetch":
		fetchAll()
	case "update":
		updateBranches()
	case "merge":
		mergeMainIntoCurrent(originalBranch)
	case "cleanup":
		cleanupBranches()
	default:
		fmt.Println("Invalid command")
	}
}

// Fetch all remotes
func fetchAll() {
	output := runGitCommand("fetch", "--all")
	fmt.Println(output)
}

// Update main and QA branches
func updateBranches() {
	temp1 := runGitCommand("checkout", "main")
	fmt.Println(temp1)
	temp2 := runGitCommand("pull", "origin", "main")
	fmt.Println(temp2)

	temp3 := runGitCommand("checkout", "qa")
	fmt.Println(temp3)
	temp4 := runGitCommand("pull", "origin", "qa")
	fmt.Println(temp4)
}

// Merge main into the original branch
func mergeMainIntoCurrent(originalBranch string) {
	temp1 := runGitCommand("checkout", originalBranch)
	fmt.Println(temp1)
	temp2 := runGitCommand("merge", "main")
	fmt.Println(temp2)
}

// Cleanup old merged branches
func cleanupBranches() {
	branches := listLocalBranches()
	oldBranches := filterOldBranches(branches)
	confirmAndDelete(oldBranches)
}

// Utility functions

// Get the current branch name
func getCurrentBranch() string {
	output := runGitCommand("rev-parse", "--abbrev-ref", "HEAD")
	return strings.TrimSpace(output)
}

// List local branches
func listLocalBranches() []string {
	output := runGitCommand("branch")
	branches := strings.Split(output, "\n")
	return branches
}

// Filter out old branches and already merged
func filterOldBranches(branches []string) []string {
	var filteredBranches []string
	for _, branch := range branches {
		if branch != "main" && branch != "qa" && !strings.Contains(branch, "*") {
			filteredBranches = append(filteredBranches, branch)
		}
	}
	return filteredBranches
}

// Confirm and delete old branches
func confirmAndDelete(branches []string) {
	fmt.Println("Branches to delete: ")
	for _, branch := range branches {
		fmt.Println(branch)
	}

	var input string
	fmt.Print("Do you wanna delete these branches? (y/n): ")
	_, err := fmt.Scanln(&input)

	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	if input == "y" {
		for _, branch := range branches {
			runGitCommand("branch", "-D", branch)
		}
		fmt.Println("Branches deleted")
	} else {
		fmt.Println("Branches not deleted")
	}
}

// Run a Git command and return its output
func runGitCommand(args ...string) string {
	cmd := exec.Command("git", args...)
	fmt.Println("The command returned: ", cmd.String())
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	return string(output)
}
