package file

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func FindFile(fPath, fName string) string {
	f, err := findFile(fPath, fName)
	if err != nil {
		log.Println("error on find file")
		return ""
	}

	return f
}

func findFile(fPath, fName string) (string, error) {
	var f string

	err := filepath.WalkDir(fPath, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			log.Println("error on walk dir callback", err)
			return err
		}
		if !info.IsDir() && info.Name() == fName {
			f = path
		}
		return nil
	})

	if err != nil {
		log.Println("error on walk dir", err)
		return "", err
	}

	return f, nil
}

// ExecReplace
func dotenv(filepath string) {
	// read file into memory
	f, err := os.Open(filepath)
	if err != nil {
		log.Println("error open file:", err)
	}
	defer f.Close()

	// var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		env := scanner.Text()
		if env == "" {
			continue
		}
		key := strings.Split(env, "=")[0]
		value := strings.Split(env, "=")[1]
		os.Setenv(key, value)
	}
}

func LoadEnv(fRootPath, fName string) {
	f, err := findFile(fRootPath, fName)
	if err != nil {
		log.Fatal("file not found", err)
	}
	dotenv(f)
}

type FileContext interface {
	Action(s string) string
}

// ExecReplace
func Exec(ctx FileContext, filepath string) {
	// read file into memory
	f, err := os.Open(filepath)
	if err != nil {
		log.Println("error on file open:", err)
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		line = ctx.Action(line) // Performs Action based on passed Context type
		if len(line) == 0 {     // Skips empty lines write
			continue
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		log.Println("erron on scanner:", err)
	}

	file, err := os.Create(filepath)
	if err != nil {
		log.Println("error on create file:", err)
	}
	defer file.Close()

	// write modified contents back to file
	for _, line := range lines {
		fmt.Fprintln(file, line)
	}
}
