# Niagara - Golang Implementation

This is the Golang implementation of the Niagara real-time data streaming platform. It provides a WebSocket server that generates random stock data at configurable rates.

## Features

- WebSocket server for real-time data streaming
- Configurable data generation rate
- Performance monitoring of HTTP and WebSocket operations
- Chaos mode for randomly changing data rates
- Multiple worker support
- Web interface for visualizing the data stream
- Comprehensive logging of I/O activities with duration tracking

## Requirements

- Go 1.19 or higher
- gorilla/websocket package

## Installation

### Clone the Repository

```bash
git clone https://github.com/yourusername/niagara.git
cd niagara/golang
```

### Install Dependencies

```bash
go mod download
```

## Running the Server

### Basic Usage

```bash
go run .
```

This starts the server on `localhost:8080` with default settings.

### Command Line Arguments

The server supports several command line arguments:

- `--host`: Host address (default: "localhost")
- `--port`: Port number (default: 8080)
- `--rate`: Data generation rate per minute (default: 60)
- `--worker`: Number of worker processes (default: 1)
- `--random-rate`: Enable random rate changes (chaos mode) (default: false)

Example:

```bash
go run . --host=0.0.0.0 --port=9090 --rate=120 --worker=4 --random-rate=true
```

## Testing

Run the unit tests with:

```bash
go test -v
```

## Performance Monitoring

The server automatically logs performance statistics:

- HTTP operations should complete in less than 300ms
- WebSocket operations should complete in less than 300ms
- Database operations (if implemented) should complete in less than 100ms

Warnings are logged when operations exceed these thresholds.

## Accessing the Web Interface

Open a web browser and navigate to `http://localhost:8080` (or your configured host/port).
The web interface allows you to connect to the WebSocket server and view the real-time data stream.

## Data Format

The data is sent in the following format:

```text
StockCode|Price|B/S|Amount
```

Example:

```text
AAPL|100|B|42
```

Where:

- `StockCode` is the stock symbol (e.g., AAPL)
- `Price` is the stock price
- `B/S` indicates Buy (B) or Sell (S)
- `Amount` is the quantity

## Building for Production

To build a standalone binary:

```bash
go build -o niagara-server
```

Then run the server:

```bash
./niagara-server --host=0.0.0.0 --port=8080
```
