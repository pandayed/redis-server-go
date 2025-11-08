package main

import (
	"log"
	"net"
	"sync"
)

type ConnectionManager struct {
	count int
	mu    sync.Mutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{}
}

func (cm *ConnectionManager) Increment(addr net.Addr) {
	cm.mu.Lock()
	cm.count++
	count := cm.count
	cm.mu.Unlock()
	log.Printf("New connection from %s. Total connections: %d", addr, count)
}

func (cm *ConnectionManager) Decrement(addr net.Addr) {
	cm.mu.Lock()
	cm.count--
	count := cm.count
	cm.mu.Unlock()
	log.Printf("Connection from %s closed. Total connections: %d", addr, count)
}

