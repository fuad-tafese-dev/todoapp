package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	pb "todoapp/proto/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatalln("usage: client [IP_ADDR] (e.g., client localhost:50051)")
	}

	addr := args[0]
	
	// Set up connection to the server
	conn, err := grpc.Dial(addr, 
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create client
	c := pb.NewTodoServiceClient(conn)

	reader := bufio.NewReader(os.Stdin)
	
	for {
		fmt.Println("\n===== Todo App Menu =====")
		fmt.Println("1. Add Task")
		fmt.Println("2. List All Tasks")
		fmt.Println("3. View Today's Tasks")
		fmt.Println("4. View Yesterday's Tasks")
		fmt.Println("5. Mark Task as Done")
		fmt.Println("6. Exit")
		fmt.Print("Choose option: ")
		
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		
		switch input {
		case "1":
			addTaskInteractive(c, reader)
		case "2":
			listTasks(c)
		case "3":
			viewTasksByDate(c, 0) // Today
		case "4":
			viewTasksByDate(c, -1) // Yesterday
		case "5":
			markTaskDoneInteractive(c, reader)
		case "6":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}

func addTaskInteractive(c pb.TodoServiceClient, reader *bufio.Reader) {
	fmt.Print("Enter task description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)
	
	if description == "" {
		fmt.Println("❌ Task description cannot be empty!")
		return
	}
	
	fmt.Print("Enter due date (days from today, e.g., 0 for today, 1 for tomorrow): ")
	daysInput, _ := reader.ReadString('\n')
	daysInput = strings.TrimSpace(daysInput)
	
	days := 0 // Default to today
	if daysInput != "" {
		d, err := strconv.Atoi(daysInput)
		if err != nil {
			fmt.Printf("❌ Invalid number. Using today.\n")
		} else {
			days = d
		}
	}
	
	fmt.Print("Enter time (HH:MM, e.g., 14:30 for 2:30 PM, or press Enter for all day): ")
	timeInput, _ := reader.ReadString('\n')
	timeInput = strings.TrimSpace(timeInput)
	
	// Calculate due date
	dueDate := time.Now().AddDate(0, 0, days)
	
	// Parse time if provided
	if timeInput != "" {
		parsedTime, err := time.Parse("15:04", timeInput)
		if err == nil {
			// Combine date and time
			dueDate = time.Date(
				dueDate.Year(), dueDate.Month(), dueDate.Day(),
				parsedTime.Hour(), parsedTime.Minute(), 0, 0, dueDate.Location(),
			)
		} else {
			fmt.Printf("❌ Invalid time format. Using all day.\n")
		}
	}
	
	id := addTask(c, description, dueDate)
	
	fmt.Printf("✅ Task added with ID: %d (Due: %s)\n", id, formatDateTitle(dueDate))
}

// NEW: View tasks by date
func viewTasksByDate(c pb.TodoServiceClient, daysOffset int) {
	date := time.Now().AddDate(0, 0, daysOffset)
	
	var title string
	switch daysOffset {
	case 0:
		title = "Today's Tasks"
	case -1:
		title = "Yesterday's Tasks"
	default:
		title = fmt.Sprintf("Tasks for %s", formatDateTitle(date))
	}
	
	fmt.Printf("\n📅 %s\n", title)
	fmt.Println("===============")
	
	tasks := getTasksByDate(c, date)
	displayTasksForDate(tasks)
}

// NEW: Interactive mark task as done
func markTaskDoneInteractive(c pb.TodoServiceClient, reader *bufio.Reader) {
	// First show current tasks
	fmt.Println("\nCurrent tasks:")
	tasks := getTasksByDate(c, time.Now())
	if len(tasks) == 0 {
		tasks = getTasksByDate(c, time.Now().AddDate(0, 0, -1))
	}
	
	if len(tasks) == 0 {
		fmt.Println("No tasks found to mark as done.")
		return
	}
	
	// Show tasks
	for _, task := range tasks {
		status := "❌"
		if task.GetDone() {
			status = "✅"
		}
		fmt.Printf("%s [%d] %s\n", status, task.GetId(), task.GetDescription())
	}
	
	// Ask for task ID
	fmt.Print("\nEnter task ID to mark as done: ")
	idInput, _ := reader.ReadString('\n')
	idInput = strings.TrimSpace(idInput)
	
	taskID, err := strconv.ParseUint(idInput, 10, 64)
	if err != nil {
		fmt.Println("❌ Invalid task ID")
		return
	}
	
	// Mark as done
	success := markTaskDone(c, taskID)
	if success {
		fmt.Printf("✅ Task %d marked as done!\n", taskID)
	} else {
		fmt.Printf("❌ Could not mark task %d as done (task not found)\n", taskID)
	}
}