package main

import (
	"log"
	"sync"
	"time"
)

// OperationType represents the type of operation being measured
type OperationType string

const (
	// HTTPOperation represents an HTTP request/response operation
	HTTPOperation OperationType = "HTTP"
	// DBOperation represents a database operation
	DBOperation OperationType = "DB"
	// WSOperation represents a WebSocket operation
	WSOperation OperationType = "WS"
)

// Performance tracks operation metrics for performance monitoring
type Performance struct {
	metrics  map[OperationType][]time.Duration
	mutex    sync.RWMutex
	logEvery int
	counter  map[OperationType]int
}

// NewPerformance creates a new Performance tracker
func NewPerformance(logEveryN int) *Performance {
	return &Performance{
		metrics:  make(map[OperationType][]time.Duration),
		mutex:    sync.RWMutex{},
		logEvery: logEveryN,
		counter:  make(map[OperationType]int),
	}
}

// Track measures the duration of a function execution and logs the performance
func (p *Performance) Track(opType OperationType, fn func()) time.Duration {
	start := time.Now()
	fn()
	duration := time.Since(start)
	
	p.mutex.Lock()
	defer p.mutex.Unlock()
	
	p.metrics[opType] = append(p.metrics[opType], duration)
	p.counter[opType]++
	
	if p.counter[opType]%p.logEvery == 0 {
		p.logPerformanceStats(opType)
	}
	
	// Log warnings for operations exceeding thresholds
	switch opType {
	case DBOperation:
		if duration > 100*time.Millisecond {
			log.Printf("WARNING: DB operation took %v, exceeding target of 100ms", duration)
		}
	case HTTPOperation, WSOperation:
		if duration > 300*time.Millisecond {
			log.Printf("WARNING: %s operation took %v, exceeding target of 300ms", opType, duration)
		}
	}
	
	return duration
}

// logPerformanceStats calculates and logs performance statistics
func (p *Performance) logPerformanceStats(opType OperationType) {
	if len(p.metrics[opType]) == 0 {
		return
	}
	
	var total time.Duration
	var max time.Duration
	for _, d := range p.metrics[opType] {
		total += d
		if d > max {
			max = d
		}
	}
	
	avg := total / time.Duration(len(p.metrics[opType]))
	
	log.Printf("[Performance] %s operations - Count: %d, Avg: %v, Max: %v", 
		opType, 
		len(p.metrics[opType]), 
		avg,
		max)
	
	threshold := 100 * time.Millisecond
	if opType == HTTPOperation || opType == WSOperation {
		threshold = 300 * time.Millisecond
	}
	
	if avg > threshold {
		log.Printf("[Performance WARNING] %s operations average (%v) exceeding target of %v", 
			opType, avg, threshold)
	}
}
