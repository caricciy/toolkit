# Go Toolkit Library

This is a simple proof-of-concept (POC) project demonstrating how to write libraries in Go that can be imported into other projects.

## Overview

The library provides a `Tools` struct with utility functions. Currently, it includes a function to generate a random string of a given length.

## Installation

To install this library, run the following command:

```bash
go get github.com/caricciy/toolkit
```

## Usage

First, import the toolkit package into your Go file:

```go
import "github.com/caricciy/toolkit"
```

Instantiate the `Tools` struct:

```go
tools := toolkit.NewTools()
```

### RandonString

Use the `RandomString` function to generate a random string of a given length:

```go
randomString := tools.RandomString(10)
fmt.Println(randomString)
```

### UploadFiles

The `UploadFiles` function is used to handle file uploads. It takes in an `http.Request` object, a string representing the upload directory, and an optional boolean to determine if the file should be renamed.

```go
t := toolkit.Tools{
		MaxFileSize: 1024 * 1024 * 1024,
		AllowedFileTypes: []string{"image/jpeg", "image/png", "image/gif"},
	}
uploadedFiles, err := t.UploadFiles(request, "/path/to/upload/directory", false)
```

The `UploadFiles` function returns an `UploadedFile` struct used to save information about the uploaded file:

```go
type UploadedFile struct {
    NewFileName      string
    OriginalFileName string
    FileSize         int64
}
```

## License

This project is licensed under the MIT License - see the LICENSE.md file for details.