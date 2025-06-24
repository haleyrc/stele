---
title: "Exploring Go Generics"
description: "A deep dive into Go's generic programming features"
draft: true
tags: ["go", "generics", "advanced"]
---

# Understanding Go Generics

Go 1.18 introduced generics, fundamentally changing how we write reusable code. This feature allows you to write functions and data structures that work with multiple types while maintaining type safety.

## Why Generics Matter

Before generics, we had two unsatisfying options:

- **Interface{}** - Type-safe but requires type assertions
- **Code duplication** - Type-safe but repetitive and error-prone

Generics provide the best of both worlds.

## Basic Generic Function

Here's a simple example of a generic function that works with any comparable type:

```go
func Contains[T comparable](slice []T, item T) bool {
    for _, v := range slice {
        if v == item {
            return true
        }
    }
    return false
}
```

You can use this with any type that supports equality:

```go
numbers := []int{1, 2, 3, 4, 5}
hasThree := Contains(numbers, 3)  // true

words := []string{"hello", "world"}
hasGoodbye := Contains(words, "goodbye")  // false
```

## Generic Data Structures

Generics shine when building reusable data structures:

```go
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
    if len(s.items) == 0 {
        var zero T
        return zero, false
    }
    item := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return item, true
}
```

## Constraints

Type constraints let you specify what operations your generic types must support:

```go
type Numeric interface {
    ~int | ~int64 | ~float64
}

func Sum[T Numeric](values []T) T {
    var total T
    for _, v := range values {
        total += v
    }
    return total
}
```

## Best Practices

When working with generics:

1. Use them when you have genuine code duplication across types
2. Prefer simple constraints over complex ones
3. Don't over-genericize - sometimes concrete types are clearer
4. Consider whether an interface might be simpler

## Conclusion

Generics are a powerful addition to Go's toolkit. While they add complexity, they enable elegant solutions to problems that previously required awkward workarounds or duplication.
