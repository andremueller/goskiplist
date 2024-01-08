# Skip Lists in Go

[![Go](https://github.com/andremueller/goskiplist/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/andremueller/goskiplist/actions/workflows/go.yml)

## Introduction

This is a new skip list library for Go supporting access by a key and index.

## Objectives

- Using Go generics for key and value.
- Simple maintainable code
- Allowing fast random access (GetByPos) of the nth element by combining the basic skip list algorithm in with the linear list operations (c.f. 1. Section 3.4). The indexed access has an average runtime of $O(log(n))$.

## Installation

```bash
go get github.com/andremueller/goskiplist/pkg/skiplist
```



## Usage Example

```go
// creates a skip list with key type `int` and value type `string`
s := skiplist.NewSkipList[int, string]()

s.Set(1, "cat")
s.Set(2, "dog")

x, _ := s.Get(1)
fmt.Printf("Value: %s", x.Value)
// Output  "Value: cat"

x, _ = s.GetByPos(1)
fmt.Printf("Value: %s", x.Value)
// Output  "Value: dog"
```

## References

- 1. W. Pugh, A Skip List Cookbook. (2001).

## MIT License

See [License](LICENSE).
