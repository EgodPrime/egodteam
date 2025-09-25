package communication

import (
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

// TestNew creates a new AgentBus instance and verifies it works correctly
func TestNew(t *testing.T) {
	// Create a temporary socket path for testing
	socketPath := filepath.Join(os.TempDir(), "test-communication-bus.sock")
	defer os.Remove(socketPath) // Clean up after test

	// Start a server to listen on the socket
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		t.Fatalf("Failed to start listener: %v", err)
	}
	defer listener.Close()

	// Create a new AgentBus instance
	a, err := New(socketPath)
	if err != nil {
		t.Fatalf("Failed to create AgentBus: %v", err)
	}
	defer a.close()

	// Verify the bus was created successfully
	if a == nil {
		t.Error("AgentBus should not be nil")
	}

	// Test that we can publish and receive messages
	testMessage := &AgentMessage{
		ID:       "test-1",
		From:     POAgent,
		To:       DevAgent,
		Type:     MsgStoryCreated,
		Priority: HighPriority,
		Payload:  []byte(`{"title":"Test Story","description":"This is a test story"}`),
	}

	// Publish the message
	err = a.publish(testMessage)
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}

	// Accept a connection from the client
	conn, err := listener.Accept()
	if err != nil {
		t.Fatalf("Failed to accept connection: %v", err)
	}
	defer conn.Close()

	// Read the message from the socket
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}

	// Parse the JSON message
	var receivedMessage AgentMessage
	err = json.Unmarshal(buffer[:n], &receivedMessage)
	if err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}

	// Verify the message content
	if receivedMessage.ID != testMessage.ID {
		t.Errorf("Expected ID %s, got %s", testMessage.ID, receivedMessage.ID)
	}
	if receivedMessage.From != testMessage.From {
		t.Errorf("Expected From %s, got %s", testMessage.From, receivedMessage.From)
	}
	if receivedMessage.To != testMessage.To {
		t.Errorf("Expected To %s, got %s", testMessage.To, receivedMessage.To)
	}
	if receivedMessage.Type != testMessage.Type {
		t.Errorf("Expected Type %s, got %s", testMessage.Type, receivedMessage.Type)
	}
	if receivedMessage.Priority != testMessage.Priority {
		t.Errorf("Expected Priority %d, got %d", testMessage.Priority, receivedMessage.Priority)
	}

	// Verify the payload is correct
	if string(receivedMessage.Payload) != string(testMessage.Payload) {
		t.Errorf("Expected Payload %s, got %s", testMessage.Payload, receivedMessage.Payload)
	}
}

