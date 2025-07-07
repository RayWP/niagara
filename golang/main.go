package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// DataHub manages WebSocket connections and broadcasts data
type DataHub struct {
	connections []*Connection
	rate        int
	randomRate  bool
	mutex       sync.Mutex
	generator   *RandomDataGenerator
	performance *Performance
	stopCh      chan struct{}
}

// Connection represents a WebSocket connection
type Connection struct {
	id         int
	conn       *websocket.Conn
	hub        *DataHub
	data       chan string
	performance *Performance
}

// NewDataHub creates a new DataHub with the specified rate
func NewDataHub(rate int, randomRate bool) *DataHub {
	perf := NewPerformance(100) // Log performance stats every 100 operations
	return &DataHub{
		rate:        rate,
		randomRate:  randomRate,
		generator:   NewRandomDataGenerator(),
		performance: perf,
		stopCh:      make(chan struct{}),
	}
}

// AddConnection adds a new connection to the hub
func (hub *DataHub) AddConnection(conn *Connection) {
	hub.performance.Track(WSOperation, func() {
		hub.mutex.Lock()
		defer hub.mutex.Unlock()
		log.Printf("New Connection added, id: %d", conn.id)
		hub.connections = append(hub.connections, conn)
		if len(hub.connections) == 1 {
			go hub.broadcastData()
		}
	})
}

// RemoveConnection removes a connection from the hub
func (hub *DataHub) RemoveConnection(conn *Connection) {
	hub.performance.Track(WSOperation, func() {
		hub.mutex.Lock()
		defer hub.mutex.Unlock()
		for i, c := range hub.connections {
			if c.id == conn.id {
				log.Printf("Connection removed, id: %d", conn.id)
				hub.connections = append(hub.connections[:i], hub.connections[i+1:]...)
				break
			}
		}
		if len(hub.connections) == 0 {
			log.Println("No more connections, stop generating data")
			close(hub.stopCh)
			hub.stopCh = make(chan struct{})
		}
	})
}

// broadcastData sends random data to all connected clients at the specified rate
func (hub *DataHub) broadcastData() {
	log.Println("Starting data broadcast")
	baseInterval := time.Duration(60/hub.rate) * time.Second
	ticker := time.NewTicker(baseInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			if hub.randomRate && rand.Intn(10) < 3 { // 30% chance of rate change
				newRate := hub.rate + rand.Intn(20) - 10 // +/- 10
				if newRate < 10 {
					newRate = 10 // Minimum rate
				}
				log.Printf("Chaos mode: Rate changed from %d to %d", hub.rate, newRate)
				hub.rate = newRate
				ticker.Reset(time.Duration(60/hub.rate) * time.Second)
			}
			
			// Generate and broadcast data with performance tracking
			hub.performance.Track(WSOperation, func() {
				data := hub.generator.GetData()
				log.Printf("Broadcasting data: %s", data)
				
				hub.mutex.Lock()
				for _, conn := range hub.connections {
					select {
					case conn.data <- data:
					default:
						// Non-blocking send to prevent slow clients from affecting others
						log.Printf("Warning: Client %d too slow, dropping message", conn.id)
					}
				}
				hub.mutex.Unlock()
			})
		case <-hub.stopCh:
			log.Println("Stopping data broadcast")
			return
		}
	}
}

// handleConnection processes incoming WebSocket connections
func (hub *DataHub) handleConnection(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all connections
		},
	}

	hub.performance.Track(HTTPOperation, func() {
		log.Printf("Received connection request from %s", r.RemoteAddr)
	})
	
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}

	// Create new connection
	connection := &Connection{
		id:         len(hub.connections) + 1,
		conn:       conn,
		hub:        hub,
		data:       make(chan string, 256), // Buffered channel to prevent blocking
		performance: hub.performance,
	}

	// Add connection to hub
	hub.AddConnection(connection)

	// Start message sender goroutine
	go connection.writePump()
	
	// Start message reader goroutine
	go connection.readPump()
}

// writePump sends messages to the WebSocket connection
func (conn *Connection) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		conn.conn.Close()
	}()

	for {
		select {
		case message, ok := <-conn.data:
			if !ok {
				// Channel closed, send close message
				conn.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			conn.performance.Track(WSOperation, func() {
				w, err := conn.conn.NextWriter(websocket.TextMessage)
				if err != nil {
					log.Printf("Error getting writer: %v", err)
					return
				}
				w.Write([]byte(message))
				
				// Add queued messages to the current websocket message
				n := len(conn.data)
				for i := 0; i < n; i++ {
					w.Write([]byte("\n"))
					w.Write([]byte(<-conn.data))
				}
				
				if err := w.Close(); err != nil {
					log.Printf("Error closing writer: %v", err)
					return
				}
			})
			
		case <-ticker.C:
			// Send ping
			if err := conn.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
				log.Printf("Error sending ping: %v", err)
				return
			}
		}
	}
}

