package models

import (
	"testing"
	"time"
)

func TestNewDevTask(t *testing.T) {
	dt := NewDevTask("story-1", "Implement login feature", "dev-1")

	// Check that the task has a valid ID
	if dt.ID == "" || len(dt.ID) < 5 {
		t.Errorf("Expected non-empty ID, got %s", dt.ID)
	}

	// Check that the story ID is set correctly
	if dt.StoryID != "story-1" {
		t.Errorf("Expected StoryID to be 'story-1', got %s", dt.StoryID)
	}

	// Check that the title is set correctly
	if dt.Title != "Implement login feature" {
		t.Errorf("Expected Title to be 'Implement login feature', got %s", dt.Title)
	}

	// Check that the type is set to TaskDevelopment (default)
	if dt.Type != TaskDevelopment {
		t.Errorf("Expected Type to be TaskDevelopment, got %s", dt.Type)
	}

	// Check that the estimate is 4 hours (default)
	expectedEstimate := time.Hour * 4
	if dt.Estimate != expectedEstimate {
		t.Errorf("Expected Estimate to be %v, got %v", expectedEstimate, dt.Estimate)
	}

	// Check that the status is TaskTodo (default)
	if dt.Status != TaskTodo {
		t.Errorf("Expected Status to be TaskTodo, got %s", dt.Status)
	}

	// Check that the assignee is set correctly
	if dt.Assignee != "dev-1" {
		t.Errorf("Expected Assignee to be 'dev-1', got %s", dt.Assignee)
	}

	// Check that dependencies is empty (default)
	if len(dt.Dependencies) > 0 {
		t.Errorf("Expected Dependencies to be empty, got %v", dt.Dependencies)
	}
}

func TestDevTaskUpdate(t *testing.T) {
	dt := NewDevTask("story-1", "Implement login feature", "dev-1")

	// Update the task with new values
	dt.Update("Fix authentication bug", TaskTesting, time.Hour*2, "dev-2")

	// Check that all fields were updated correctly
	if dt.Title != "Fix authentication bug" {
		t.Errorf("Expected Title to be 'Fix authentication bug', got %s", dt.Title)
	}

	if dt.Type != TaskTesting {
		t.Errorf("Expected Type to be TaskTesting, got %s", dt.Type)
	}

	expectedEstimate := time.Hour * 2
	if dt.Estimate != expectedEstimate {
		t.Errorf("Expected Estimate to be %v, got %v", expectedEstimate, dt.Estimate)
	}

	if dt.Assignee != "dev-2" {
		t.Errorf("Expected Assignee to be 'dev-2', got %s", dt.Assignee)
	}
}

func TestDevTaskSetStatus(t *testing.T) {
	dt := NewDevTask("story-1", "Implement login feature", "dev-1")

	// Set status to InProgress
	dt.SetStatus(TaskInProgress)

	// Check that the status was updated correctly
	if dt.Status != TaskInProgress {
		t.Errorf("Expected Status to be TaskInProgress, got %s", dt.Status)
	}
}

func TestDevTaskAddDependency(t *testing.T) {
	dt := NewDevTask("story-1", "Implement login feature", "dev-1")

	// Add a dependency
	dt.AddDependency("task-2")

	// Check that the dependency was added
	if len(dt.Dependencies) != 1 {
		t.Errorf("Expected Dependencies length to be 1, got %d", len(dt.Dependencies))
	}

	if dt.Dependencies[0] != "task-2" {
		t.Errorf("Expected first dependency to be 'task-2', got %s", dt.Dependencies[0])
	}

	// Add the same dependency again (should not duplicate)
	dt.AddDependency("task-2")

	// Check that there's still only one instance of the dependency
	if len(dt.Dependencies) != 1 {
		t.Errorf("Expected Dependencies length to be 1 after adding duplicate, got %d", len(dt.Dependencies))
	}
}

func TestDevTaskRemoveDependency(t *testing.T) {
	dt := NewDevTask("story-1", "Implement login feature", "dev-1")
	dt.AddDependency("task-2")
	dt.AddDependency("task-3")

	// Remove a dependency that exists
	dt.RemoveDependency("task-2")

	// Check that the dependency was removed and only task-3 remains
	if len(dt.Dependencies) != 1 {
		t.Errorf("Expected Dependencies length to be 1 after removal, got %d", len(dt.Dependencies))
	}

	if dt.Dependencies[0] != "task-3" {
		t.Errorf("Expected remaining dependency to be 'task-3', got %s", dt.Dependencies[0])
	}

	// Try to remove a dependency that doesn't exist
	dt.RemoveDependency("task-4")

	// Check that the dependencies list remains unchanged
	if len(dt.Dependencies) != 1 {
		t.Errorf("Expected Dependencies length to be 1 after trying to remove non-existent dependency, got %d", len(dt.Dependencies))
	}
}

func TestDevTaskgenerateIDDT(t *testing.T) {
	// Test that generateIDDT creates unique IDs
	id1 := generateIDDT()
	id2 := generateIDDT()

	// Check that the IDs are different
	if id1 == id2 {
		t.Errorf("Expected unique IDs, got %s and %s", id1, id2)
	}

	// Check that the ID starts with "task-"
	if !startsWith(id1, "task-") || !startsWith(id2, "task-") {
		t.Errorf("Expected IDs to start with 'task-', got %s and %s", id1, id2)
	}

	// Check that the ID has a timestamp format (YYYYMMDDHHMMSS)
	if len(id1) < 15 || len(id2) < 15 {
		t.Errorf("Expected IDs to have timestamp format, got %s and %s", id1, id2)
	}
}

// Helper function to check if a string starts with another string
func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
