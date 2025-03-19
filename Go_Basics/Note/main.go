package main

import (
	"Go-Basics/note/struct"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	title, content := getNoteData()

	userNote, err := note.New(title, content)

	if err != nil {
		panic(err)
	}

	userNote.Display()
	err = userNote.Save()

	if err != nil {
		fmt.Println("Saving note failed")
		return
	}

	fmt.Println("Saving note succeeded!")
}

func getNoteData() (string, string) {
	title := getUserInput("Title: ")
	content := getUserInput("Content: ")
	return title, content
}

func getUserInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}
	text = strings.TrimSuffix(text, "\n")
	text = strings.TrimSuffix(text, "\r")
	return text
}
