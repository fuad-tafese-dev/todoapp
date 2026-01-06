package main

import (
	"context"
	"fmt"
	"time"

	pb "todoapp/proto/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// AddTask calls the AddTask RPC endpoint
func addTask(c pb.TodoServiceClient, description string, dueDate time.Time) uint64 {
	req := &pb.AddTaskRequest{
		Description: description,
		DueDate:     timestamppb.New(dueDate),
	}

	res, err := c.AddTask(context.Background(), req)
	if err != nil {
		panic(err)
	}

	return res.Id
}

// listTasks calls the ListTasks RPC endpoint
func listTasks(c pb.TodoServiceClient) {
	fmt.Println("\n📋 All Tasks:")
	fmt.Println("===============")
	
	req := &pb.ListTasksRequest{}
	res, err := c.ListTasks(context.Background(), req)
	if err != nil {
		fmt.Printf("Error getting tasks: %v\n", err)
		return
	}
	
	displayTasks(res.GetTasks())
}

// NEW: markTaskDone marks a task as done
func markTaskDone(c pb.TodoServiceClient, taskID uint64) bool {
	req := &pb.MarkTaskDoneRequest{
		Id: taskID,
	}
	
	res, err := c.MarkTaskDone(context.Background(), req)
	if err != nil {
		fmt.Printf("Error marking task as done: %v\n", err)
		return false
	}
	
	return res.GetSuccess()
}

// NEW: getTasksByDate gets tasks for a specific date
func getTasksByDate(c pb.TodoServiceClient, date time.Time) []*pb.Task {
	req := &pb.GetTasksByDateRequest{
		Date: timestamppb.New(date),
	}
	
	res, err := c.GetTasksByDate(context.Background(), req)
	if err != nil {
		fmt.Printf("Error getting tasks by date: %v\n", err)
		return nil
	}
	
	return res.GetTasks()
}

// NEW: displayTasks helper function
func displayTasks(tasks []*pb.Task) {
	// Check if no tasks
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		fmt.Println("===============")
		return
	}
	
	// Group tasks by date
	tasksByDate := make(map[string][]*pb.Task)
	
	for _, task := range tasks {
		var dateKey string
		if task.GetDueDate() != nil {
			dueDate := task.GetDueDate().AsTime()
			dateKey = dueDate.Format("2006-01-02")
		} else {
			dateKey = "No Date"
		}
		
		tasksByDate[dateKey] = append(tasksByDate[dateKey], task)
	}
	
	// Sort dates (simplified - in real app you'd sort properly)
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	
	// Display today's tasks first
	if todayTasks, ok := tasksByDate[today]; ok {
		fmt.Printf("\n📅 TODAY (%s)\n", formatDateTitle(time.Now()))
		displayTasksForDate(todayTasks)
		delete(tasksByDate, today)
	}
	
	// Display yesterday's tasks
	if yesterdayTasks, ok := tasksByDate[yesterday]; ok {
		fmt.Printf("\n📅 YESTERDAY (%s)\n", formatDateTitle(time.Now().AddDate(0, 0, -1)))
		displayTasksForDate(yesterdayTasks)
		delete(tasksByDate, yesterday)
	}
	
	// Display other dates
	for dateStr, dateTasks := range tasksByDate {
		date, _ := time.Parse("2006-01-02", dateStr)
		fmt.Printf("\n📅 %s\n", formatDateTitle(date))
		displayTasksForDate(dateTasks)
	}
	
	fmt.Printf("\nTotal tasks: %d\n", len(tasks))
	fmt.Println("===============")
}

// NEW: Helper to display tasks for a specific date
func displayTasksForDate(tasks []*pb.Task) {
	doneCount := 0
	
	for _, task := range tasks {
		// Status icon
		status := "❌"
		if task.GetDone() {
			status = "✅"
			doneCount++
		}
		
		// Time if available
		var timeStr string
		if task.GetDueDate() != nil {
			dueDate := task.GetDueDate().AsTime()
			timeStr = dueDate.Format("15:04")
		}
		
		fmt.Printf("   %s [%d] %s", status, task.GetId(), task.GetDescription())
		if timeStr != "" {
			fmt.Printf(" (at %s)", timeStr)
		}
		fmt.Println()
	}
	
	// Show summary
	total := len(tasks)
	if total > 0 {
		percentage := (doneCount * 100) / total
		fmt.Printf("   📊 %d/%d done (%d%%)\n", doneCount, total, percentage)
	}
}

// NEW: Helper to format date title
func formatDateTitle(date time.Time) string {
	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)
	tomorrow := today.AddDate(0, 0, 1)
	
	dateStr := date.Format("Jan 02, 2006")
	
	switch {
	case date.Year() == today.Year() && date.Month() == today.Month() && date.Day() == today.Day():
		return "Today, " + dateStr
	case date.Year() == yesterday.Year() && date.Month() == yesterday.Month() && date.Day() == yesterday.Day():
		return "Yesterday, " + dateStr
	case date.Year() == tomorrow.Year() && date.Month() == tomorrow.Month() && date.Day() == tomorrow.Day():
		return "Tomorrow, " + dateStr
	default:
		return date.Format("Monday, Jan 02, 2006")
	}
}