package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

//SavedServers stores server cache - initalized on loadCache
var SavedServers map[int][2]string

//CurrentServerIndex will store the parsed server index values
var CurrentServerIndex map[int]string

//IndexDownloadLoc is where the active server's index is downloaded
var IndexDownloadLoc string = "/tmp/"

// MainDownloadsLoc is the location where downloads are stored
var MainDownloadsLoc string = "/home/marcus/Downloads"

func main() {
	clearScreen()
	// if config not found then run the first launch wizard
	if !loadCache() {
		firstRunWizard()
		browseNow()
		loadCache()
	}

	for true {
		userChoice := menuHandler()
		serverIndexMenuHandler(userChoice)
	}
}

// =========== User Interface ===========
func menuHandler() string {
	serverPublicKey := ""
	consoleInput := bufio.NewReader(os.Stdin)
ServerMenu:

	renderServerBrowser()

	fmt.Print("(Press A to Add server, D to Delete a server, Q to quit): ")
	userChoice, _ := consoleInput.ReadString('\n')
	userChoice = strings.ToUpper(removeNewline(userChoice))
	switch userChoice {
	case "Q":
		os.Exit(1)
	case "A":
		clearScreen()
		addServer()
		loadCache()
	case "P":
		//TODO
	case "N":
		//TODO
	case "D":
		deleteServerWizard()
	default:
		userInt, err := strconv.Atoi(userChoice)
		if err != nil {
			break
		}
		if userInt >= 1 && userInt <= len(SavedServers) {
			serverPublicKey = SavedServers[userInt-1][1]
			// download file from server
			clearScreen()
			fmt.Println("Downloading Server Index...")
			dmsggetWrapper(serverPublicKey, IndexDownloadLoc, "index", "index."+serverPublicKey)
			loadServerIndex(serverPublicKey)
			goto ExitLoop
		} else {
			break
		}
	}
	goto ServerMenu
ExitLoop:
	return serverPublicKey
}

func serverIndexMenuHandler(serverPublicKey string) {
	consoleInput := bufio.NewReader(os.Stdin)
ServerIndexMenu:

	renderServerIndexBrowser()

	fmt.Print("(R to Refresh Server Index, E to Exit Server File Browser, Q to quit): ")
	userChoice, _ := consoleInput.ReadString('\n')
	userChoice = strings.ToUpper(removeNewline(userChoice))
	switch userChoice {
	case "Q":
		os.Exit(1)
	case "e":
		goto ExitLoop
	case "P":
		//TODO
		goto ServerIndexMenu
	case "N":
		//TODO
		goto ServerIndexMenu
	case "R":
		//TODOq
		goto ServerIndexMenu

	default:
		userInt, err := strconv.Atoi(userChoice)
		if err != nil {
			break
		}
		if userInt >= 1 && userInt <= len(CurrentServerIndex) {
			filenameDownload := CurrentServerIndex[userInt-1]
			// download file
			clearScreen()
			downloadInfo := fmt.Sprintf("Downloading %s to %s/", filenameDownload, MainDownloadsLoc)
			fmt.Println(downloadInfo)
			dmsggetWrapper(serverPublicKey, MainDownloadsLoc, filenameDownload, "")
		} else {
			break
		}
		goto ServerIndexMenu
	}
ExitLoop:
}

func renderServerBrowser() {
	pageStatus := fmt.Sprintf("page (%d / %d)", 1, 20)
	divider := "----------------------"
	clearScreen()

	fmt.Println("DMSG HTTP SERVER LIST")
	fmt.Println(divider)

	for i := 0; i < len(SavedServers); i++ {
		listEntry := fmt.Sprintf("%d) %s", i+1, SavedServers[i][0])
		fmt.Println(listEntry)
	}

	fmt.Println(divider)
	fmt.Println(pageStatus)
	fmt.Println("<< P  |  N >>")
}