// TestSubscribeUnsubscribe verifies that subscribing and unsubscribing works correctly
func TestSubscribeUnsubscribe(t *testing.T) {
	// Create a temporary socket path for testing
	socketPath := filepath.Join(os.TempDir(), "test-communication-bus.sock")
	defer os.Remove(socketPath) // Clean up after test

	// Start a server to listen on the socket
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		t.Fatalf("Failed to start listener: %v", err)
	}
	defer listener.Close()

	// Create a new AgentBus instance
	a, err := New(socketPath)
	if err != nil {
		t.Fatalf("Failed to create AgentBus: %v", err)
	}
	defer a.close()

	// Define a test handler function
	var receivedMessage *AgentMessage
	var messageReceived bool
	var receivedMutex sync.Mutex

	handler := func(msg *AgentMessage) {
		t.Logf("Handler called with message ID: %s", msg.ID)
		receivedMutex.Lock()
		defer receivedMutex.Unlock()

		receivedMessage = msg
		messageReceived = true
	}

	// Subscribe to messages of type MsgStoryCreated
	err = a.subscribe(MsgStoryCreated, handler)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	// Publish a message of the subscribed type
	testMessage := &AgentMessage{
		ID:       "test-2",
		From:     POAgent,
		To:       DevAgent,
		Type:     MsgStoryCreated,
		Priority: HighPriority,
		Payload:  []byte(`{"title":"Test Story","description":"This is a test story"}`),
	}

	// Publish the message
	err = a.publish(testMessage)
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}

	// Accept a connection from the client
	conn, err := listener.Accept()
	if err != nil {
		t.Fatalf("Failed to accept connection: %v", err)
	}
	defer conn.Close()

	// Read the message from the socket
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}

	// Parse the JSON message
	var receivedMessageFromSocket AgentMessage
	err = json.Unmarshal(buffer[:n], &receivedMessageFromSocket)
	if err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}

	// Verify the handler was called - wait for it to complete with timeout
	timeout := time.After(5 * time.Second)
	tick := time.NewTicker(10 * time.Millisecond)
	defer tick.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Timeout waiting for handler to be called")
		case <-tick.C:
			receivedMutex.Lock()
			if messageReceived {
				receivedMutex.Unlock()
				break
			}
			receivedMutex.Unlock()
		}
	}

	// Check if handler was called
	receivedMutex.Lock()
	defer receivedMutex.Unlock()

	if !messageReceived {
		t.Error("Handler should have been called")
	}
	if receivedMessage == nil {
		t.Error("Received message should not be nil")
	}
	if receivedMessage.ID != testMessage.ID {
		t.Errorf("Expected ID %s, got %s", testMessage.ID, receivedMessage.ID)
	}

	// Unsubscribe from messages of type MsgStoryCreated
	err = a.unsubscribe(MsgStoryCreated, handler)
	if err != nil {
		t.Fatalf("Failed to unsubscribe: %v", err)
	}

	// Publish another message of the same type
	testMessage2 := &AgentMessage{
		ID:       "test-3",
		From:     POAgent,
		To:       DevAgent,
		Type:     MsgStoryCreated,
		Priority: HighPriority,
		Payload:  []byte(`{"title":"Test Story 2","description":"This is another test story"}`),
	}

	// Publish the message
	err = a.publish(testMessage2)
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}

	// Accept a connection from the client
	conn2, err := listener.Accept()
	if err != nil {
		t.Fatalf("Failed to accept connection: %v", err)
	}
	defer conn2.Close()

	// Read the message from the socket
	n2, err := conn2.Read(buffer)
	if err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}

	// Parse the JSON message
	var receivedMessageFromSocket2 AgentMessage
	err = json.Unmarshal(buffer[:n2], &receivedMessageFromSocket2)
	if err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}

	// Verify the handler was not called (should be unsubscribed) - wait for it to complete with timeout
	timeout = time.After(5 * time.Second)
	tick = time.NewTicker(10 * time.Millisecond)
	defer tick.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Timeout waiting for handler not to be called")
		case <-tick.C:
			receivedMutex.Lock()
			if !messageReceived {
				receivedMutex.Unlock()
				break
			}
			receivedMutex.Unlock()
		}
	}

	// Check if handler was not called after unsubscribe
	receivedMutex.Lock()
	defer receivedMutex.Unlock()

	if messageReceived {
		t.Error("Handler should not have been called after unsubscribe")
	}
}

// TestClose verifies that closing the bus works correctly
func TestClose(t *testing.T) {
	// Create a temporary socket path for testing
	socketPath := filepath.Join(os.TempDir(), "test-communication-bus.sock")
	defer os.Remove(socketPath) // Clean up after test

	// Start a server to listen on the socket
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		t.Fatalf("Failed to start listener: %v", err)
	}
	defer listener.Close()

	// Create a new AgentBus instance
	a, err := New(socketPath)
	if err != nil {
		t.Fatalf("Failed to create AgentBus: %v", err)
	}

	// Close the bus
	err = a.close()
	if err != nil {
		t.Fatalf("Failed to close AgentBus: %v", err)
	}

	// Verify that we can't publish after closing
	testMessage := &AgentMessage{
		ID:       "test-4",
		From:     POAgent,
		To:       DevAgent,
		Type:     MsgStoryCreated,
		Priority: HighPriority,
		Payload:  []byte(`{"title":"Test Story","description":"This is a test story"}`),
	}

	// Try to publish after closing - should fail
	err = a.publish(testMessage)
	if err == nil {
		t.Error("Expected error when publishing after bus is closed")
	}

	// Verify the error is ErrBusClosed
	if err != ErrBusClosed {
		t.Errorf("Expected error %v, got %v", ErrBusClosed, err)
	}

	// Try to subscribe after closing - should fail
	handler := func(*AgentMessage) {}

	err = a.subscribe(MsgStoryCreated, handler)
	if err == nil {
		t.Error("Expected error when subscribing after bus is closed")
	}

	// Verify the error is ErrBusClosed
	if err != ErrBusClosed {
		t.Errorf("Expected error %v, got %v", ErrBusClosed, err)
	}

	// Try to unsubscribe after closing - should fail
	err = a.unsubscribe(MsgStoryCreated, handler)
	if err == nil {
		t.Error("Expected error when unsubscribing after bus is closed")
	}

	// Verify the error is ErrBusClosed
	if err != ErrBusClosed {
		t.Errorf("Expected error %v, got %v", ErrBusClosed, err)
	}
}
