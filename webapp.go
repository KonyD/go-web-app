package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
)

var pl = fmt.Println

type ToDoList struct {
	ToDoCount int
	ToDos []string
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func write(writer http.ResponseWriter, msg string) {
	_, err := writer.Write([]byte(msg))
	checkError(err)
}

func getStrings(fileName string) []string  {
	var lines []string
	file, err := os.Open("todos.txt")
	if os.IsNotExist(err) {
		return nil
	}
	checkError(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	checkError(scanner.Err())
	return lines
}

func englishHandler(writer http.ResponseWriter, request *http.Request) {
	write(writer, "Hello internet")
}

func spanishHandler(writer http.ResponseWriter, request *http.Request) {
	write(writer, "Hola internet")
}

func turkishHandler(writer http.ResponseWriter, request *http.Request) {
	write(writer, "Merhaba internet")
}

func interactHandler(writer http.ResponseWriter, request *http.Request) {
	todoVals := getStrings("todos.txt")
	fmt.Printf("%#v\n", todoVals)
	tmpl, err := template.ParseFiles("view.html")
	checkError(err)
	todos := ToDoList {
		ToDoCount: len(todoVals),
		ToDos: todoVals,
	}
	err = tmpl.Execute(writer, todos)
}

func newHandler(writer http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFiles("new.html")
	checkError(err)
	err = tmpl.Execute(writer, nil)
}

func createHandler(writer http.ResponseWriter, request *http.Request) {
	todo := request.FormValue("todo")
	options := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	file, err := os.OpenFile("todos.txt", options, os.FileMode(0600))
	checkError(err)
	_, err = fmt.Fprintln(file, todo)
	checkError(err)
	err = file.Close()
	checkError(err)
	http.Redirect(writer, request, "/interact", http.StatusFound)
}

func main() {
	http.HandleFunc("/en", englishHandler)
	http.HandleFunc("/es", spanishHandler)
	http.HandleFunc("/tr", turkishHandler)
	http.HandleFunc("/interact", interactHandler)
	http.HandleFunc("/new", newHandler)
	http.HandleFunc("/create", createHandler)
	pl("Listening on http://localhost:8080")
	err := http.ListenAndServe("localhost:8080", nil)
	checkError(err)	
}