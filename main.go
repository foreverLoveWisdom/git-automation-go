package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// GitCommander defines the interface for running Git commands.
type GitCommander interface {
	RunGitCommand(args ...string) string
}

// BranchHandler defines the interface for branch-related operations.
type BranchHandler interface {
	GetCurrentBranch() string
	DeleteMergedLocalBranches() []string
	SyncWithMain(originalBranch string)
	CleanupBranches()
}

// Core functionality implementation

type GitCommandExecutor struct{}

// RunGitCommand executes a Git command and returns its output.
func (g GitCommandExecutor) RunGitCommand(args ...string) string {
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

// BranchManager handles branch operations using a GitCommander.
type BranchManager struct {
	GitCmd GitCommander
}

func (b BranchManager) GetCurrentBranch() string {
	output := b.GitCmd.RunGitCommand("rev-parse", "--abbrev-ref", "HEAD")
	return strings.TrimSpace(output)
}

func (b BranchManager) SyncWithMain(originalBranch string) {
	// Fetch all remotes
	output := b.GitCmd.RunGitCommand("fetch", "--all")
	fmt.Println(output)

	// Whitelist of branches to update
	branchesToUpdate := []string{"main", "qa"}

	// Loop through the whitelist and update each branch
	for _, branch := range branchesToUpdate {
		temp := b.GitCmd.RunGitCommand("checkout", branch)
		fmt.Println(temp)
		temp = b.GitCmd.RunGitCommand("pull", "origin", branch)
		fmt.Println(temp)
	}

	// Check if we should merge into the current branch
	blacklistedBranches := []string{"main", "qa", "production"} // Add any other branches as needed
	if isBranchBlacklisted(originalBranch, blacklistedBranches) {
		fmt.Printf("Merge operation aborted: Cannot merge 'main' into '%s'.\n", originalBranch)
		return
	}

	// Merge main into the current branch
	temp := b.GitCmd.RunGitCommand("checkout", originalBranch)
	fmt.Println(temp)
	temp = b.GitCmd.RunGitCommand("merge", "main")
	fmt.Println(temp)
}

func (b BranchManager) DeleteMergedLocalBranches() []string {
	b.GitCmd.RunGitCommand("fetch", "--prune")

	// List all branches merged into main
	output := b.GitCmd.RunGitCommand("branch", "--merged", "main")
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

func (b BranchManager) CleanupBranches() {
	branchesToDelete := b.DeleteMergedLocalBranches()
	// Return early if there are no branches to delete
	if len(branchesToDelete) == 0 {
		fmt.Println("No branches to delete. All merged branches are already cleaned up!")
		return
	}
	confirmAndDelete(branchesToDelete)
}

// Confirm and delete old branches.
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

func runGitCommand(args ...string) string {
	cmd := exec.Command("git", args...)
	fmt.Println("The command returned: ", cmd.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
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

// Helper function to check if the branch is blacklisted.
func isBranchBlacklisted(branch string, blacklistedBranches []string) bool {
	for _, blacklistedBranch := range blacklistedBranches {
		if branch == blacklistedBranch {
			return true
		}
	}
	return false
}

func main() {
	// Initialize the GitCommandExecutor and BranchManager
	gitExecutor := GitCommandExecutor{}
	branchManager := BranchManager{GitCmd: gitExecutor}

	// Get the current branch
	originalBranch := branchManager.GetCurrentBranch()

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
		branchManager.SyncWithMain(originalBranch)
	case 2:
		branchManager.CleanupBranches()
	default:
		fmt.Println("Invalid command")
	}
}
