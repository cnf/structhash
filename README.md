# structhash [![GoDoc](https://godoc.org/github.com/cnf/structhash?status.svg)](https://godoc.org/github.com/cnf/structhash)

structhash is a Go library for generating hash strings of arbitrary data structures.

## Installation

Standard `go get`:

```
$ go get github.com/cnf/structhash
```

## Documentation

For usage and examples see the [Godoc](http://godoc.org/github.com/cnf/structhash).

## Quick start

```go
package main

import (
    "fmt"
    "crypto/md5"
    "crypto/sha1"
    "github.com/gavv/structhash"
)

type S struct {
    Str string
    Num int
}

func main() {
    s := S{"hello", 123}

    hash, err := structhash.Hash(s, 1)
    if err != nil {
        panic(err)
    }
    fmt.Println(hash)
    // Prints: v1_55743877f3ffd5fc834e97bc43a6e7bd

    fmt.Printf("%x\n", structhash.Md5(s, 1))
    // Prints: 55743877f3ffd5fc834e97bc43a6e7bd

    fmt.Printf("%x\n", structhash.Sha1(s, 1))
    // Prints: 00f550e11183e2bb70f8bf12699c3866e5c8fcb3

    fmt.Printf("%x\n", md5.Sum(structhash.Dump(s, 1)))
    // Prints: 55743877f3ffd5fc834e97bc43a6e7bd

    fmt.Printf("%x\n", sha1.Sum(structhash.Dump(s, 1)))
    // Prints: 00f550e11183e2bb70f8bf12699c3866e5c8fcb3
}
```
