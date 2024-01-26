# Go Toolkit Library

This is a simple proof-of-concept (POC) project demonstrating how to write libraries in Go that can be imported into other projects.

## Overview

The library provides a `Tools` struct with utility functions. Currently, it includes a function to generate a random string of a given length.

## Usage

First, import the toolkit package into your Go file:

```go
import "github.com/caricciy/toolkit"
```

Instantiate the `Tools` struct:

```go
tools := toolkit.NewTools()
```

Use the `RandomString` function to generate a random string of a given length:

```go
randomString := tools.RandomString(10)
fmt.Println(randomString)
```

## Installation

To install this library, run the following command:

```bash
go get github.com/caricciy/toolkit
```



## License

This project is licensed under the MIT License - see the LICENSE.md file for details.