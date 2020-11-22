package main

import (
	"dmsggui"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	//default interval
	sleepInterval := time.Duration(30) * time.Second

	//parse program arguments
	programArguments := ""
	if len(os.Args) > 1 {
		programArguments = os.Args[1]
		if strings.Contains(programArguments, "-h") {
			dmsggui.ClearScreen()
			printUseage()
			os.Exit(0)
		}
	}
	intergerValue, err := strconv.Atoi(programArguments)
	if err != nil && programArguments != "" {
		fmt.Println("Error interpreting user input", err)
		printUseage()
		os.Exit(1)
	} else if intergerValue > 0 {
		sleepInterval = time.Duration(intergerValue) * time.Second
	}

	fmt.Println("Indexing with an interval of:", sleepInterval)
	// enter main file monitor loop
	for true {
		directory, err := filePathWalk(".")
		if err != nil {
			fmt.Println("An error occured while reading the directory.")
		}
		clearCurrentIndex() //todo add diff check before rewriting index?
		for entry := range directory {

			if directory[entry][0] != "index" {
				appendToIndex(directory[entry])
			}
		}
		time.Sleep(sleepInterval)
	}
}

func printUseage() {
	fmt.Println("Program usage:  indexer [index_interval_in_seconds - (Default=30s)]")
	fmt.Println("--------------------------------------------------------------------")
	fmt.Println("Program is meant to be installed as a service. This makes switching the working index directory a lot easier.")
	fmt.Println()
}

// filePathWalk will list all absolute file paths and their sizes
func filePathWalk(root string) ([][2]string, error) {
	var files [][2]string
	var appendData [2]string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			fileInfo, err := os.Stat(path)
			if err != nil {
				fmt.Println(err)
			}
			fileSize := fmt.Sprint(fileInfo.Size())
			appendData[0] = path
			appendData[1] = string(fileSize)

			files = append(files, appendData)
		}
		return nil
	})
	return files, err
}

func clearCurrentIndex() {
	configFile := "./index"
	file, err := os.Create(configFile)

	if err != nil {
		fmt.Println(err)
		fmt.Println(file)
	}

}

func appendToIndex(filename [2]string) {
	filename[0] = removeNewline(filename[0])
	rawData := fmt.Sprintf("%s;%s\n", filename[0], filename[1])

	dataToWrite := []byte(rawData)

	f, err := os.OpenFile("./index", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
	}

	if _, err := f.Write(dataToWrite); err != nil {
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func removeNewline(userInput string) string {
	return strings.TrimRight(userInput, "\n")
}