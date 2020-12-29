package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	clipboard "github.com/atotto/clipboard"
	htgotts "github.com/weiztd/htgo-tts"
	handlers "github.com/weiztd/htgo-tts/handlers"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

func main() {
	var langCode, tempFolder, input string

	tempFolder = os.Getenv("APPDATA") + "/clipboardTTS"

	for {
		fmt.Print("Language code: \n")
		_, err := fmt.Scanln(&input)
		if err != nil {
			log.Fatal(err)
		}
		if isLangCodeLegit(input) {
			langCode = input
			break
		}
	}

	err := speechFromClipboard(langCode, tempFolder)
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(tempFolder)
}

func isLangCodeLegit(input string) bool {
	tag, err := language.Parse(input)
	if err != nil {
		fmt.Println("invalid language code: " + input)
		return false
	}
	fmt.Printf("Language set to: %s (%s)\n",
		display.English.Tags().Name(tag),
		display.Self.Name(tag))
	return true
}

func speechFromClipboard(langCode, tempFolder string) error {
	quitChan := make(chan os.Signal)
	signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	temp := ""
	clipboard.WriteAll(temp)
	for {
		select {
		case <-time.After(200 * time.Millisecond):
			clipboardString, err := clipboard.ReadAll()
			if err != nil {
				return fmt.Errorf("fail to read from clipboard: %v", err)
			}
			if len(clipboardString) > 1 && clipboardString != temp {
				log.Println(clipboardString)
				speech := htgotts.Speech{Folder: tempFolder, Language: langCode, Handler: &handlers.MPlayer{}}
				speech.Speak(clipboardString)
				temp = clipboardString
			}
		case <-quitChan:
			return nil
		}
	}
}
