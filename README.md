Perfect — you want **a complete knowledge map for Go (Golang)** 👏

Below is a **comprehensive roadmap** covering **everything from beginner to advanced** — including **language syntax, concurrency, design patterns, performance tuning, and real-world applications**.
This is the **“full knowledge base”** for mastering Go professionally.

---

# 🧭 GO (GOLANG) FULL KNOWLEDGE ROADMAP

---

## 🪜 1. **Fundamentals**

### 🔤 Syntax & Basics

* Go environment setup (`go install`, `GOROOT`, `GOPATH`)
* File structure (`main.go`, `package main`)
* `main()` function entry point
* Variables & constants (`var`, `:=`, `const`)
* Data types: `int`, `float64`, `string`, `bool`
* Type inference and zero values
* Control flow:

  * `if / else`
  * `switch / case`
  * `for` (Go’s only loop)

### 🧱 Functions

* Basic function declaration
* Named return values
* Variadic functions
* Anonymous functions & closures
* Defer statements (for cleanup)
* Panic & recover (error-like control flow)

---

## 🧩 2. **Core Language Features**

### 📦 Packages & Modules

* `go mod init`, `go mod tidy`
* Importing and aliasing packages
* Creating reusable modules
* `go get` for dependencies

### 🧰 Structs & Interfaces

* Struct declaration & initialization
* Embedding & composition
* Methods with receiver types
* Interface definition (`type Shape interface { Area() float64 }`)
* Implicit interface implementation (no `implements` keyword)
* Empty interface (`interface{}`) and type assertions

### 🧶 Pointers

* Declaration (`var p *int`)
* Dereferencing (`*p`)
* Pass by reference
* Avoiding nil pointer dereference

---

## ⚙️ 3. **Data Structures**

### 🗃️ Built-in Types

* Arrays (`[5]int`)
* Slices (`[]int`, `make`, `append`, `copy`)
* Maps (`map[string]int`)
* Strings and runes (`rune`, `byte`)

### 📊 Custom Structures

* Structs with tags (e.g. JSON: `` `json:"name"` ``)
* Nested structs
* Embedding for inheritance-like behavior

### 🧠 Generics (Go 1.18+)

* Type parameters: `func PrintSlice[T any](s []T)`
* Constraints & interfaces
* Generic data structures and utilities

---

## 🔄 4. **Error Handling**

* `error` interface (`Error() string`)
* `fmt.Errorf`, `errors.New`, `errors.Is`, `errors.As`
* Wrapping errors (`%w`)
* Custom error types
* Panic vs. error — when to use which

---

## ⚡ 5. **Concurrency (Go’s Superpower)**

### 🧵 Goroutines

* `go func() { ... }()`
* Lightweight threads
* Synchronization

### 📬 Channels

* Unbuffered & buffered channels
* Send (`ch <- x`) / receive (`x := <-ch`)
* Channel directions (`chan<-`, `<-chan`)
* Closing channels (`close(ch)`)
* Select statements for multiplexing

### 🧭 Sync Primitives

* `sync.Mutex`, `sync.RWMutex`
* `sync.WaitGroup`
* `sync.Once`
* `sync.Cond`
* `context.Context` for cancellation and timeouts

---

## 💾 6. **I/O, File, and Network**

### 📄 Files & I/O

* Reading/writing files with `os` and `io`
* Buffered I/O (`bufio`)
* JSON/YAML/XML encoding (`encoding/json`)

### 🌐 Networking

* `net/http` for REST APIs
* HTTP client & server
* WebSocket (`gorilla/websocket`)
* TCP/UDP sockets (`net` package)

### 🔌 Socket Programming (Chi tiết)

Socket programming cho phép giao tiếp trực tiếp giữa các process qua mạng, không qua HTTP. Go cung cấp package `net` mạnh mẽ cho TCP, UDP và Unix Domain Sockets.

#### **TCP Socket Server**

TCP (Transmission Control Protocol) đảm bảo dữ liệu được truyền đáng tin cậy và có thứ tự.

**Example: TCP Echo Server**

