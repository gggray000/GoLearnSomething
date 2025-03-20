package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Todo struct {
	Text string `json:"text"`
}

func New(text string) (Todo, error) {

	if text == "" {
		return Todo{}, errors.New("Invalid Input ")
	}

	return Todo{
		Text: text,
	}, nil
}

func (todo Todo) Display() {
	fmt.Printf("Todo:%v\n", todo.Text)
}

func (todo Todo) Save() error {
	fileName := "todo.json"
	jsonTodo, err := json.Marshal(todo)
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, jsonTodo, 0644)
}
