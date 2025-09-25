// Sprint represents an iteration in the agile development process
package models

import (
	"math/rand"
	"strconv"
	"time"
)

// BurnDownPoint represents a data point in the burn-down chart
type BurnDownPoint struct {
	Date time.Time `json:"date"`
	Left int       `json:"left"` // Remaining story points
}

// Sprint represents an iteration in the agile development process
type Sprint struct {
	ID        string          `json:"id"`
	Goal      string          `json:"goal"`
	StartDate time.Time       `json:"start_date"`
	EndDate   time.Time       `json:"end_date"`
	Velocity  int             `json:"velocity"`  // Team velocity
	Committed []string        `json:"committed"` // IDs of committed user stories
	Completed []string        `json:"completed"` // IDs of completed user stories
	BurnDown  []BurnDownPoint `json:"burn_down"` // Burn-down chart data
}

// generateIDS generates a unique ID for a sprint
func generateIDS() string {
	return "sprint-" + time.Now().Format("20060102150405.000000") + "." + strconv.FormatInt(rand.Int63(), 10)
}

// NewSprint creates a new sprint with default values
func NewSprint(goal string, startDate, endDate time.Time) *Sprint {
	return &Sprint{
		ID:        generateIDS(),
		Goal:      goal,
		StartDate: startDate,
		EndDate:   endDate,
		Velocity:  0,
		Committed: []string{},
		Completed: []string{},
		BurnDown:  []BurnDownPoint{},
	}
}

// AddCommittedStory adds a user story to the committed list
func (s *Sprint) AddCommittedStory(storyID string) {
	if !contains(s.Committed, storyID) {
		s.Committed = append(s.Committed, storyID)
	}
}

// RemoveCommittedStory removes a user story from the committed list
func (s *Sprint) RemoveCommittedStory(storyID string) {
	for i, id := range s.Committed {
		if id == storyID {
			s.Committed = append(s.Committed[:i], s.Committed[i+1:]...)
			break
		}
	}
}

// AddCompletedStory adds a user story to the completed list
func (s *Sprint) AddCompletedStory(storyID string) {
	if !contains(s.Completed, storyID) {
		s.Completed = append(s.Completed, storyID)
	}
}

// SetVelocity sets the team velocity
func (s *Sprint) SetVelocity(velocity int) {
	s.Velocity = velocity
}

// AddBurnDownPoint adds a data point to the burn-down chart
func (s *Sprint) AddBurnDownPoint(date time.Time, left int) {
	point := BurnDownPoint{
		Date: date,
		Left: left,
	}
	s.BurnDown = append(s.BurnDown, point)
}

// contains checks if a slice contains an element
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
