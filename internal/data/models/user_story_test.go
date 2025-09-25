package models

import (
	"testing"
	"time"
)

func TestNewUserStory(t *testing.T) {
	title := "Implement login feature"
	description := "Create a secure login system for users"

	us := NewUserStory(title, description)

	// Check that the story has a valid ID
	if us.ID == "" || len(us.ID) < 5 {
		t.Errorf("Expected non-empty ID, got %s", us.ID)
	}

	// Check that title and description are set correctly
	if us.Title != title {
		t.Errorf("Expected Title to be '%s', got %s", title, us.Title)
	}
	if us.Description != description {
		t.Errorf("Expected Description to be '%s', got %s", description, us.Description)
	}

	// Check that acceptance criteria list is empty (default)
	if len(us.AcceptanceCriteria) > 0 {
		t.Errorf("Expected AcceptanceCriteria length to be 0, got %d", len(us.AcceptanceCriteria))
	}

	// Check that business value is 5 (default)
	if us.BusinessValue != 5 {
		t.Errorf("Expected BusinessValue to be 5, got %d", us.BusinessValue)
	}

	// Check that priority is Medium (default)
	if us.Priority != PriorityMedium {
		t.Errorf("Expected Priority to be PriorityMedium, got %s", us.Priority)
	}

	// Check that status is Draft (default)
	if us.Status != StoryDraft {
		t.Errorf("Expected Status to be StoryDraft, got %s", us.Status)
	}

	// Check that estimate points are 0 (default)
	if us.Estimate.Points != 0 {
		t.Errorf("Expected Estimate Points to be 0, got %d", us.Estimate.Points)
	}

	// Check that created and updated times are set correctly
	now := time.Now()
	if us.CreatedAt.IsZero() || us.UpdatedAt.IsZero() {
		t.Errorf("Expected CreatedAt and UpdatedAt to be non-zero")
	}
	// Allow for small difference due to timing
	if us.CreatedAt.After(now.Add(time.Second)) || us.UpdatedAt.After(now.Add(time.Second)) {
		t.Errorf("Expected CreatedAt and UpdatedAt to be close to current time, got %v and %v", us.CreatedAt, us.UpdatedAt)
	}
}

func TestUserStoryUpdate(t *testing.T) {
	us := NewUserStory("Implement login feature", "Create a secure login system for users")

	// Update the story with new values
	newTitle := "Fix authentication bug"
	newDescription := "Improve security of login system"
	newPriority := PriorityHigh
	newBusinessValue := 8

	us.Update(newTitle, newDescription, newPriority, newBusinessValue)

	// Check that all fields were updated correctly
	if us.Title != newTitle {
		t.Errorf("Expected Title to be '%s', got %s", newTitle, us.Title)
	}
	if us.Description != newDescription {
		t.Errorf("Expected Description to be '%s', got %s", newDescription, us.Description)
	}
	if us.Priority != newPriority {
		t.Errorf("Expected Priority to be %s, got %s", newPriority, us.Priority)
	}
	if us.BusinessValue != newBusinessValue {
		t.Errorf("Expected BusinessValue to be %d, got %d", newBusinessValue, us.BusinessValue)
	}

	// Check that updated time was set correctly
	now := time.Now()
	if us.UpdatedAt.IsZero() || us.UpdatedAt.Before(us.CreatedAt) {
		t.Errorf("Expected UpdatedAt to be non-zero and after CreatedAt")
	}
	// Allow for small difference due to timing
	if us.UpdatedAt.After(now.Add(time.Second)) {
		t.Errorf("Expected UpdatedAt to be close to current time, got %v", us.UpdatedAt)
	}
}

