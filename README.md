# Load Test CLI Program

This is a powerful and flexible load testing CLI program implemented in Golang. It allows you to simulate concurrent users, customize test duration and request rate, and provides real-time monitoring and detailed result reporting.

## Project Structure

1. `main.go`: Entry point of the application, handles CLI setup and execution
2. `config/config.go`: Manages configuration and command-line flags
3. `loadgen/loadgen.go`: Implements the load generator and request handling
4. `monitor/monitor.go`: Handles real-time monitoring and statistics
5. `results/collector.go`: Collects and aggregates test results
6. `scenarios/scenarios.go`: Manages customizable test scenarios

## Installation

1. Ensure you have Go installed on your system
2. Clone the repository and navigate to the project directory
3. Run `go mod tidy` to download dependencies
4. Build the program: `go build -o loadtest`

## Usage

Run the program with desired flags, e.g.:

```
./loadtest --url https://example.com --users 10 --duration 1m --rate 100
```

### Using Custom Scenarios

1. Create a JSON file with your scenarios (e.g., scenarios.json)
2. Run the program with the scenario file:

```
./loadtest --scenario-file scenarios.json --users 10 --duration 1m
```

## Features

- Concurrent user simulation
- Customizable test duration and request rate
- HTTP method and custom headers support
- Real-time monitoring and detailed result reporting
- Customizable test scenarios
- Warm-up period
- Basic response validation

## Future Enhancements

- More advanced response validation
- Detailed resource monitoring (CPU, memory usage)
- Support for more protocols (e.g., WebSocket, gRPC)
- Distributed load testing capabilities
- Add hist plot
- Data Visualization


## License

This project is licensed under the MIT License - see the LICENSE file for details.