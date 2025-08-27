# DotSync: A Concurrent Dotfile Backup Utility

DotSync is a simple and efficient command-line tool for backing up configuration files (dotfiles) to a Git repository. This project was built to solve a personal need—automating configuration backups—while serving as a practical exercise to solidify foundational Go programming concepts.

## Project Genesis & Purpose

As a developer, my configuration files are crucial. I wanted an automated way to back them up to a Git repository. DotSync was the perfect project to build this solution while applying key Go concepts I was learning, such as concurrency, file I/O, and external process execution. It's a tool I use personally, built with the technologies I wanted to master.

## Technical Deep Dive: Go Concepts in Practice

This project was a deliberate exercise in using Go's standard library to build a robust tool. The key concepts implemented are:

* **Concurrency (`Goroutines` & `WaitGroups`):** Leverages goroutines to perform parallel file I/O, copying multiple sources concurrently for significant speed improvements. A `sync.WaitGroup` ensures all copy operations are complete before the final Git commit.

* **File System Operations (`os`, `path/filepath`):** Uses the `os` and `path/filepath` packages for platform-agnostic file reading, writing, and recursive directory traversal.

* **JSON Configuration (`encoding/json`):** Employs a `config.json` file for all settings, parsed into a custom Go `struct` for clean, type-safe configuration management.

* **External Process Execution (`os/exec`):** Integrates directly with Git by using the `os/exec` package to run `git add` and `git commit` commands, automating the versioning process.

* **Robust Error Handling & Logging:** Implements Go's standard error handling patterns (`if err != nil`) and directs all operational output to a `logs.txt` file for easy debugging.

## Getting Started

### Prerequisites

* Go (version 1.18 or higher)
* Git installed and configured on your machine.

### 1. Configuration

Create a `config.json` file in the same directory as the program.

**Example `config.json`:**
```json
{
  "srcPaths": [
    "~/.zshrc",
    "~/.config/nvim"
  ],
  "dstPath": "~/backups/dotfiles",
  "repoPath": "~/backups/dotfiles"
}
```

* `srcPaths`: A list of files or directories to back up.
* `dstPath`: The folder inside your Git repository where items will be copied.
* `repoPath`: The path to the root of the local Git repository.

### 2. Running the Program

Navigate to the project directory and run:

```bash
go run main.go
```

The program will back up your files and commit the changes with a timestamp. Or you could set up a cron job to automate it even further, this is what I'm planning on doing.

## Future Improvements
* Add `--config` flag to add paths.
* Add a `--push` flag to push the commit to a remote.
* Implement smarter syncing to only copy changed files.
