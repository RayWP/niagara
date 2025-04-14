# Niagara C++ Implementation

This is the C++ implementation of the Niagara real-time data streaming platform.

## Requirements

- CMake 3.10 or higher
- C++17 compatible compiler
- Git (to clone external dependencies)

## External Dependencies

- [websocketpp](https://github.com/zaphoyd/websocketpp) (header-only WebSocket library)
- [asio](https://github.com/chriskohlhoff/asio) (header-only networking library)

## Setup

### Clone External Dependencies

```bash
# Create external directory
mkdir -p external

# Clone websocketpp
git clone https://github.com/zaphoyd/websocketpp.git external/websocketpp

# Clone asio
git clone https://github.com/chriskohlhoff/asio.git external/asio
```

### Build Instructions

```bash
# Create build directory
mkdir build && cd build

# Configure with CMake
cmake ..

# Build
cmake --build .
```

## Run

```bash
./niagara
```

## Run Options

```bash
--host=string   # The host address of the machine (default: localhost)
--port=int      # The port number for the service (default: 8080)
--rate=int      # How many data messages per minute to be sent (default: 60)
--worker=int    # How many worker threads (default: 1)
```

## Run Example

```bash
./niagara --host=localhost --port=8080 --rate=10 --worker=2
```

Open a WebSocket client and connect to `ws://{host}:{port}` to receive real-time data.

## Data Format

The server sends stock market data in the following format:

```
StockCode|Price|B/S|Amount
```

For example: `AAPL|100|B|5`
