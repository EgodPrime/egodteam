// DevTask represents a development task in the agile process
package models

import (
	"math/rand"
	"strconv"
	"time"
)

// TaskType represents the type of development task
type TaskType string

const (
	TaskDevelopment TaskType = "Development"
	TaskTesting     TaskType = "Testing"
	TaskDeployment  TaskType = "Deployment"
)

// TaskStatus represents the status of a development task
type TaskStatus string

const (
	TaskTodo       TaskStatus = "Todo"
	TaskInProgress TaskStatus = "InProgress"
	TaskDone       TaskStatus = "Done"
	TaskBlocked    TaskStatus = "Blocked"
)

// DevTask represents a development task in the agile process
type DevTask struct {
	ID           string        `json:"id"`
	StoryID      string        `json:"story_id"`
	Title        string        `json:"title"`
	Type         TaskType      `json:"type"`
	Estimate     time.Duration `json:"estimate"` // Estimated duration in hours
	Status       TaskStatus    `json:"status"`
	Assignee     string        `json:"assignee"`     // Developer ID
	Dependencies []string      `json:"dependencies"` // IDs of dependent tasks
}

// NewDevTask creates a new development task
func NewDevTask(storyID, title, assignee string) *DevTask {
	return &DevTask{
		ID:           generateIDDT(),
		StoryID:      storyID,
		Title:        title,
		Type:         TaskDevelopment,
		Estimate:     time.Hour * 4, // Default to 4 hours
		Status:       TaskTodo,
		Assignee:     assignee,
		Dependencies: []string{},
	}
}

// Update updates the task with new values
func (dt *DevTask) Update(title string, taskType TaskType, estimate time.Duration, assignee string) {
	dt.Title = title
	dt.Type = taskType
	dt.Estimate = estimate
	dt.Assignee = assignee
}

// SetStatus updates the status of the task
func (dt *DevTask) SetStatus(status TaskStatus) {
	dt.Status = status
}

// AddDependency adds a new dependency, preventing duplicates
func (dt *DevTask) AddDependency(dependencyID string) {
	for _, dep := range dt.Dependencies {
		if dep == dependencyID {
			return // Don't add if already exists
		}
	}
	dt.Dependencies = append(dt.Dependencies, dependencyID)
}

// RemoveDependency removes a dependency
func (dt *DevTask) RemoveDependency(dependencyID string) {
	for i, dep := range dt.Dependencies {
		if dep == dependencyID {
			dt.Dependencies = append(dt.Dependencies[:i], dt.Dependencies[i+1:]...)
			break
		}
	}
}

// generateIDDT generates a unique ID for a development task
func generateIDDT() string {
	return "task-" + time.Now().Format("20060102150405.000000") + "." + strconv.FormatInt(rand.Int63(), 10)
}
