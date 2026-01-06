package main

import (
	"sync"
	"time"

	pb "todoapp/proto/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// In-memory task store
type taskStore struct {
	tasks  map[uint64]*pb.Task
	mutex  sync.RWMutex
	nextID uint64
}

func New() *taskStore {
	return &taskStore{
		tasks:  make(map[uint64]*pb.Task),
		nextID: 1,
	}
}

func (s *taskStore) AddTask(description string, dueDate *timestamppb.Timestamp) uint64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	id := s.nextID
	s.nextID++

	s.tasks[id] = &pb.Task{
		Id:          id,
		Description: description,
		Done:        false,
		DueDate:     dueDate,
	}

	return id
}

// NEW: Mark task as done
func (s *taskStore) MarkTaskDone(id uint64) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	task, exists := s.tasks[id]
	if !exists {
		return false
	}
	
	task.Done = true
	return true
}

// Get all tasks
func (s *taskStore) GetAllTasks() []*pb.Task {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	tasks := make([]*pb.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	
	return tasks
}

// NEW: Get tasks for a specific date (ignores time part, only compares year/month/day)
func (s *taskStore) GetTasksByDate(date time.Time) []*pb.Task {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	tasks := make([]*pb.Task, 0)
	
	// Convert input date to year/month/day for comparison
	targetYear, targetMonth, targetDay := date.Date()
	
	for _, task := range s.tasks {
		if task.DueDate == nil {
			continue
		}
		
		taskTime := task.DueDate.AsTime()
		taskYear, taskMonth, taskDay := taskTime.Date()
		
		// Compare only year, month, and day (ignore time)
		if taskYear == targetYear && taskMonth == targetMonth && taskDay == targetDay {
			tasks = append(tasks, task)
		}
	}
	
	return tasks
}

// Optional: Get a single task by ID
func (s *taskStore) GetTask(id uint64) (*pb.Task, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	task, exists := s.tasks[id]
	return task, exists
}