func renderServerIndexBrowser() {
	pageStatus := fmt.Sprintf("page (%d / %d)", 1, 20)
	divider := "----------------------"
	clearScreen()

	fmt.Println("SERVER DOWNLOAD INDEX")
	fmt.Println(divider)

	for i := 0; i < len(CurrentServerIndex); i++ {
		listEntry := fmt.Sprintf("%d) %s", i+1, CurrentServerIndex[i])
		fmt.Println(listEntry)
	}

	fmt.Println(divider)
	fmt.Println(pageStatus)
	fmt.Println("<< P  |  N >>")
}

func generateConfigAbsPath() string {

	homeDir, err := os.UserHomeDir()

	if err != nil {
		log.Fatal(err)
	}

	configPath := "/.config/dmsg-http-browser.config"

	return homeDir + configPath
}

func generateServerIndexAbsPath(serverPublicKey string) string {
	indexPath := "/tmp/index." + serverPublicKey

	return indexPath
}

func clearScreen() {
	//TODO find a more elegant way of accomplishing this
	fmt.Print("\033[H\033[2J")
}

func dmsggetWrapper(publicKey string, downloadLoc string, file string, alternateFileName string) bool {
	fetchString := fmt.Sprintf("dmsg://%s:80/%s", publicKey, file)
	returnValue := true
	dmsggetPath, err := exec.LookPath("dmsgget")
	if err != nil {
		fmt.Println(err)
	}

	dmsggetCmd := &exec.Cmd{
		Path: dmsggetPath,
		Args: []string{dmsggetPath, "-O", downloadLoc + alternateFileName, fetchString},
		//Stdout: os.Stdout,
		Stderr: os.Stdout,
	}
	if err := dmsggetCmd.Run(); err != nil {
		fmt.Println("There was an error fetching the file")
		// file exists?
		returnValue = false
	}
	return returnValue
}

// =========== File I/O ===========
func loadServerIndex(serverPublicKey string) bool {
	returnBool := true
	file, err := os.Open(generateServerIndexAbsPath(serverPublicKey))
	defer file.Close()
	defer func() {
		if err := recover(); err != nil {
		}
	}()

	if err != nil {
		panic(err.Error())
	}

	fileStats, err := file.Stat()

	if err != nil {
		panic(err.Error())
	}

	if fileStats.Size() == 0 {
		returnBool = false
	}

	parseServerIndex(&file)
	return returnBool
}

func clearCacheConfig() {
	configFile := generateConfigAbsPath()
	file, err := os.Create(configFile)

	if err != nil {
		fmt.Println(err)
		fmt.Println(file)
	}

}

func appendToConfig(friendlyName string, serverPublicKey string) {

	friendlyName = removeNewline(friendlyName)
	serverPublicKey = removeNewline(serverPublicKey)
	rawData := friendlyName + ";" + serverPublicKey + string('\n')

	dataToWrite := []byte(rawData)

	f, err := os.OpenFile(generateConfigAbsPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

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

func loadCache() bool {
	returnBool := true
	file, err := os.Open(generateConfigAbsPath())
	defer file.Close()
	defer func() {
		if err := recover(); err != nil {
		}
	}()

	if err != nil {
		panic(err.Error())
	}

	fileStats, err := file.Stat()

	if err != nil {
		panic(err.Error())
	}

	if fileStats.Size() == 0 {
		returnBool = false
	}

	// load up map values
	parseConfigFile(&file)
	return returnBool
}

func parseConfigFile(file **os.File) {
	savedServers := make(map[int][2]string)
	friendlyNameIndex := 0
	serverPubKeyIndex := 1
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			fmt.Println("Error parsing configuration file.")
		}
	}()

	fileScan := bufio.NewScanner(*file)

	i := 0
	for fileScan.Scan() {
		splitStringArray := [2]string{"", ""}
		tmpString := fileScan.Text()
		tmpSplitString := strings.Split(tmpString, ";")
		splitStringArray[0] = tmpSplitString[friendlyNameIndex]
		splitStringArray[1] = tmpSplitString[serverPubKeyIndex]
		savedServers[i] = splitStringArray
		i++
	}
	SavedServers = savedServers
}

func parseServerIndex(file **os.File) {
	currentServerIndex := make(map[int]string)

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			fmt.Println("Error parsing server index file.")
		}
	}()

	fileScan := bufio.NewScanner(*file)

	i := 0
	for fileScan.Scan() {
		currentServerIndex[i] = fileScan.Text()
		i++
	}
	CurrentServerIndex = currentServerIndex
}

