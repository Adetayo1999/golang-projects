package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/fatih/color"
)

type TaskStatus string

type Todo struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type Todos map[string]Todo

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in-progress"
	StatusDone       TaskStatus = "done"
)

const DataStoreName = "data.json";

var yellow = color.New(color.FgYellow)
var green = color.New(color.FgGreen)
var red = color.New(color.FgRed)

func main() {

	var cmdHelpers = []string{
		 `1. add command; usage "add "<todo>" " `,
		 `2. update command; usage "update <todo-id>  "<todo>" "`,
		 `3. delete command; usage "delete <todo-id>"`,
		 `4. list command; usage "list <status>"; status can be (done, todo, in-progress) and it is optional; `,
		 `5. mark in progress command; usage "mark-in-progress <todo-id>"`,
		 `6. mark done command; usage "mark-done <todo-id>"`,
	}

	var todos = Todos{}

	file, err := os.OpenFile("data.json", os.O_CREATE|os.O_RDWR, 0644)

	if err != nil {
		panic(err);
	} 

	 defer file.Close()

	err = json.NewDecoder(file).Decode(&todos);
	
	if err != nil {
		if err.Error() != "EOF" {
		panic(err);
		}
	}
		

	for {

		fmt.Print("task-cli: ")
		reader := bufio.NewReader(os.Stdin)

		cmdString, err := reader.ReadString('\n')

		if err != nil {
			red.Println(err.Error())
			continue
		}

		cmdStringArr := splitCommand(cmdString)

		switch {

		case cmdStringArr[0] == "add" && len(cmdStringArr) == 2:
			unquotedDesc, err := strconv.Unquote(cmdStringArr[1])
			if err != nil {
				red.Println(err.Error())
				continue
			}
			todos.Add(unquotedDesc, file)

		case cmdStringArr[0] == "update" && len(cmdStringArr) == 3:
			unquotedDesc, err := strconv.Unquote(cmdStringArr[2])
			if err != nil {
				red.Println(err.Error())
				continue
			}
			id, err := strconv.ParseInt(cmdStringArr[1], 10, 64)
			if err != nil {
				red.Println(err.Error())
				continue
			}
			err = todos.Update(int(id), unquotedDesc, file)
			if err != nil {
				red.Println(err.Error())
				continue
			}

		case cmdStringArr[0] == "delete" && len(cmdStringArr) == 2:
			id, err := strconv.ParseInt(cmdStringArr[1], 10, 64)
			if err != nil {
				red.Println(err.Error())
				continue
			}
			todos.Delete(int(id), file)

		case cmdStringArr[0] == "list" && len(cmdStringArr) == 1:
			todos.List(nil)

		case cmdStringArr[0] == "list" && len(cmdStringArr) == 2:

			var status = StatusTodo

			statusFromCmd := cmdStringArr[1]

			if statusFromCmd == string(StatusDone) || statusFromCmd == string(StatusInProgress) || statusFromCmd == string(StatusTodo) {
				status = TaskStatus(statusFromCmd)
			}

			todos.List(&status)

		case cmdStringArr[0] == "mark-in-progress" && len(cmdStringArr) == 2:
			id, err := strconv.ParseInt(cmdStringArr[1], 10, 64)
			if err != nil {
				red.Println(err.Error())
				continue
			}
			err = todos.UpdateTodoStatus(int(id), StatusInProgress, file)
			if err != nil {
				red.Println(err.Error())
				continue
			}

		case cmdStringArr[0] == "mark-done" && len(cmdStringArr) == 2:
			id, err := strconv.ParseInt(cmdStringArr[1], 10, 64)
			if err != nil {
				red.Println(err.Error())
				continue
			}
			err = todos.UpdateTodoStatus(int(id), StatusDone, file)
			if err != nil {
				red.Println(err.Error())
				continue
			}
		
		case cmdStringArr[0] == "exit" && len(cmdStringArr) == 1:
			red.Println("terminating task tracker app...")
			os.Exit(0);	

		case cmdStringArr[0] == "help" && len(cmdStringArr) == 1:
			fmt.Println("-------------------------------")
			for i := range cmdHelpers{
				green.Println(cmdHelpers[i])
			}
			fmt.Println("-------------------------------")
			
		default:
			red.Println("Invalid command")
			fmt.Println("-------------------------------")
			for i := range cmdHelpers{
				green.Println(cmdHelpers[i])
			}
			fmt.Println("-------------------------------")
		}

	}

}

