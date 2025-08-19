package main

import (
	"encoding/json"
	"fmt"
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

func main () {
	configBytes, err := os.ReadFile(configFile)
	if err != nil {
		log.Println("Erro ao abrir o arquivo de configuração: ", err)
	}

	var config Config
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		log.Fatalf("Erro ao extrair configuração do arquivo JSON: %s", err)
	}

	config.DstPath, err = removesTilde(config.DstPath)
	if err != nil {
		log.Fatalf("Erro ao extrair tio: %s", err)
	}

	config.RepoPath, err = removesTilde(config.RepoPath)
	if err != nil {
		log.Fatalf("Erro ao extrair tio: %s", err)
	}

	var wg sync.WaitGroup
	wg.Add(len(config.SrcPaths))

	for _,srcPath := range config.SrcPaths{
		go copy(srcPath, config.DstPath, wg)
	}

	wg.Wait()

	err = runCommandAt(config.RepoPath, "git", "add", "./")
	if err != nil {
		log.Println("Erro ao adicionar arquivos: ", err)
	}

	err = runCommandAt(config.RepoPath, "git", "commit", "-m", fmt.Sprintf("Auto-commit em %s: Backup por DotSync", time.Now()))
	if err != nil {
		log.Println("Erro ao fazer commit: ", err)
	}

	err = runCommandAt(config.RepoPath, "git", "push", "origin", "main")
	if err != nil {
		log.Println("Erro ao subir o commit para o github: ", err)
	}
}
		func coy(src, det string, wg sync.WaitGroup){
			srcPath, err = removesTilde(srcPath)
			if err != nil {
				log.Fatalf("Erro ao extrair tio: %s", err)
			}

			fs := os.DirFS(srcPath)

			err := os.CopyFS(config.DstPath, fs)
			if err != nil {
				log.Printf("Erro ao copiar diretório: %s - %s", srcPath, err)
			}	

			wg.Done()
		}()

func runCommandAt (path string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func removesTilde(path string) (string,error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return path, err
	}

	if(path[0] != '~'){
		return path, nil
	}

	if path == "~" {
		return home, nil
	}

	return filepath.Join(home, path[2:]), nil
}
