package main

import (
	"sync"
	"time"
)

type Status = int

const (
	Unknown Status = iota
	Connected
	Disconnected
	Reconnecting
)

type ConnectionInfo struct {
	mu                      sync.RWMutex
	status                  Status
	lastError               error
	lastPingTime            time.Time
	reconnectTime           time.Time
	allReconnectTriesFailed bool
}

// GetStatus returns the current connection status
func (ci *ConnectionInfo) GetStatus() Status {
	ci.mu.RLock()
	defer ci.mu.RUnlock()
	return ci.status
}

// SetStatus sets the connection status
func (ci *ConnectionInfo) SetStatus(status Status) {
	ci.mu.Lock()
	defer ci.mu.Unlock()
	ci.status = status
}

// GetLastError returns the last error encountered
func (ci *ConnectionInfo) GetLastError() error {
	ci.mu.RLock()
	defer ci.mu.RUnlock()
	return ci.lastError
}

// SetLastError sets the last error encountered
func (ci *ConnectionInfo) SetLastError(err error) {
	ci.mu.Lock()
	defer ci.mu.Unlock()
	ci.lastError = err
}

// GetLastPingTime returns the time of the last successful ping
func (ci *ConnectionInfo) GetLastPingTime() time.Time {
	ci.mu.RLock()
	defer ci.mu.RUnlock()
	return ci.lastPingTime
}

// SetLastPingTime sets the time of the last successful ping
func (ci *ConnectionInfo) SetLastPingTime() {
	ci.mu.Lock()
	defer ci.mu.Unlock()
	ci.lastPingTime = time.Now()
}

// GetReconnectTime returns the last reconnection attempt time
func (ci *ConnectionInfo) GetReconnectTime() time.Time {
	ci.mu.RLock()
	defer ci.mu.RUnlock()
	return ci.reconnectTime
}

// SetReconnectTime sets the last reconnection attempt time
func (ci *ConnectionInfo) SetReconnectTime() {
	ci.mu.Lock()
	defer ci.mu.Unlock()
	ci.reconnectTime = time.Now()
}

// GetAllReconnectTriesFailed returns whether all reconnection attempts have failed
func (ci *ConnectionInfo) GetAllReconnectTriesFailed() bool {
	ci.mu.RLock()
	defer ci.mu.RUnlock()
	return ci.allReconnectTriesFailed
}

// SetAllReconnectTriesFailed sets the status of reconnection attempts
func (ci *ConnectionInfo) SetAllReconnectTriesFailed(failed bool) {
	ci.mu.Lock()
	defer ci.mu.Unlock()
	ci.allReconnectTriesFailed = failed
}
