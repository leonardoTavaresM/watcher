# Watcher

A lightweight file system watcher written in Go that monitors directory changes and publishes file events in real-time.

## Overview

Watcher is a simple yet powerful file system monitoring tool that uses the `fsnotify` library to track file and directory changes. It follows clean architecture principles with a domain-driven design approach, making it easy to extend and maintain.

## Features

- **Real-time file monitoring**: Watches directories recursively for file changes
- **Event types**: Tracks CREATE, MODIFY, and REMOVE operations
- **Smart filtering**: Automatically ignores common directories like `node_modules`, `.git`, `vendor`, and `dist`
- **JSON output**: Publishes events in structured JSON format
- **Docker support**: Ready-to-use Docker container with volume mounting
- **Clean architecture**: Modular design with domain, service, and adapter layers

## Architecture

The project follows clean architecture principles:

```
├── cmd/api/           # Application entry point
├── internal/
│   ├── domain/        # Core business entities and interfaces
│   ├── service/       # Business logic layer
│   └── adapter/       # External dependencies (fsnotify, console)
```

### Components

- **Domain**: Defines the `FileEvent` struct and `EventPublisher` interface
- **Service**: Contains the business logic for handling file events
- **Adapters**: 
  - `fsnotify`: File system monitoring using fsnotify library
  - `consolepub`: Console-based event publisher (JSON output)

## Installation

### Prerequisites

- Go 1.24.6 or later
- Docker (optional)

### Local Development

1. Clone the repository:
```bash
git clone https://github.com/leonardoTavaresM/watcher.git
cd watcher
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build -o watcher ./cmd/api
```

4. Run the watcher:
```bash
# Set the directory to watch
export WATCH_PATH=/path/to/watch
./watcher
```

## Docker Usage

The project includes Docker support for easy deployment and testing.

### Build the Docker image:
```bash
make build
```

### Run with Docker:
```bash
make run
```

This will:
- Mount your local development directory (`/home/leonardomalt/Documents/dev`) to `/app/dev` in the container
- Set `WATCH_PATH=/app/dev` as the default watch directory
- Run the watcher in interactive mode

### Other Docker commands:
```bash
make build-nc    # Build without cache
make shell       # Enter container shell
make clean       # Remove Docker image
make rebuild     # Clean, build, and run
```

## Configuration

### Environment Variables

- `WATCH_PATH`: Directory path to monitor (default: `/app/dev`)

### Ignored Directories

The watcher automatically ignores the following directories:
- `node_modules`
- `.git`
- `vendor`
- `dist`

## Usage Examples

### Basic Usage
```bash
# Watch current directory
export WATCH_PATH=.
./watcher
```

### Docker Usage
```bash
# Watch your development directory
make run
```

### Sample Output
When files are created, modified, or removed, you'll see JSON output like:
```json
publish {"Timestamp":"2024-01-15T10:30:45.123Z","FilePath":"/app/dev/src/main.go","Ext":".go","Event":"CREATE"}
publish {"Timestamp":"2024-01-15T10:30:46.456Z","FilePath":"/app/dev/src/main.go","Ext":".go","Event":"MODIFY"}
```

## Development

### Project Structure
```
watcher/
├── cmd/api/main.go              # Application entry point
├── internal/
│   ├── domain/
│   │   └── fileevent.go         # Core domain entities
│   ├── service/watcher/
│   │   └── watcher.go           # Business logic
│   └── adapter/
│       ├── fsnotify/
│       │   ├── fsnotify.go      # File system adapter
│       │   └── utils.go         # Utility functions
│       └── consolepub/
│           └── consolepub.go    # Console publisher
├── dockerfile                   # Docker configuration
├── Makefile                     # Build automation
├── go.mod                       # Go module definition
└── go.sum                       # Dependency checksums
```

### Adding New Publishers

To add a new event publisher (e.g., webhook, message queue), implement the `EventPublisher` interface:

```go
type CustomPublisher struct{}

func (p *CustomPublisher) Publish(event domain.FileEvent) error {
    // Your custom publishing logic
    return nil
}
```

### Extending Event Types

To add new event types, modify the `fsnotify.go` adapter to handle additional `fsnotify` operations.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Author

Leonardo Tavares Malt - [@leonardoTavaresM](https://github.com/leonardoTavaresM)
