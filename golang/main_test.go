package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestParseFlags(t *testing.T) {
	// Create a new test flag set
	testFlags := flag.NewFlagSet("test", flag.ContinueOnError)
	
	// Define flags for the test flag set
	rate := testFlags.Int("rate", 60, "")
	host := testFlags.String("host", "localhost", "")
	port := testFlags.Int("port", 8080, "")
	worker := testFlags.Int("worker", 1, "")
	randomRate := testFlags.Bool("random-rate", false, "")
	
	// Parse test arguments
	err := testFlags.Parse([]string{
		"--rate=120", 
		"--host=127.0.0.1", 
		"--port=9090", 
		"--worker=4",
		"--random-rate=true",
	})
	
	if err != nil {
		t.Fatalf("Failed to parse flags: %v", err)
	}

	// Check that the values were parsed correctly
	if *rate != 120 {
		t.Errorf("Expected rate 120, got %d", *rate)
	}
	if *host != "127.0.0.1" {
		t.Errorf("Expected host 127.0.0.1, got %s", *host)
	}
	if *port != 9090 {
		t.Errorf("Expected port 9090, got %d", *port)
	}
	if *worker != 4 {
		t.Errorf("Expected worker 4, got %d", *worker)
	}
	if !*randomRate {
		t.Errorf("Expected randomRate true, got %v", *randomRate)
	}
}

func TestRandomDataGenerator(t *testing.T) {
	generator := NewRandomDataGenerator()
	
	// Test multiple data generations
	for i := 0; i < 10; i++ {
		data := generator.GetData()
		
		// Check format: StockCode|Price|B/S|Amount
		parts := strings.Split(data, "|")
		if len(parts) != 4 {
			t.Errorf("Expected 4 parts in data, got %d: %s", len(parts), data)
		}
		
		// Check if stock code is one of the expected ones
		stockCodeFound := false
		for _, code := range generator.stockCodes {
			if parts[0] == code {
				stockCodeFound = true
				break
			}
		}
		if !stockCodeFound {
			t.Errorf("Stock code %s not in expected list", parts[0])
		}
		
		// Check buy/sell indicator
		if parts[2] != "B" && parts[2] != "S" {
			t.Errorf("Expected B or S for buy/sell indicator, got %s", parts[2])
		}
	}
}

func TestPerformanceTracking(t *testing.T) {
	perf := NewPerformance(1) // Log every operation
	
	// Test HTTP operation tracking
	duration := perf.Track(HTTPOperation, func() {
		time.Sleep(10 * time.Millisecond)
	})
	
	if duration < 10*time.Millisecond {
		t.Errorf("Expected duration >= 10ms, got %v", duration)
	}
	
	// Test metrics recording
	if len(perf.metrics[HTTPOperation]) != 1 {
		t.Errorf("Expected 1 HTTP operation recorded, got %d", len(perf.metrics[HTTPOperation]))
	}
}

func TestDataHub(t *testing.T) {
	hub := NewDataHub(60, false)
	
	// Test adding connection
	conn := &Connection{
		id:         1,
		data:       make(chan string, 10),
		performance: hub.performance,
	}
	
	hub.AddConnection(conn)
	
	if len(hub.connections) != 1 {
		t.Errorf("Expected 1 connection, got %d", len(hub.connections))
	}
	
	// Test removing connection
	hub.RemoveConnection(conn)
	
	if len(hub.connections) != 0 {
		t.Errorf("Expected 0 connections, got %d", len(hub.connections))
	}
}

func TestHealthEndpoint(t *testing.T) {
	hub := NewDataHub(60, false)
	
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hub.performance.Track(HTTPOperation, func() {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "OK")
		})
	}).ServeHTTP(w, req)
	
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %d", resp.StatusCode)
	}
}
