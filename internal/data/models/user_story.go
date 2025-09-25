// UserStory represents a user story in the agile development process
package models

import (
	"math/rand"
	"strconv"
	"time"
)

// PriorityLevel represents the priority of a user story
type PriorityLevel string

const (
	PriorityHigh   PriorityLevel = "High"
	PriorityMedium PriorityLevel = "Medium"
	PriorityLow    PriorityLevel = "Low"
)

// StoryStatus represents the status of a user story
type StoryStatus string

const (
	StoryDraft      StoryStatus = "Draft"
	StoryReady      StoryStatus = "Ready"
	StoryInProgress StoryStatus = "InProgress"
	StoryDone       StoryStatus = "Done"
)

// UserStory represents a user story in the agile development process
type UserStory struct {
	ID                 string         `json:"id"`
	Title              string         `json:"title"`
	Description        string         `json:"description"`
	AcceptanceCriteria []string       `json:"acceptance_criteria"`
	BusinessValue      int            `json:"business_value"` // 1-10
	Priority           PriorityLevel  `json:"priority"`
	Status             StoryStatus    `json:"status"`
	Estimate           *StoryEstimate `json:"estimate"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}

// StoryEstimate represents the estimation of a user story
type StoryEstimate struct {
	Points int `json:"points"` // Story points estimate
}

// NewUserStory creates a new UserStory with default values
func NewUserStory(title, description string) *UserStory {
	return &UserStory{
		ID:                 generateIDUS(),
		Title:              title,
		Description:        description,
		AcceptanceCriteria: []string{},
		BusinessValue:      5,
		Priority:           PriorityMedium,
		Status:             StoryDraft,
		Estimate:           &StoryEstimate{Points: 0},
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
}

func generateIDUS() string {
	return "story-" + time.Now().Format("20060102150405.000000") + "." + strconv.FormatInt(rand.Int63(), 10)
}

// Update updates the user story with new values
func (us *UserStory) Update(title, description string, priority PriorityLevel, businessValue int) {
	us.Title = title
	us.Description = description
	us.Priority = priority
	us.BusinessValue = businessValue
	us.UpdatedAt = time.Now()
}

// SetStatus updates the status of the user story
func (us *UserStory) SetStatus(status StoryStatus) {
	us.Status = status
	us.UpdatedAt = time.Now()
}

// AddAcceptanceCriterion adds a new acceptance criterion
func (us *UserStory) AddAcceptanceCriterion(criterion string) {
	us.AcceptanceCriteria = append(us.AcceptanceCriteria, criterion)
	us.UpdatedAt = time.Now()
}

// SetEstimate sets the story point estimate
func (us *UserStory) SetEstimate(points int) {
	us.Estimate.Points = points
	us.UpdatedAt = time.Now()
}
