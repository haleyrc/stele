---
title: "Building REST APIs with Go"
description: "Learn how to build scalable and maintainable REST APIs using Go"
date: 2025-09-18T14:30:00Z
tags: ["go", "api", "web", "rest"]
---

# Crafting Robust REST APIs in Go

Building web APIs in Go combines the language's performance benefits with excellent HTTP handling capabilities. This guide walks through creating production-ready REST services.

## Foundation: HTTP Server Setup

Go's standard library provides powerful HTTP tools:

```go
package main

import (
    "net/http"
    "log"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/api/users", handleUsers)
    
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

## Structuring Your API

Organization matters for maintainable APIs:

### Handler Functions

Keep handlers focused and testable:

```go
func handleUsers(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        getUserList(w, r)
    case http.MethodPost:
        createUser(w, r)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}
```

### JSON Processing

Handle request and response data efficiently:

```go
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func createUser(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    // Process user creation...
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}
```

## Best Practices

### Error Handling

Implement consistent error responses:

```go
type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

func writeError(w http.ResponseWriter, code int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(APIError{Code: code, Message: message})
}
```

### Middleware Integration

Add cross-cutting concerns like logging and authentication:

```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}
```

## Testing Your API

Go's testing tools make API validation straightforward:

```go
func TestUserAPI(t *testing.T) {
    req := httptest.NewRequest("GET", "/api/users", nil)
    w := httptest.NewRecorder()
    
    handleUsers(w, req)
    
    if w.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", w.Code)
    }
}
```

Building REST APIs in Go leverages the language's strengths while maintaining code clarity and performance.