// =========== Wizards ===========
func browseNow() {
	consoleInput := bufio.NewReader(os.Stdin)

Browse:
	fmt.Print("Would you like to browse this server now? (Y/N): ")
	userAnswer, _ := consoleInput.ReadString('\n')

	switch formattedInput := strings.ToUpper(removeNewline(userAnswer)); formattedInput {
	case "Y":
		print("YES!... Attempting to load server index")
		//load server index
	case "N":
		// continue to main menu
	default:
		goto Browse
	}
}
func firstRunWizard() {
	fmt.Println("It looks like this is your frist time running the dmsg-http CLI browser.")
	addServer()
}

func addServer() {
	keyLength := 66
	consoleInput := bufio.NewReader(os.Stdin)
	fmt.Print("Please enter the public key for the dmsg-http server you want to add: ")

PubKeyInput:
	publicKey, _ := consoleInput.ReadString('\n')
	publicKey = removeNewline(publicKey)

	if len(publicKey) == keyLength {
		fmt.Print("Add a friendly name to this public key (default: [public_key]): ")
		friendlyName, _ := consoleInput.ReadString('\n')
		friendlyName = removeNewline(friendlyName)
		friendlyName = removeSemiColon(friendlyName)
		if len(friendlyName) == 0 {
			appendToConfig(publicKey, publicKey)
		} else {
			appendToConfig(friendlyName, publicKey)
		}
		fmt.Println("Entry cached.")

	} else {
		errorInfo := fmt.Sprintf("Provided key has length of %d. Expected length of %d.", len(publicKey), keyLength)
		fmt.Println(errorInfo)
		fmt.Print("Invalid key length please enter public key again: ")
		goto PubKeyInput
	}

}

func deleteServerIndex(indexToDelete int) {
	// SavedServers[indexToDelete]
	clearCacheConfig()

	for index := 0; index < len(SavedServers); index++ {
		if index == indexToDelete-1 {
			continue
		} else {
			appendToConfig(SavedServers[index][0], SavedServers[index][1])
		}
	}

	loadCache()
}

func deleteServerWizard() {
DeletePrompt:
	fmt.Print("Which server do you want to delete? (Enter C to Cancel): ")
	consoleInputWhichServer := bufio.NewReader(os.Stdin)
	userDelete, _ := consoleInputWhichServer.ReadString('\n')
	userDelete = strings.ToUpper(removeNewline(userDelete))

	switch userDelete {
	case "C":
		goto ExitLoop
	default:
		userInt, err := strconv.Atoi(userDelete)
		if err != nil {
			break
		}
	ConfirmDelete:
		if userInt >= 1 && userInt <= len(SavedServers) {
			deleteConfirmPrompt := fmt.Sprintf("Are you sure you want to delete (Y/N)? { %s | %s }", SavedServers[userInt-1][0], SavedServers[userInt-1][1])
			fmt.Println(deleteConfirmPrompt)
			deleteConfirmInput := bufio.NewReader(os.Stdin)

			deleteConfirm, _ := deleteConfirmInput.ReadString('\n')
			deleteConfirm = strings.ToUpper(removeNewline(deleteConfirm))
			//deleteIndex, err := strconv.Atoi(deleteConfirm)

			switch deleteConfirm {
			case "Y":

				deleteServerIndex(userInt)
				goto ExitLoop
			case "N":
				goto ExitLoop
			default:
				goto ConfirmDelete
			}
		} else {
			break
		}
	}
	goto DeletePrompt
ExitLoop:
}

// =========== String formatting functions ===========

func removeNewline(userInput string) string {
	return strings.TrimRight(userInput, "\n")
}

func removeSemiColon(stringToScan string) string {
	semiColonCode := byte(59)
	spaceBarCode := byte(32)
	tmpByteString := []byte(stringToScan)
	for i, v := range tmpByteString {
		if v == semiColonCode {
			tmpByteString[i] = spaceBarCode
		}
	}
	return string(tmpByteString)
}