```go
package main

import (
    "fmt"
    "io"
    "net"
    "strings"
)

func handleConnection(conn net.Conn) {
    defer conn.Close()
    
    // Đọc dữ liệu từ client
    buffer := make([]byte, 1024)
    for {
        n, err := conn.Read(buffer)
        if err != nil {
            if err == io.EOF {
                fmt.Println("Client disconnected")
                return
            }
            fmt.Printf("Read error: %v\n", err)
            return
        }
        
        message := strings.TrimSpace(string(buffer[:n]))
        fmt.Printf("Received: %s from %s\n", message, conn.RemoteAddr())
        
        // Echo lại cho client
        if message == "quit" {
            conn.Write([]byte("Goodbye!\n"))
            return
        }
        
        response := fmt.Sprintf("Echo: %s\n", message)
        conn.Write([]byte(response))
    }
}

func main() {
    // Lắng nghe trên port 8080
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        fmt.Printf("Error listening: %v\n", err)
        return
    }
    defer listener.Close()
    
    fmt.Println("TCP Server listening on :8080")
    
    // Chấp nhận kết nối mới
    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Printf("Accept error: %v\n", err)
            continue
        }
        
        // Xử lý mỗi connection trong một goroutine riêng
        go handleConnection(conn)
    }
}
```

#### **TCP Socket Client**

```go
package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
)

func main() {
    // Kết nối đến server
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        fmt.Printf("Error connecting: %v\n", err)
        return
    }
    defer conn.Close()
    
    fmt.Println("Connected to server. Type 'quit' to exit.")
    
    // Đọc input từ user
    reader := bufio.NewReader(os.Stdin)
    
    // Goroutine để đọc phản hồi từ server
    go func() {
        for {
            response := make([]byte, 1024)
            n, err := conn.Read(response)
            if err != nil {
                fmt.Printf("Read error: %v\n", err)
                return
            }
            fmt.Printf("Server: %s", string(response[:n]))
        }
    }()
    
    // Gửi messages đến server
    for {
        fmt.Print("You: ")
        text, _ := reader.ReadString('\n')
        text = strings.TrimSpace(text)
        
        if text == "quit" {
            conn.Write([]byte("quit\n"))
            break
        }
        
        _, err := conn.Write([]byte(text + "\n"))
        if err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }
    }
}
```

#### **UDP Socket Server**

UDP (User Datagram Protocol) không đảm bảo thứ tự và độ tin cậy nhưng nhanh hơn, phù hợp cho real-time applications.

```go
package main

import (
    "fmt"
    "net"
)

func main() {
    // Lắng nghe UDP trên port 8081
    addr, err := net.ResolveUDPAddr("udp", ":8081")
    if err != nil {
        fmt.Printf("Error resolving address: %v\n", err)
        return
    }
    
    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        fmt.Printf("Error listening: %v\n", err)
        return
    }
    defer conn.Close()
    
    fmt.Println("UDP Server listening on :8081")
    
    buffer := make([]byte, 1024)
    
    for {
        // Đọc dữ liệu từ client
        n, clientAddr, err := conn.ReadFromUDP(buffer)
        if err != nil {
            fmt.Printf("Read error: %v\n", err)
            continue
        }
        
        message := string(buffer[:n])
        fmt.Printf("Received from %s: %s\n", clientAddr, message)
        
        // Gửi phản hồi
        response := fmt.Sprintf("Echo: %s", message)
        conn.WriteToUDP([]byte(response), clientAddr)
    }
}
```

#### **UDP Socket Client**

```go
package main

import (
    "fmt"
    "net"
    "time"
)

func main() {
    // Kết nối đến UDP server
    serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8081")
    if err != nil {
        fmt.Printf("Error resolving server address: %v\n", err)
        return
    }
    
    conn, err := net.DialUDP("udp", nil, serverAddr)
    if err != nil {
        fmt.Printf("Error dialing: %v\n", err)
        return
    }
    defer conn.Close()
    
    // Gửi messages
    messages := []string{"Hello", "World", "From UDP Client"}
    
    for _, msg := range messages {
        _, err := conn.Write([]byte(msg))
        if err != nil {
            fmt.Printf("Write error: %v\n", err)
            continue
        }
        
        // Đọc phản hồi
        buffer := make([]byte, 1024)
        conn.SetReadDeadline(time.Now().Add(1 * time.Second))
        n, err := conn.Read(buffer)
        if err != nil {
            fmt.Printf("Read error: %v\n", err)
            continue
        }
        
        fmt.Printf("Server response: %s\n", string(buffer[:n]))
        time.Sleep(500 * time.Millisecond)
    }
}
```

#### **Unix Domain Sockets (Local IPC)**

Unix Domain Sockets dùng cho giao tiếp giữa các process trên cùng máy, nhanh hơn TCP/UDP.

