// Package communication provides a message bus for agent-to-agent communication
// in the agile team intelligent agent system.
package communication

import (
	"encoding/json"
	"net"
	"sync"
	"time"
)

// AgentType represents the type of an agent in the system
type AgentType string

const (
	POAgent  AgentType = "po"
	DevAgent AgentType = "dev"
	SMAgent  AgentType = "sm"
)

// PriorityLevel represents the priority level of a message
type PriorityLevel int

const (
	LowPriority PriorityLevel = iota
	MediumPriority
	HighPriority
)

// MessageType represents the type of communication message
type MessageType string

const (
	MsgStoryCreated      MessageType = "story.created"
	MsgBacklogUpdated    MessageType = "backlog.updated"
	MsgAcceptanceRequest MessageType = "acceptance.request"
	MsgTaskBreakdown     MessageType = "task.breakdown"
	MsgProgressUpdate    MessageType = "progress.update"
	MsgObstacleReport    MessageType = "obstacle.report"
	MsgSprintStart       MessageType = "sprint.start"
	MsgDailyStandup      MessageType = "daily.standup"
	MsgRetrospective     MessageType = "retrospective"
	MsgAck               MessageType = "acknowledgment"
	MsgError             MessageType = "error"
)

// AgentMessage represents a message sent between agents
type AgentMessage struct {
	ID          string        `json:"id"`
	Timestamp   time.Time     `json:"timestamp"`
	From        AgentType     `json:"from"`
	To          AgentType     `json:"to"`
	Type        MessageType   `json:"type"`
	Priority    PriorityLevel `json:"priority"`
	Correlation string        `json:"correlation"`
	Payload     []byte        `json:"payload"`
}

// AgentBus defines the interface for the message bus
// that handles publishing and subscribing to messages between agents.
type AgentBus interface {
	// Publish sends a message to all subscribers of the specified type
	publish(message *AgentMessage) error

	// Subscribe registers a handler function for messages of a specific type
	subscribe(messageType MessageType, handler func(*AgentMessage)) error

	// Unsubscribe removes a handler for messages of a specific type
	unsubscribe(messageType MessageType, handler func(*AgentMessage)) error

	// Close shuts down the message bus and releases resources
	close() error
}

// communicationError implements the error interface
type communicationError string

func (e communicationError) Error() string {
	return string(e)
}

// ErrBusClosed is returned when trying to use a closed message bus
var ErrBusClosed = communicationError("message bus is closed")

// agentBus implements the AgentBus interface using Unix Domain Sockets
type agentBus struct {
	socket   net.Conn
	handlers map[MessageType][]func(*AgentMessage)
	mutex    sync.RWMutex
	closed   bool
}

// New creates a new instance of AgentBus with Unix Domain Socket communication
func New(socketPath string) (AgentBus, error) {
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, err
	}

	a := &agentBus{
		handlers: make(map[MessageType][]func(*AgentMessage)),
		socket:   conn,
	}

	// Start a goroutine to handle incoming messages
	go a.handleMessages()

	return a, nil
}

// publish sends a message to all subscribers of the specified type
func (a *agentBus) publish(message *AgentMessage) error {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	if a.closed {
		return ErrBusClosed
	}

	// Serialize the message to JSON
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Send the message through the socket
	_, err = a.socket.Write(data)
	return err
}

// subscribe registers a handler function for messages of a specific type
func (a *agentBus) subscribe(messageType MessageType, handler func(*AgentMessage)) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.closed {
		return ErrBusClosed
	}

	// Add the handler to the map of handlers for this message type
	a.handlers[messageType] = append(a.handlers[messageType], handler)

	return nil
}

// unsubscribe removes a handler for messages of a specific type
func (a *agentBus) unsubscribe(messageType MessageType, handler func(*AgentMessage)) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.closed {
		return ErrBusClosed
	}

	// Find and remove the handler from the map of handlers for this message type
	handlers := a.handlers[messageType]
	for i, h := range handlers {
		// Use pointer comparison since functions can't be compared with ==
		if &h == &handler {
			a.handlers[messageType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}

	return nil
}

// close shuts down the message bus and releases resources
func (a *agentBus) close() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.closed {
		return nil // Already closed
	}

	a.closed = true

	// Close the socket connection
	err := a.socket.Close()

	// Clear all handlers
	a.handlers = make(map[MessageType][]func(*AgentMessage))

	return err
}

// handleMessages reads incoming messages from the socket and dispatches them to appropriate handlers
func (a *agentBus) handleMessages() {
	buffer := make([]byte, 4096)

	for {
		// Read data from the socket
		n, err := a.socket.Read(buffer)
		if err != nil {
			// Handle connection error or closure
			break
		}

		// Process the received message
		messageData := buffer[:n]
		var msg AgentMessage
		err = json.Unmarshal(messageData, &msg)
		if err != nil {
			continue // Skip malformed messages
		}

		// Create a copy of the message to pass to handlers
		copiedMsg := msg

		// Dispatch to appropriate handlers
		a.mutex.RLock()
		handlers := a.handlers[msg.Type]
		a.mutex.RUnlock()

		for _, handler := range handlers {
			handler(&copiedMsg)
		}
	}
}
