package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

const configFile string = "./config.json"

type Config struct {
	SrcPaths []string `json:"srcPaths"`
	DstPath  string   `json:"dstPath"`
	RepoPath string   `json:"repoPath"`
}

func main() {

	logFile, err := os.Create("./logs.txt")
	if err != nil {
		log.Println("Error opening log file")
	} else {
		defer logFile.Close()
		log.SetOutput(logFile)
	}

	log.Println("Starting DotSync...")

	configBytes, err := os.ReadFile(configFile)
	if err != nil {
		log.Println("Error opening config file: ", err)
	}

	var config Config
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON file: %s", err)
	}

	config.DstPath, err = removesTilde(config.DstPath)
	if err != nil {
		log.Printf("Error removing tilde: %s\n", err)
	}

	config.RepoPath, err = removesTilde(config.RepoPath)
	if err != nil {
		log.Printf("Error removing tilde: %s\n", err)
	}

	var wg sync.WaitGroup
	wg.Add(len(config.SrcPaths))

	for _, src := range config.SrcPaths {
		go func(src string) {
			defer wg.Done()

			src, err = removesTilde(src)
			if err != nil {
				log.Println(err)
			}

			err = copy(src, config.DstPath)
			if err != nil {
				log.Printf("Error copying %s to %s: %s\n", src, config.DstPath, err)
			}
		}(src)
	}

	wg.Wait()

	timeStamp := time.Now().Format("2006-01-02 15:04:05")

	err = runCommandAt(config.RepoPath, "git", "add", "./")
	if err != nil {
		log.Println("Error adding files to git: ", err)
	}

	err = runCommandAt(config.RepoPath, "git", "commit", "-m", fmt.Sprintf("Auto-commit em %s: Backup por DotSync", timeStamp))
	if err != nil {
		log.Println("Error commiting files to git: ", err)
	}

	// err = runCommandAt(config.RepoPath, "git", "push", "origin", "main")
	// if err != nil {
	// 	log.Println("Erro ao subir o commit para o github: ", err)
	// }

	log.Println("All righty! Your configs are safe now.")
}

func copy(src, dst string) error {
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	base := filepath.Base(src)
	dstDir := filepath.Join(dst, base)

	if !sourceInfo.IsDir() {
		return copyFile(src, dstDir)
	}

	if err := os.MkdirAll(dstDir, 0744); err != nil {
		return err
	}

	err = filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dstDir, relPath)

		if d.IsDir() {

			return os.MkdirAll(dstPath, 0744)

		} else {
			return copyFile(path, dstPath)

		}
	})

	return err
}

func copyFile(src, dst string) error {

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err

}

func runCommandAt(path string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func removesTilde(path string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return path, err
	}

	if path[0] != '~' {
		return path, nil
	}

	if path == "~" {
		return home, nil
	}

	return filepath.Join(home, path[2:]), nil
}