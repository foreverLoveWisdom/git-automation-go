package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	// Get the current branch
	originalBranch := getCurrentBranch()

	fmt.Println("Select a command:")
	fmt.Println("1. Fetch")
	fmt.Println("2. Update")
	fmt.Println("3. Merge")
	fmt.Println("4. Cleanup")
	fmt.Print("Enter your choice (1-4): ")

	var command int
	_, err := fmt.Scanln(&command)

	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	// Execute function based on user input
	switch command {
	case 1:
		fetchAll()
	case 2:
		updateBranches()
	case 3:
		mergeMainIntoCurrent(originalBranch)
	case 4:
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
	branchesToDelete := deleteMergedLocalBranches()
	// Return early if there are no branches to delete
	if len(branchesToDelete) == 0 {
		fmt.Println("No branches to delete. All merged branches are already cleaned up!")
		return
	}
	confirmAndDelete(branchesToDelete)
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

	var cleanedBranches []string
	for _, branch := range branches {
		branch = strings.TrimSpace(branch)
		if strings.HasPrefix(branch, "*") {
			branch = strings.TrimPrefix(branch, "* ")
		}
		if branch != "" {
			cleanedBranches = append(cleanedBranches, branch)
		}
	}
	return cleanedBranches
}

// func isFullyMerged(branch string) bool {
// 	output := runGitCommand("rev-list", "--no-walk", branch, "^main")
// 	return output == ""
// }

func deleteMergedLocalBranches() []string {
	runGitCommand("fetch", "--prune")

	// List all branches merged into main
	output := runGitCommand("branch", "--merged", "main")
	branches := strings.Split(output, "\n")

	var branchesToDelete []string
	for _, branch := range branches {
		branch = strings.TrimSpace(branch)
		// Exclude the main branch and any other protected branches
		if branch != "" && branch != "main" && branch != "qa" && branch != "* main" {
			branchesToDelete = append(branchesToDelete, branch)
		}
	}

	return branchesToDelete
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
		exitErr, ok := err.(*exec.ExitError)
		if ok {
			if exitErr.ExitCode() == 128 {
				fmt.Printf("Error running command: %v - likely a branch issue\n", err)
			} else {
				fmt.Printf("Error running command: %v\n", err)
			}
		} else {
			fmt.Printf("Error: %v\n", err)
		}
		return ""
	}

	return string(output)
}