func splitCommand(input string) []string {
	// Regular expression to match quoted text and words separately
	// It matches either quoted strings or individual words.
	re := regexp.MustCompile(`"([^"]+)"|(\S+)`)

	// Find all matches
	matches := re.FindAllString(input, -1)

	return matches
}

func (todos *Todos) Add(description string, file *os.File) error  {

	newTodo := Todo{
		ID:          len(*todos) + 1,
		Description: description,
		Status:      StatusTodo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	(*todos)[string(len(*todos)+1)] = newTodo

	_, err := file.Seek(0, 0)
	if err != nil {
		panic(err)
	}

	err = file.Truncate(0) 
	if err != nil {
		panic(err)
	}

	err = json.NewEncoder(file).Encode(*todos)
	if err != nil {
		panic(err)
	}


	green.Printf("%+v \n", *todos)

	return nil
}

func (todos *Todos) Update(id int, description string, file *os.File) error {

	todo, exists := (*todos)[string(id)]

	if !exists {
		return fmt.Errorf("no todo with ID: %d", id)
	}

	todo.Description = description
	todo.UpdatedAt = time.Now()

	(*todos)[string(id)] = todo

		_, err := file.Seek(0, 0)
	if err != nil {
		panic(err)
	}

	err = file.Truncate(0) 
	if err != nil {
		panic(err)
	}

	err = json.NewEncoder(file).Encode(*todos)
	if err != nil {
		panic(err)
	}


	green.Printf("%+v \n", *todos)
	return nil
}

func (todos *Todos) Delete(id int, file *os.File) {
	delete(*todos, string(id))

		_, err := file.Seek(0, 0)
	if err != nil {
		panic(err)
	}

	err = file.Truncate(0) 
	if err != nil {
		panic(err)
	}

	err = json.NewEncoder(file).Encode(*todos)
	if err != nil {
		panic(err)
	}


	red.Printf("Todo with ID:%d deleted \n", id)
}

func (todos *Todos) List(status *TaskStatus) {


	for _, value := range *todos {

		if status != nil && value.Status != *status {
			continue
		}

		stmt := fmt.Sprintf("ID: %d, Description: %v, Status: %v, CreatedAt: %s, UpdatedAt: %s", value.ID, value.Description, value.Status, value.CreatedAt.Format(time.RFC850), value.UpdatedAt.Format(time.RFC850))

		switch true {

		case value.Status == StatusDone:
			green.Println(stmt)
		case value.Status == StatusInProgress:
			yellow.Println(stmt)
		default:
			red.Println(stmt)

		}
	}

}

func (todos *Todos) UpdateTodoStatus(id int, status TaskStatus, file *os.File) error {
	todo, exists := (*todos)[string(id)]

	if !exists {
		return fmt.Errorf("no todo with ID: %d", id)
	}

	previousTodo := todo.Status

	todo.Status = status

	(*todos)[string(id)] = todo

		_, err := file.Seek(0, 0)
	if err != nil {
		panic(err)
	}

	err = file.Truncate(0) 
	if err != nil {
		panic(err)
	}

	err = json.NewEncoder(file).Encode(*todos)
	if err != nil {
		panic(err)
	}


	green.Printf("Updated todo with ID: %d from %s to %s \n", todo.ID, previousTodo, todo.Status)

	return nil
}
