# Git Automation Tool

## Overview

This is a simple CLI tool built in Go to automate common Git tasks. The tool allows you to:

- Sync your current branch with the `main` branch.
- Clean up old local branches that have been merged into `main`.

This tool is designed to be run from within any Git project directory and can be executed globally on macOS.

## Prerequisites

- **Go** must be installed on your macOS machine. You can download and install it from [golang.org](https://golang.org/dl/).

## Installation

1. **Clone the Repository**

   ```bash
   git clone <repository-url>
   cd <repository-directory>
   ```

2. **Compile the Program**

   - Run the following command to compile the Go program into an executable file:

     ```bash
     go build -o git-automation
     ```

   - This will generate an executable file named `git-automation` in your current directory.

3. **Move the Executable to a Global Path**

   - To make the `git-automation` tool accessible from any directory, move it to a directory that is included in your system's `PATH`, such as `/usr/local/bin`:

     ```bash
     sudo mv git-automation /usr/local/bin/
     ```

   - This allows you to run the tool from any directory without needing to specify the path.

## Usage

1. **Navigate to Your Git Project Directory**

   - Go to the root of any Git project where you want to use the tool:

     ```bash
     cd /path/to/your/git/project
     ```

2. **Run the Tool**

   - Execute the tool by typing:

     ```bash
     git-automation
     ```

   - The program will prompt you to enter a command. You can choose from the following options:
     - `1`: Sync with the `main` branch
     - `2`: Clean up old local branches

- Note: Before running the script, ensure that your current branch is clean (i.e., no uncommitted changes or untracked files). This is important because there are still improvements to be made, and the script might otherwise cause issues if the working directory is not clean.

## Example Usage

```bash
cd /path/to/your/git/project
git-automation
```

**Example Commands**:

- To sync your current branch with the `main` branch:

  ```bash
  1
  ```

- To clean up old local branches:

  ```bash
  2
  ```

## Notes

- You must run this tool from within a Git project directory.
- Ensure that the `git-automation` executable is in your `PATH` for global access.
- The tool is designed to avoid merging `main` into protected branches such as `main`, `qa`, and `production`.

## Troubleshooting

- If you encounter any issues, ensure that the Go binary is compiled correctly and is executable:

```bash
chmod +x /usr/local/bin/git-automation
```

- This simple setup allows you to efficiently manage your Git workflow across different projects, with the convenience of a globally accessible command-line tool.

## License

- This project is licensed under the MIT License - see the [LICENSE](https://github.com/foreverLoveWisdom/git-automation-go/blob/main/LICENSE) file for details.
