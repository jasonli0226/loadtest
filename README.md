# Load Test CLI Program

This is a powerful and flexible load testing CLI program implemented in Golang. It allows you to simulate concurrent users, customize test duration and request rate, and provides real-time monitoring and detailed result reporting.

## Installation

1. Ensure you have Go installed on your system
2. Clone the repository and navigate to the project directory
3. Run `go mod tidy` to download dependencies
4. Build the program: `go build -o loadtest`

## Usage

Run the program with desired flags, e.g.:

```
loadtest --url https://example.com --headers 'token=abc' \
    --users 10 --duration 1m --rate 10
```

Result example:

```bash
Starting load test...

Load Test Results:
Duration: 1m6.537116126s

Test Duration: 1m0.389113036s
Total Requests: 600
Successful Requests: 600
Failed Requests: 0
Requests per second: 9.94
Average Latency: 180.015439ms

Response Time Statistics:
Min: 156.00 ms
Max: 488.00 ms
Mean: 179.48 ms
Median: 169.00 ms
50th percentile: 169.00 ms
95th percentile: 253.00 ms
99th percentile: 337.00 ms

Response Time Histogram:
  0.00 ms -  16.60 ms | ################################################## (386)
 16.60 ms -  33.20 ms | ################# (136)
 33.20 ms -  49.80 ms | ### (24)
 49.80 ms -  66.40 ms | # (8)
 66.40 ms -  83.00 ms | # (9)
 83.00 ms -  99.60 ms | # (10)
 99.60 ms - 116.20 ms |  (6)
116.20 ms - 132.80 ms |  (6)
132.80 ms - 149.40 ms |  (3)
149.40 ms - 166.00 ms |  (2)
166.00 ms - 182.60 ms |  (4)
182.60 ms - 199.20 ms |  (0)
199.20 ms - 215.80 ms |  (3)
215.80 ms - 232.40 ms |  (2)
232.40 ms - 249.00 ms |  (0)
249.00 ms - 265.60 ms |  (0)
265.60 ms - 282.20 ms |  (0)
282.20 ms - 298.80 ms |  (0)
298.80 ms - 315.40 ms |  (0)
315.40 ms - 332.00 ms |  (1)


Status Codes:
200: 600

Total Errors: 0
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