// readPump reads messages from the WebSocket connection
func (conn *Connection) readPump() {
	defer func() {
		conn.hub.RemoveConnection(conn)
		conn.conn.Close()
	}()
	
	conn.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.conn.SetPongHandler(func(string) error { 
		conn.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil 
	})
	
	for {
		_, _, err := conn.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}
	}
}

// parseFlags parses command-line flags and returns configuration values
func parseFlags() (int, string, int, int, bool) {
	rate := flag.Int("rate", 60, "Data generation rate per minute")
	host := flag.String("host", "localhost", "Server host")
	port := flag.Int("port", 8080, "Server port")
	worker := flag.Int("worker", 1, "Number of worker processes")
	randomRate := flag.Bool("random-rate", false, "Enable random rate changes (chaos mode)")
	
	// Parse flags
	flag.Parse()
	
	return *rate, *host, *port, *worker, *randomRate
}

func main() {
	// Configure logging
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC | log.Lshortfile)
	
	// Parse command-line flags
	rate, host, port, workers, randomRate := parseFlags()
	
	// Create data hub
	hub := NewDataHub(rate, randomRate)
	
	// Set up HTTP routes
	http.HandleFunc("/connect", hub.handleConnection)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		hub.performance.Track(HTTPOperation, func() {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "OK")
		})
	})
	
	// Set up static file serving
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	
	// Create static directory and index.html if they don't exist
	if _, err := os.Stat("./static"); os.IsNotExist(err) {
		os.Mkdir("./static", 0755)
		
		indexHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Niagara - Go Implementation</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
        h1 { color: #333; }
        #messages { border: 1px solid #ccc; height: 300px; overflow-y: auto; padding: 10px; margin-bottom: 20px; }
        .message { margin: 5px 0; }
        button { padding: 10px; background: #4CAF50; color: white; border: none; cursor: pointer; }
        button:hover { background: #45a049; }
    </style>
</head>
<body>
    <h1>Niagara - Real-time Data Stream</h1>
    <div id="status">Status: Disconnected</div>
    <div id="messages"></div>
    <button id="connect">Connect</button>
    <button id="disconnect" disabled>Disconnect</button>

    <script>
        const messagesDiv = document.getElementById('messages');
        const statusDiv = document.getElementById('status');
        const connectBtn = document.getElementById('connect');
        const disconnectBtn = document.getElementById('disconnect');
        let ws;

        function connect() {
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const wsUrl = protocol + '//' + window.location.host + '/connect';
            
            ws = new WebSocket(wsUrl);
            
            ws.onopen = function() {
                statusDiv.textContent = 'Status: Connected';
                connectBtn.disabled = true;
                disconnectBtn.disabled = false;
                addMessage('Connected to server');
            };
            
            ws.onmessage = function(event) {
                addMessage('Received: ' + event.data);
            };
            
            ws.onclose = function() {
                statusDiv.textContent = 'Status: Disconnected';
                connectBtn.disabled = false;
                disconnectBtn.disabled = true;
                addMessage('Disconnected from server');
            };
            
            ws.onerror = function(error) {
                addMessage('Error: ' + error);
                statusDiv.textContent = 'Status: Error';
            };
        }
        
        function disconnect() {
            if (ws) {
                ws.close();
            }
        }
        
        function addMessage(message) {
            const messageElem = document.createElement('div');
            messageElem.className = 'message';
            messageElem.textContent = new Date().toLocaleTimeString() + ' - ' + message;
            messagesDiv.appendChild(messageElem);
            messagesDiv.scrollTop = messagesDiv.scrollHeight;
        }
        
        connectBtn.addEventListener('click', connect);
        disconnectBtn.addEventListener('click', disconnect);
    </script>
</body>
</html>`;
		
		err := os.WriteFile("./static/index.html", []byte(indexHTML), 0644)
		if err != nil {
			log.Fatalf("Error creating index.html: %v", err)
		}
	}
	
	// Configure server with workers
	address := fmt.Sprintf("%s:%d", host, port)
	log.Printf("Starting Niagara server at %s", address)
	log.Printf("Configuration - Rate: %d messages/min, Workers: %d, Random Rate: %v", 
		rate, workers, randomRate)
	
	// Start server
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
