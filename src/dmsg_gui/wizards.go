package dmsg_gui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func DeleteServerWizard() {
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

				DeleteServerIndex(userInt)
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

func addServer() string {
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
			AppendToConfig(publicKey, publicKey)
		} else {
			AppendToConfig(friendlyName, publicKey)
		}
		fmt.Println("Entry cached.")

	} else {
		errorInfo := fmt.Sprintf("Provided key has length of %d. Expected length of %d.", len(publicKey), keyLength)
		fmt.Println(errorInfo)
		fmt.Print("Invalid key length please enter public key again: ")
		goto PubKeyInput
	}
	return publicKey
}

func FirstRunWizard() {
	fmt.Println("It looks like this is your frist time running the dmsg-http CLI browser.")
	serverPublicKey := addServer()
	browseNow(serverPublicKey)
	LoadCache()
}

func browseNow(serverPublicKey string) {
	consoleInput := bufio.NewReader(os.Stdin)

Browse:
	fmt.Print("Would you like to browse this server now? (Y/N): ")
	userAnswer, _ := consoleInput.ReadString('\n')

	switch formattedInput := strings.ToUpper(removeNewline(userAnswer)); formattedInput {
	case "Y":
		RefreshServerIndex(serverPublicKey, true)
		ServerIndexMenuHandler(serverPublicKey)
		//load server index
	case "N":
		// continue to main menu
	default:
		goto Browse
	}
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