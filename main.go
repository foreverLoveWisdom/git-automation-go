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
	fmt.Println("1. Sync with the Main Branch")
	fmt.Println("2. Cleanup")
	fmt.Print("Enter your choice (1-2): ")

	var command int
	_, err := fmt.Scanln(&command)

	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	// Execute function based on user input
	switch command {
	case 1:
		SyncWithMain(originalBranch)
	case 2:
		cleanupBranches()
	default:
		fmt.Println("Invalid command")
	}
}

// Core functionality

func SyncWithMain(originalBranch string) {
	// Fetch all remotes
	output := runGitCommand("fetch", "--all")
	fmt.Println(output)

	// Whitelist of branches to update
	branchesToUpdate := []string{"main", "qa"}

	// Loop through the whitelist and update each branch
	for _, branch := range branchesToUpdate {
		temp := runGitCommand("checkout", branch)
		fmt.Println(temp)
		temp = runGitCommand("pull", "origin", branch)
		fmt.Println(temp)
	}

	// Check if we should merge into the current branch
	blacklistedBranches := []string{"main", "qa", "production"} // Add any other branches as needed
	if isBranchBlacklisted(originalBranch, blacklistedBranches) {
		fmt.Printf("Merge operation aborted: Cannot merge 'main' into '%s'.\n", originalBranch)
		return
	}

	// Merge main into the current branch
	temp3 := runGitCommand("checkout", originalBranch)
	fmt.Println(temp3)
	temp4 := runGitCommand("merge", "main")
	fmt.Println(temp4)
}

func isBranchBlacklisted(branch string, blacklistedBranches []string) bool {
	for _, blacklistedBranch := range blacklistedBranches {
		if branch == blacklistedBranch {
			return true
		}
	}
	return false
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

// Helper functions

// Get the current branch name
func getCurrentBranch() string {
	output := runGitCommand("rev-parse", "--abbrev-ref", "HEAD")
	return strings.TrimSpace(output)
}

// Delete merged local branches
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

// Confirm and delete old branches
func confirmAndDelete(branches []string) {
	fmt.Println("Branches to delete: ")
	for _, branch := range branches {
		fmt.Println(branch)
	}

	var input int
	fmt.Println("Do you want to delete these branches?")
	fmt.Println("1. Yes")
	fmt.Println("2. No")
	fmt.Print("Enter your choice (1 or 2): ")
	_, err := fmt.Scanln(&input)

	if err != nil || (input != 1 && input != 2) {
		fmt.Println("Invalid input. Please enter 1 for Yes or 2 for No.")
		return
	}

	if input == 1 {
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