func TestUserStorySetStatus(t *testing.T) {
	us := NewUserStory("Implement login feature", "Create a secure login system for users")

	// Set status to InProgress
	newStatus := StoryInProgress
	us.SetStatus(newStatus)

	// Check that the status was updated correctly
	if us.Status != newStatus {
		t.Errorf("Expected Status to be %s, got %s", newStatus, us.Status)
	}

	// Check that updated time was set correctly
	now := time.Now()
	if us.UpdatedAt.IsZero() || us.UpdatedAt.Before(us.CreatedAt) {
		t.Errorf("Expected UpdatedAt to be non-zero and after CreatedAt")
	}
	// Allow for small difference due to timing
	if us.UpdatedAt.After(now.Add(time.Second)) {
		t.Errorf("Expected UpdatedAt to be close to current time, got %v", us.UpdatedAt)
	}
}

func TestUserStoryAddAcceptanceCriterion(t *testing.T) {
	us := NewUserStory("Implement login feature", "Create a secure login system for users")

	// Add an acceptance criterion
	criterion1 := "User can log in with valid credentials"
	us.AddAcceptanceCriterion(criterion1)

	// Check that the criterion was added
	if len(us.AcceptanceCriteria) != 1 {
		t.Errorf("Expected AcceptanceCriteria length to be 1, got %d", len(us.AcceptanceCriteria))
	}
	if us.AcceptanceCriteria[0] != criterion1 {
		t.Errorf("Expected first acceptance criterion to be '%s', got %s", criterion1, us.AcceptanceCriteria[0])
	}

	// Add another criterion
	criterion2 := "User cannot log in with invalid credentials"
	us.AddAcceptanceCriterion(criterion2)

	// Check that the second criterion was added
	if len(us.AcceptanceCriteria) != 2 {
		t.Errorf("Expected AcceptanceCriteria length to be 2, got %d", len(us.AcceptanceCriteria))
	}
	if us.AcceptanceCriteria[1] != criterion2 {
		t.Errorf("Expected second acceptance criterion to be '%s', got %s", criterion2, us.AcceptanceCriteria[1])
	}

	// Check that updated time was set correctly
	now := time.Now()
	if us.UpdatedAt.IsZero() || us.UpdatedAt.Before(us.CreatedAt) {
		t.Errorf("Expected UpdatedAt to be non-zero and after CreatedAt")
	}
	// Allow for small difference due to timing
	if us.UpdatedAt.After(now.Add(time.Second)) {
		t.Errorf("Expected UpdatedAt to be close to current time, got %v", us.UpdatedAt)
	}
}

func TestUserStorySetEstimate(t *testing.T) {
	us := NewUserStory("Implement login feature", "Create a secure login system for users")

	// Set estimate to 5 points
	points := 5
	us.SetEstimate(points)

	// Check that the estimate was set correctly
	if us.Estimate.Points != points {
		t.Errorf("Expected Estimate Points to be %d, got %d", points, us.Estimate.Points)
	}

	// Check that updated time was set correctly
	now := time.Now()
	if us.UpdatedAt.IsZero() || us.UpdatedAt.Before(us.CreatedAt) {
		t.Errorf("Expected UpdatedAt to be non-zero and after CreatedAt")
	}
	// Allow for small difference due to timing
	if us.UpdatedAt.After(now.Add(time.Second)) {
		t.Errorf("Expected UpdatedAt to be close to current time, got %v", us.UpdatedAt)
	}
}

func TestUserStorygenerateIDUS(t *testing.T) {
	// Test that generateIDUS creates unique IDs
	id1 := generateIDUS()
	id2 := generateIDUS()

	// Check that the IDs are different
	if id1 == id2 {
		t.Errorf("Expected unique IDs, got %s and %s", id1, id2)
	}

	// Check that the ID starts with "story-"
	if !startsWith(id1, "story-") || !startsWith(id2, "story-") {
		t.Errorf("Expected IDs to start with 'story-', got %s and %s", id1, id2)
	}

	// Check that the ID has a timestamp format (YYYYMMDDHHMMSS)
	if len(id1) < 15 || len(id2) < 15 {
		t.Errorf("Expected IDs to have timestamp format, got %s and %s", id1, id2)
	}
}
