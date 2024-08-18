package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func GetTodoFileName() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	filename := filepath.Join(home, "./.todos.json")
	return filename
}

func GetTodoFile() string {
	filename := GetTodoFileName()

	file, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(filename)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
		} else {
			log.Fatal("Unknown error occurred!")
		}
	}

	return file.Name()
}

func GetItemsFromFile() []Item {
	var items []Item
	filename := GetTodoFileName()
	bytes, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(bytes, &items)
	if err != nil {
		log.Fatal(err)
	}
	return items
}

func SeedFile() {
	file := GetTodoFileName()

	items := []Item{
		Item{"Learn Go"},
		Item{"Play with Ollama LLM"},
		Item{"Github APIs"},
		Item{"Dive into frontend"},
	}
	bytes, _ := json.Marshal(items)
	os.WriteFile(file, bytes, 0666)
}
