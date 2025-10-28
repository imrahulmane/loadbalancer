# Load Balancer in Go

A production-ready HTTP load balancer built from scratch in Go, implementing round-robin distribution and health checking across multiple backend servers.

## Project Overview

This project demonstrates building a Layer 7 (HTTP) reverse proxy load balancer without using any frameworks. Built to understand distributed systems concepts, concurrent programming, and network protocols at a fundamental level.

**Built with:** Pure Go (no external dependencies except standard library)

---

## Features

### Level 1: Basic Load Balancing

- ✅ TCP connection handling
- ✅ Round-robin load balancing algorithm
- ✅ Concurrent request handling (goroutine per connection)
- ✅ Bidirectional data streaming with `io.Copy`
- ✅ Error handling (502 Bad Gateway for failed backends)
- ✅ Thread-safe round-robin state management

### Level 2: Health Checking

- ✅ Background health checker (checks every 10 seconds)
- ✅ Automatic unhealthy server detection
- ✅ Automatic recovery detection
- ✅ Smart round-robin (skips unhealthy servers)
- ✅ Thread-safe health status tracking (RWMutex)
- ✅ Graceful handling when all backends are down

---

## Architecture

```
                                  ┌─────────────────┐
                                  │  Load Balancer  │
                                  │   Port 8090     │
                                  └────────┬────────┘
                                           │
                    ┌──────────────────────┼──────────────────────┐
                    │                      │                      │
                    ▼                      ▼                      ▼
            ┌──────────────┐      ┌──────────────┐      ┌──────────────┐
            │  Backend 1   │      │  Backend 2   │      │  Backend 3   │
            │  Port 9001   │      │  Port 9002   │      │  Port 9003   │
            └──────────────┘      └──────────────┘      └──────────────┘
```

---

## How It Works

### Round Robin Algorithm

Distributes requests evenly across healthy backends:

```
Backends: [A, B, C]
Requests: 1→A, 2→B, 3→C, 4→A, 5→B, 6→C...

If B becomes unhealthy:
Requests: 1→A, 2→C, 3→A, 4→C...
```

**Implementation:**

```go
currentIndex = (currentIndex + 1) % len(servers)
```

### Health Checking

**Active health checks** run every 10 seconds:

1. Dial each backend with 2-second timeout
2. If successful → Mark healthy
3. If failed → Mark unhealthy
4. Log status changes only

**Smart routing** skips unhealthy servers automatically.

### Bidirectional Data Flow

```go
// Goroutine: Client → Backend (concurrent)
go io.Copy(backendConn, clientConn)

// Main goroutine: Backend → Client (blocks to keep function alive)
io.Copy(clientConn, backendConn)
```

**Why this pattern?**

- Streams data without buffering entire request/response
- Constant memory usage (works for 1KB or 1GB requests)
- Main goroutine blocks until response completes
- Keeps connections alive until data transfer done

---

## Project Structure

```
loadbalancer/
├── go.mod
├── main.go              # Entry point, configuration
|__ backend-servers      # Server for testing
    ├── server1.js      # Test backend server 1
    ├── server2.js      # Test backend server 2
    ├── server3.        # Test backend server 3
└── balancer/
    ├── balancer.go      # Core load balancer logic
    └── handler.go       # Connection handling
```

### Package Organization

**balancer/balancer.go:**

- LoadBalancer struct definition
- Round-robin logic (`getNextServer()`)
- Health checking (`checkHealth()`, `startHealthChecker()`)
- Server startup (`Start()`)

**balancer/handler.go:**

- Connection handling (`handleConnection()`)
- Bidirectional data copying
- Error responses (502 Bad Gateway)

**main.go:**

- Backend server configuration
- LoadBalancer initialization
- Server startup

---

## Installation & Setup

### Prerequisites

- Go 1.16+ installed
- Node.js installed (for test backend servers)

### Clone/Setup

```bash
mkdir loadbalancer
cd loadbalancer
go mod init loadbalancer
```

### Project Files

Create the directory structure and copy the Go files and Node.js backend servers.

---

## Running the Load Balancer

### Start Backend Servers

**Terminal 1:**

```bash
node backend1.js
# Output: Backend Server 1 listening on port 9001
```

**Terminal 2:**

```bash
node backend2.js
# Output: Backend Server 2 listening on port 9002
```

**Terminal 3:**

```bash
node backend3.js
# Output: Backend Server 3 listening on port 9003
```

### Start Load Balancer

**Terminal 4:**

```bash
go run main.go
# Or build and run:
go build
./loadbalancer
```

**Expected output:**

```
Starting load balancer...
Load balancer listening on :8090
Forwarding to backends: [localhost:9001 localhost:9002 localhost:9003]
Health checker started (checking every 10 seconds)
Running health checks...
Backend localhost:9001 marked as HEALTHY
Backend localhost:9002 marked as HEALTHY
Backend localhost:9003 marked as HEALTHY
```

---

## Testing

### Test 1: Basic Round Robin

```bash
for i in {1..9}; do curl http://localhost:8090; done
```

