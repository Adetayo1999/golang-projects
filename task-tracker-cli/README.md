# Task Tracker CLI

Cli App for [Roadmap Tast Tracker](https://roadmap.sh/projects/task-tracker)

This basic cli app can be used to add, create, update and view todos.

## How To Run The Program

```bash

# To build the todo app
go build ./cmd/cli -o task-tracker-cli

# To start the todo app
./task-tracker-cli

# Add todo
add "Buy groceries"

# Update Todo
update 1 "Buy groceries updated"

# Delete Todo
delete 1

# List all todos
list # list the entire todos
list done # list all todos with the done status
list in-progress # list all todos with the in progress status
list todo  #list all todos with the todo status

# Update todo status
mark-in-progress 1 # update the todo with ID 1 to in-progress
mark-done 1 # update the todo with ID 1 to done

# Exit program
exit

```
