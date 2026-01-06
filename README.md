Set-Content -Path README.md @"

# TodoApp with gRPC and Go

A Todo application built with gRPC, Protocol Buffers, and Go. This project demonstrates microservices communication using gRPC.

## Features

- Add tasks with descriptions and due dates
- List all tasks grouped by date (Today, Yesterday, etc.)
- Mark tasks as done
- View tasks by specific date
- Interactive CLI interface
- In-memory task storage

## Technology Stack

- **Go** (Golang) - Backend language
- **gRPC** - Remote Procedure Call framework
- **Protocol Buffers** - Data serialization
- **Git** - Version control

## Project Structure

TodoApp/
├── proto/ # Protocol Buffer definitions
│ ├── todo.proto # Service and message definitions
│ └── pb/ # Generated Go code
│ ├── todo.pb.go
│ └── todo_grpc.pb.go
├── server/ # gRPC server implementation
│ ├── main.go # Server entry point
│ ├── server.go # TodoService implementation
│ └── storage.go # In-memory task storage
├── client/ # gRPC client implementation
│ ├── main.go # Client entry point
│ └── client.go # Client helper functions
├── go.mod # Go module definition
├── go.sum # Dependency lock file
├── .gitignore # Git ignore rules
└── README.md # This file

text

## Prerequisites

- Go 1.21 or higher
- Protocol Buffer Compiler (protoc)
- Go plugins: protoc-gen-go, protoc-gen-go-grpc

## Installation

1. Clone the repository:
   \`\`\`bash
   git clone https://github.com/fuad-tafese-dev/todoapp.git
   cd todoapp
   \`\`\`

2. Install dependencies:
   \`\`\`bash
   go mod tidy
   \`\`\`

3. Install protoc and plugins:
   \`\`\`bash

   # Install protoc (see official protobuf releases)

   # Install Go plugins

   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   \`\`\`

4. Generate protobuf code:
   \`\`\`bash
   protoc --proto_path=proto --go_out=proto/pb --go_opt=paths=source_relative --go-grpc_out=proto/pb --go-grpc_opt=paths=source_relative proto/todo.proto
   \`\`\`

## Running the Application

### Start the Server:

\`\`\`bash
go run ./server localhost:50051
\`\`\`

### Run the Client (in another terminal):

\`\`\`bash
go run ./client localhost:50051
\`\`\`

## Usage

The client provides an interactive menu:

1. Add Task - Add a new task with description and due date
2. List All Tasks - View all tasks grouped by date
3. View Today's Tasks - See tasks due today
4. View Yesterday's Tasks - See overdue tasks from yesterday
5. Mark Task as Done - Complete a task
6. Exit - Quit the application

## API Methods

- \`AddTask\` - Add a new task
- \`ListTasks\` - List all tasks
- \`MarkTaskDone\` - Mark a task as completed
- \`GetTasksByDate\` - Get tasks for a specific date

## Example

\`\`\`
===== Todo App Menu =====

1. Add Task
2. List All Tasks
3. View Today's Tasks
4. View Yesterday's Tasks
5. Mark Task as Done
6. Exit
   Choose option: 1
   Enter task description: Buy groceries
   Enter due date (days from today): 0
   Enter time (HH:MM): 18:00

✅ Task added with ID: 1 (Due: Today, Jan 06, 2026)
\`\`\`

## Future Enhancements

- Database persistence (PostgreSQL/MySQL)
- Web interface
- User authentication
- Task categories and tags
- Email notifications
- REST API gateway

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

## License

This project is open source and available under the MIT License.

## Author

Fuad Tafese

- GitHub: [@fuad-tafese-dev](https://github.com/fuad-tafese-dev)
  "@