**Expected output:**

```
Response from Backend Server 1
Response from Backend Server 2
Response from Backend Server 3
Response from Backend Server 1
Response from Backend Server 2
Response from Backend Server 3
Response from Backend Server 1
Response from Backend Server 2
Response from Backend Server 3
```

Perfect round-robin distribution!

---

### Test 2: Backend Failure & Recovery

**Kill Backend 2:**

```bash
# In Terminal 2, press Ctrl+C
```

**Wait 10 seconds for health check, then:**

```bash
for i in {1..6}; do curl http://localhost:8090; done
```

**Expected output:**

```
Response from Backend Server 1
Response from Backend Server 3
Response from Backend Server 1
Response from Backend Server 3
Response from Backend Server 1
Response from Backend Server 3
```

Backend 2 is automatically skipped! No 502 errors.

**Restart Backend 2:**

```bash
node backend2.js
```

**Wait 10 seconds, then:**

```bash
for i in {1..6}; do curl http://localhost:8090; done
```

Backend 2 is back in rotation!

---

### Test 3: All Backends Down

**Stop all backends (Ctrl+C in all terminals).**

**Try a request:**

```bash
curl http://localhost:8090
```

**Expected output:**

```
Backend Unavailable
```

**Load balancer logs:**

```
All backends are unhealthy!
```

Graceful degradation - returns 502 instead of crashing.

---

### Test 4: Concurrent Load

```bash
# Send 100 concurrent requests
for i in {1..100}; do curl http://localhost:8090 & done
wait
```

All requests handled concurrently - one goroutine per connection.

---

## Key Concepts

### 1. Reverse Proxy Pattern

- Load balancer is a **middleman**, not the destination
- Forwards requests to backends, forwards responses to clients
- Manages two connections simultaneously per request

### 2. Concurrency in Go

- Goroutines for concurrent connection handling
- One goroutine per client connection
- Background goroutine for health checking
- Lightweight and efficient (thousands of connections possible)

### 3. Thread Safety

- **Mutex** for protecting `currentIndex` (simple mutual exclusion)
- **RWMutex** for protecting `healthy` map (read-heavy workload)
  - Multiple readers can proceed simultaneously
  - Writers get exclusive access
  - Prevents concurrent map read/write crashes

### 4. Bidirectional Data Streaming

- `io.Copy()` for efficient data transfer
- Goroutine for client → backend
- Main goroutine for backend → client (blocks to keep function alive)
- Constant memory usage regardless of request/response size

### 5. Background Tasks

- `time.Ticker` for periodic health checks
- Runs independently of request handling
- Non-blocking, efficient pattern for scheduled tasks

### 6. Network Programming

- `net.Listen()` - Server side (waiting for connections)
- `net.Dial()` - Client side (creating connections)
- `net.DialTimeout()` - Client with timeout (for health checks)
- Load balancer is both server (to clients) and client (to backends)

### 7. Health Checking Strategies

- **Active health checks:** Proactively ping backends
- Timeout-based detection (2 seconds)
- Status change logging (only log transitions)
- Automatic recovery detection

---

## Limitations & Trade-offs

### Current Limitations

- No connection pooling (creates new connection per request)
- No persistent connections (HTTP/1.1 keep-alive)
- Health checks run serially (not parallelized)
- 10-second detection window (failed backends serve traffic for up to 10s)
- No weighted round-robin
- No least-connections algorithm
- No SSL/TLS termination
- No request logging or metrics

### Design Decisions

**Why Round Robin?**

- Simple and predictable
- Fair distribution
- Good enough for equal-capacity backends

**Why 10-second health checks?**

- Balance between detection speed and overhead
- Too frequent: Wastes resources
- Too infrequent: Slow failure detection

**Why single health checker goroutine?**

- Simpler than parallel checking
- Adequate for dozens of backends
- Serialized checks easier to reason about

**Why RWMutex for health map?**

- Read-heavy workload (every request checks health)
- Multiple readers don't block each other
- Better performance than regular Mutex

---

## Future Enhancements

### Possible Improvements

- [ ] Connection pooling (reuse backend connections)
- [ ] Least-connections algorithm
- [ ] Weighted round-robin
- [ ] Parallel health checking
- [ ] Configurable health check interval
- [ ] Passive health checks (mark unhealthy on request failure)
- [ ] Exponential backoff for recovery
- [ ] HTTP/1.1 persistent connections
- [ ] Request logging and metrics
- [ ] Prometheus metrics export
- [ ] Configuration file (YAML/JSON)
- [ ] Graceful shutdown
- [ ] SSL/TLS support
- [ ] Path-based routing
- [ ] Sticky sessions

---

## Acknowledgments

Built as a learning project to understand:

- Low-level network programming
- Distributed systems concepts
- Production load balancer architecture
- Go concurrency patterns

---

## License

This project is for educational purposes.

---

## Contact

**Email:** rahulmane.dev@gmail.com

Feel free to reach out with questions, feedback, or discussions about distributed systems and Go!
