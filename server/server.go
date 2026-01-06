package main

import (
	"context"
	"time"

	pb "todoapp/proto/pb"
)

// Server implements the TodoService
type server struct {
	pb.UnimplementedTodoServiceServer
	store *taskStore
}

// AddTask implements the AddTask RPC method
func (s *server) AddTask(ctx context.Context, req *pb.AddTaskRequest) (*pb.AddTaskResponse, error) {
	id := s.store.AddTask(req.GetDescription(), req.GetDueDate())
	return &pb.AddTaskResponse{Id: id}, nil
}

// ListTasks implements the ListTasks RPC method
func (s *server) ListTasks(ctx context.Context, req *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	tasks := s.store.GetAllTasks()
	return &pb.ListTasksResponse{
		Tasks: tasks,
	}, nil
}

// NEW: MarkTaskDone implements the MarkTaskDone RPC method
func (s *server) MarkTaskDone(ctx context.Context, req *pb.MarkTaskDoneRequest) (*pb.MarkTaskDoneResponse, error) {
	success := s.store.MarkTaskDone(req.GetId())
	return &pb.MarkTaskDoneResponse{
		Success: success,
	}, nil
}

// NEW: GetTasksByDate implements the GetTasksByDate RPC method
func (s *server) GetTasksByDate(ctx context.Context, req *pb.GetTasksByDateRequest) (*pb.GetTasksByDateResponse, error) {
	var date time.Time
	if req.GetDate() != nil {
		date = req.GetDate().AsTime()
	} else {
		date = time.Now() // Default to today
	}
	
	tasks := s.store.GetTasksByDate(date)
	return &pb.GetTasksByDateResponse{
		Tasks: tasks,
	}, nil
}