```go
package main

import (
    "fmt"
    "io"
    "net"
    "os"
    "os/signal"
    "syscall"
)

func handleUnixConnection(conn net.Conn) {
    defer conn.Close()
    
    buffer := make([]byte, 1024)
    for {
        n, err := conn.Read(buffer)
        if err != nil {
            if err == io.EOF {
                return
            }
            fmt.Printf("Read error: %v\n", err)
            return
        }
        
        message := string(buffer[:n])
        fmt.Printf("Received: %s\n", message)
        
        conn.Write([]byte("OK\n"))
    }
}

func main() {
    socketPath := "/tmp/go_socket.sock"
    
    // Xóa socket file nếu tồn tại
    os.Remove(socketPath)
    
    listener, err := net.Listen("unix", socketPath)
    if err != nil {
        fmt.Printf("Error listening: %v\n", err)
        return
    }
    defer listener.Close()
    defer os.Remove(socketPath)
    
    // Xử lý tín hiệu để cleanup
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    
    go func() {
        <-sigChan
        os.Remove(socketPath)
        os.Exit(0)
    }()
    
    fmt.Println("Unix Socket Server listening on", socketPath)
    
    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Printf("Accept error: %v\n", err)
            continue
        }
        
        go handleUnixConnection(conn)
    }
}
```

#### **Advanced Socket Features**

**1. Socket Options & Timeouts**

```go
package main

import (
    "fmt"
    "net"
    "time"
)

func setSocketOptions(conn net.Conn) {
    // Set read deadline
    conn.SetReadDeadline(time.Now().Add(10 * time.Second))
    
    // Set write deadline
    conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
    
    // Set keep-alive
    if tcpConn, ok := conn.(*net.TCPConn); ok {
        tcpConn.SetKeepAlive(true)
        tcpConn.SetKeepAlivePeriod(30 * time.Second)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()
    setSocketOptions(conn)
    
    buffer := make([]byte, 1024)
    n, err := conn.Read(buffer)
    if err != nil {
        fmt.Printf("Read error: %v\n", err)
        return
    }
    
    fmt.Printf("Received: %s\n", string(buffer[:n]))
    conn.Write([]byte("Response\n"))
}

func main() {
    listener, err := net.Listen("tcp", ":8082")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    defer listener.Close()
    
    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        go handleConnection(conn)
    }
}
```

**2. Concurrent Socket Pool với Worker Pool Pattern**

```go
package main

import (
    "fmt"
    "net"
    "sync"
)

type SocketPool struct {
    workers    int
    connections chan net.Conn
    wg          sync.WaitGroup
}

func NewSocketPool(workers int) *SocketPool {
    return &SocketPool{
        workers:    workers,
        connections: make(chan net.Conn, 100),
    }
}

func (sp *SocketPool) worker() {
    defer sp.wg.Done()
    
    for conn := range sp.connections {
        // Xử lý connection
        buffer := make([]byte, 1024)
        n, err := conn.Read(buffer)
        if err != nil {
            conn.Close()
            continue
        }
        
        fmt.Printf("Worker processed: %s\n", string(buffer[:n]))
        conn.Write([]byte("Processed\n"))
        conn.Close()
    }
}

func (sp *SocketPool) Start() {
    for i := 0; i < sp.workers; i++ {
        sp.wg.Add(1)
        go sp.worker()
    }
}

func (sp *SocketPool) Stop() {
    close(sp.connections)
    sp.wg.Wait()
}

func main() {
    pool := NewSocketPool(5)
    pool.Start()
    defer pool.Stop()
    
    listener, _ := net.Listen("tcp", ":8083")
    defer listener.Close()
    
    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        pool.connections <- conn
    }
}
```

**3. Socket với Context cho Graceful Shutdown**

```go
package main

import (
    "context"
    "fmt"
    "net"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    // Handle shutdown signals
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    
    go func() {
        <-sigChan
        fmt.Println("Shutting down...")
        cancel()
    }()
    
    listener, err := net.Listen("tcp", ":8084")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    defer listener.Close()
    
    // Accept connections trong goroutine riêng
    connChan := make(chan net.Conn)
    errChan := make(chan error)
    
    go func() {
        for {
            conn, err := listener.Accept()
            if err != nil {
                errChan <- err
                return
            }
            connChan <- conn
        }
    }()
    
    for {
        select {
        case <-ctx.Done():
            fmt.Println("Server stopped")
            return
        case conn := <-connChan:
            go func(c net.Conn) {
                defer c.Close()
                c.SetDeadline(time.Now().Add(5 * time.Second))
                
                buffer := make([]byte, 1024)
                n, _ := c.Read(buffer)
                fmt.Printf("Received: %s\n", string(buffer[:n]))
                c.Write([]byte("OK\n"))
            }(conn)
        case err := <-errChan:
            fmt.Printf("Accept error: %v\n", err)
        }
    }
}
```

