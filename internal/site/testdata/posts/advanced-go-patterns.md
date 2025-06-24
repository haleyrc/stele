---
title: "Advanced Go Patterns"
description: "Exploring advanced design patterns and best practices in Go"
date: 2025-09-15T09:15:00Z
tags: ["go", "patterns", "advanced", "best-practices"]
---

# Mastering Advanced Go Design Patterns

As Go applications grow in complexity, employing sophisticated design patterns becomes essential for maintainable and scalable code.

## Functional Options Pattern

This pattern provides flexible configuration for complex types:

```go
type Server struct {
    host    string
    port    int
    timeout time.Duration
}

type Option func(*Server)

func WithHost(host string) Option {
    return func(s *Server) {
        s.host = host
    }
}

func WithPort(port int) Option {
    return func(s *Server) {
        s.port = port
    }
}

func NewServer(options ...Option) *Server {
    s := &Server{
        host:    "localhost",
        port:    8080,
        timeout: 30 * time.Second,
    }
    
    for _, option := range options {
        option(s)
    }
    
    return s
}
```

Usage becomes intuitive and extensible:

```go
server := NewServer(
    WithHost("api.example.com"),
    WithPort(9000),
)
```

## Worker Pool Pattern

Manage concurrent processing efficiently:

```go
type WorkerPool struct {
    workers   int
    jobs      chan Job
    results   chan Result
    done      chan bool
}

func NewWorkerPool(workers int) *WorkerPool {
    return &WorkerPool{
        workers: workers,
        jobs:    make(chan Job, 100),
        results: make(chan Result, 100),
        done:    make(chan bool),
    }
}

func (wp *WorkerPool) Start() {
    for i := 0; i < wp.workers; i++ {
        go wp.worker()
    }
}

func (wp *WorkerPool) worker() {
    for job := range wp.jobs {
        result := processJob(job)
        wp.results <- result
    }
}
```

## Interface Segregation

Design focused interfaces for better testability:

```go
// Instead of one large interface
type Storage interface {
    Save(data []byte) error
    Load(id string) ([]byte, error)
    Delete(id string) error
    List() ([]string, error)
    Backup() error
    Restore(path string) error
}

// Prefer smaller, focused interfaces
type Reader interface {
    Load(id string) ([]byte, error)
}

type Writer interface {
    Save(data []byte) error
}

type Deleter interface {
    Delete(id string) error
}
```

## Context-Driven Design

Leverage context for cancellation and timeouts:

```go
func ProcessWithTimeout(ctx context.Context, data []byte) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    resultChan := make(chan error, 1)
    
    go func() {
        resultChan <- heavyProcessing(data)
    }()
    
    select {
    case err := <-resultChan:
        return err
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

## Builder Pattern for Complex Objects

Construct intricate objects step by step:

```go
type QueryBuilder struct {
    table   string
    fields  []string
    where   []string
    orderBy string
    limit   int
}

func NewQuery() *QueryBuilder {
    return &QueryBuilder{}
}

func (qb *QueryBuilder) From(table string) *QueryBuilder {
    qb.table = table
    return qb
}

func (qb *QueryBuilder) Select(fields ...string) *QueryBuilder {
    qb.fields = append(qb.fields, fields...)
    return qb
}

func (qb *QueryBuilder) Where(condition string) *QueryBuilder {
    qb.where = append(qb.where, condition)
    return qb
}

func (qb *QueryBuilder) Build() string {
    // Construct SQL query from builder state
    return buildQueryString(qb)
}
```

## Error Wrapping Strategy

Create meaningful error chains:

```go
type ProcessingError struct {
    Op   string
    Path string
    Err  error
}

func (e *ProcessingError) Error() string {
    return fmt.Sprintf("processing %s failed at %s: %v", e.Op, e.Path, e.Err)
}

func (e *ProcessingError) Unwrap() error {
    return e.Err
}

func processFile(path string) error {
    if err := validatePath(path); err != nil {
        return &ProcessingError{
            Op:   "validation",
            Path: path,
            Err:  err,
        }
    }
    
    // Continue processing...
    return nil
}
```

These patterns elevate Go code from functional to elegant, providing foundations for complex, maintainable applications.