package main

import (
	"Go-Basics/note/note"
	"Go-Basics/note/todo"
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Saver interface {
	Save() error
}

type outputtable interface {
	Saver
	Display()
}

func main() {

	title, content := getNoteData()
	userNote, err := note.New(title, content)
	if err != nil {
		panic(err)
	}
	err = outputData(userNote)
	if err != nil {
		panic(err)
	}

	todoText := getTodoData()
	userTodo, err := todo.New(todoText)
	if err != nil {
		panic(err)
	}
	err = outputData(userTodo)
	if err != nil {
		panic(err)
	}

}

func getNoteData() (string, string) {
	title := getUserInput("Title: ")
	content := getUserInput("Content: ")
	return title, content
}

func getTodoData() string {
	return getUserInput("Todo text: ")
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

func saveData(data Saver) error {
	err := data.Save()
	if err != nil {
		fmt.Println("Saving todo failed")
		return err
	}
	fmt.Println("Saving todo succeeded!")
	return nil
}

func outputData(data outputtable) error {
	data.Display()
	return saveData(data)
}

func printSomething(value interface{}) {
	switch value.(type) {
	case int:
		fmt.Println("Integer: ", value)
	case float64:
		fmt.Println("Float: ", value)
	default:
		fmt.Println(value)
	}
	/*typedValue, ok := value.(int)
	if ok {
		fmt.Println("Integer: ", typedValue)
		return
	}*/
}

func add[T int | float64 | string](a, b T) T {
	return a + b
}