#### **Các Khái Niệm Quan Trọng**

* **TCP vs UDP:**
  * TCP: Đáng tin cậy, có thứ tự, connection-oriented
  * UDP: Nhanh, không đảm bảo, connectionless
  
* **Buffering:**
  * TCP: Stream-based, cần buffer để xử lý messages
  * UDP: Datagram-based, mỗi `ReadFromUDP` là một message hoàn chỉnh

* **Concurrency:**
  * Mỗi connection nên được xử lý trong goroutine riêng
  * Sử dụng channels để quản lý connections
  * Worker pool pattern cho high-performance applications

* **Error Handling:**
  * Luôn kiểm tra `io.EOF` để detect client disconnect
  * Sử dụng deadlines để tránh blocking indefinitely
  * Graceful shutdown với context và signals

* **Performance Tips:**
  * Reuse buffers khi có thể
  * Set appropriate read/write deadlines
  * Sử dụng connection pooling cho clients
  * Monitor goroutine count để tránh resource exhaustion

---

## 🧱 7. **Testing & Tooling**

### ✅ Testing

* Unit testing with `testing` package
* Table-driven tests
* Benchmarks (`go test -bench`)
* Test coverage (`go test -cover`)
* Mocks (`testify/mock`)

### 🧰 Tooling

* `go fmt` (format code)
* `go vet` (linting)
* `go mod` (dependency management)
* `go build`, `go install`, `go run`
* Profiling: `pprof`, `trace`
* Static analysis: `golangci-lint`

---

## 🧩 8. **Web Development in Go**

* `net/http` basics
* Routers: `gorilla/mux`, `chi`, `gin-gonic/gin`
* Middleware (logging, CORS, authentication)
* Templates (`html/template`)
* Serving static files
* JSON APIs
* RESTful conventions

---

## 🧠 9. **Advanced Topics**

### 🏗️ Design Patterns in Go

* Singleton (via `sync.Once`)
* Factory
* Adapter
* Strategy
* Observer (via channels)
* Dependency injection (manual or via packages like `wire`)

### 🧬 Reflection

* `reflect.TypeOf`, `reflect.ValueOf`
* Dynamic field access
* Tag reading (`reflect.StructTag`)

### 🧮 Memory & Performance

* Stack vs Heap allocation
* Escape analysis
* Garbage collector behavior
* Profiling CPU and memory

### 🧳 Generics in Production

* Type constraints
* Generic repository or utility pattern

---

## 🧰 10. **DevOps & Deployment**

* Build binaries for multiple platforms (`GOOS`, `GOARCH`)
* Environment configuration
* Dockerizing Go applications
* Continuous Integration (GitHub Actions, GitLab CI)
* Logging and monitoring (`zap`, `logrus`)
* Graceful shutdowns with `context`

---

## 🌍 11. **Ecosystem & Frameworks**

| Purpose              | Library / Framework  |
| -------------------- | -------------------- |
| Web Framework        | Gin, Echo, Fiber     |
| ORM / DB             | GORM, sqlx, ent      |
| Testing              | Testify, Ginkgo      |
| Logging              | Logrus, Zap, Zerolog |
| Dependency Injection | Wire, Fx             |
| CLI Tools            | Cobra, Urfave/cli    |
| Microservices        | Go-Kit, gRPC, NATS   |

---

## 🚀 12. **Practical Projects to Build**

1. **CLI App** – e.g., file organizer or todo manager
2. **REST API** – CRUD + JWT authentication
3. **Concurrent Scraper** – using goroutines & channels
4. **WebSocket Chat Server** – real-time communication
5. **gRPC Microservice** – structured communication
6. **Go + Docker** – containerized REST API
7. **CI/CD Pipeline** – automated build & test

---

## 🧾 13. **Best Practices**

* Prefer composition over inheritance
* Use channels over shared memory
* Keep packages small and focused
* Handle errors explicitly
* Avoid global variables
* Document with `godoc` comments

