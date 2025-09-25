package models

import (
	"testing"
	"time"
)

func TestNewSprint(t *testing.T) {
	startDate := time.Date(2025, 9, 25, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 10, 2, 0, 0, 0, 0, time.UTC)

	s := NewSprint("Implement user authentication", startDate, endDate)

	// Check that the sprint has a valid ID
	if s.ID == "" || len(s.ID) < 5 {
		t.Errorf("Expected non-empty ID, got %s", s.ID)
	}

	// Check that the goal is set correctly
	if s.Goal != "Implement user authentication" {
		t.Errorf("Expected Goal to be 'Implement user authentication', got %s", s.Goal)
	}

	// Check that start and end dates are set correctly
	if !s.StartDate.Equal(startDate) {
		t.Errorf("Expected StartDate to be %v, got %v", startDate, s.StartDate)
	}
	if !s.EndDate.Equal(endDate) {
		t.Errorf("Expected EndDate to be %v, got %v", endDate, s.EndDate)
	}

	// Check that velocity is 0 (default)
	if s.Velocity != 0 {
		t.Errorf("Expected Velocity to be 0, got %d", s.Velocity)
	}

	// Check that committed and completed lists are empty (default)
	if len(s.Committed) > 0 {
		t.Errorf("Expected Committed list to be empty, got %v", s.Committed)
	}
	if len(s.Completed) > 0 {
		t.Errorf("Expected Completed list to be empty, got %v", s.Completed)
	}

	// Check that burn-down chart is empty (default)
	if len(s.BurnDown) > 0 {
		t.Errorf("Expected BurnDown list to be empty, got %v", s.BurnDown)
	}
}

func TestSprintAddCommittedStory(t *testing.T) {
	s := NewSprint("Implement user authentication", time.Now(), time.Now().Add(7*24*time.Hour))

	// Add a story to committed list
	s.AddCommittedStory("story-1")

	// Check that the story was added
	if len(s.Committed) != 1 {
		t.Errorf("Expected Committed length to be 1, got %d", len(s.Committed))
	}
	if s.Committed[0] != "story-1" {
		t.Errorf("Expected first committed story to be 'story-1', got %s", s.Committed[0])
	}

	// Add the same story again (should not duplicate)
	s.AddCommittedStory("story-1")

	// Check that there's still only one instance of the story
	if len(s.Committed) != 1 {
		t.Errorf("Expected Committed length to be 1 after adding duplicate, got %d", len(s.Committed))
	}
}

func TestSprintRemoveCommittedStory(t *testing.T) {
	s := NewSprint("Implement user authentication", time.Now(), time.Now().Add(7*24*time.Hour))
	s.AddCommittedStory("story-1")
	s.AddCommittedStory("story-2")

	// Remove a story that exists
	s.RemoveCommittedStory("story-1")

	// Check that the story was removed and only story-2 remains
	if len(s.Committed) != 1 {
		t.Errorf("Expected Committed length to be 1 after removal, got %d", len(s.Committed))
	}
	if s.Committed[0] != "story-2" {
		t.Errorf("Expected remaining committed story to be 'story-2', got %s", s.Committed[0])
	}

	// Try to remove a story that doesn't exist
	s.RemoveCommittedStory("story-3")

	// Check that the committed list remains unchanged
	if len(s.Committed) != 1 {
		t.Errorf("Expected Committed length to be 1 after trying to remove non-existent story, got %d", len(s.Committed))
	}
}

func TestSprintAddCompletedStory(t *testing.T) {
	s := NewSprint("Implement user authentication", time.Now(), time.Now().Add(7*24*time.Hour))

	// Add a story to completed list
	s.AddCompletedStory("story-1")

	// Check that the story was added
	if len(s.Completed) != 1 {
		t.Errorf("Expected Completed length to be 1, got %d", len(s.Completed))
	}
	if s.Completed[0] != "story-1" {
		t.Errorf("Expected first completed story to be 'story-1', got %s", s.Completed[0])
	}

	// Add the same story again (should not duplicate)
	s.AddCompletedStory("story-1")

	// Check that there's still only one instance of the story
	if len(s.Completed) != 1 {
		t.Errorf("Expected Completed length to be 1 after adding duplicate, got %d", len(s.Completed))
	}
}

func TestSprintSetVelocity(t *testing.T) {
	s := NewSprint("Implement user authentication", time.Now(), time.Now().Add(7*24*time.Hour))

	// Set velocity to 10
	s.SetVelocity(10)

	// Check that the velocity was set correctly
	if s.Velocity != 10 {
		t.Errorf("Expected Velocity to be 10, got %d", s.Velocity)
	}
}

func TestSprintAddBurnDownPoint(t *testing.T) {
	s := NewSprint("Implement user authentication", time.Now(), time.Now().Add(7*24*time.Hour))

	// Add a burn-down point
	date := time.Date(2025, 9, 26, 10, 0, 0, 0, time.UTC)
	left := 8
	s.AddBurnDownPoint(date, left)

	// Check that the point was added
	if len(s.BurnDown) != 1 {
		t.Errorf("Expected BurnDown length to be 1, got %d", len(s.BurnDown))
	}
	if !s.BurnDown[0].Date.Equal(date) {
		t.Errorf("Expected BurnDown point date to be %v, got %v", date, s.BurnDown[0].Date)
	}
	if s.BurnDown[0].Left != left {
		t.Errorf("Expected BurnDown point Left value to be %d, got %d", left, s.BurnDown[0].Left)
	}
}

func TestSprintContains(t *testing.T) {
	s := NewSprint("Implement user authentication", time.Now(), time.Now().Add(7*24*time.Hour))

	// Check that contains function works correctly
	if contains(s.Committed, "story-1") {
		t.Errorf("Expected contains to return false for non-existent story")
	}

	s.AddCommittedStory("story-1")

	if !contains(s.Committed, "story-1") {
		t.Errorf("Expected contains to return true for existing story")
	}
}
