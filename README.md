### What is a file in binary?
At the lowest level, a file is just a sequence of bytes (binary data) stored on disk. It doesn’t matter if it’s a .txt, .png, .exe, or .mp3 — it’s just a long stream of bytes.

For example:

If you open a file and see this in hex view:

````
48 65 6C 6C 6F 20 57 6F 72 6C 64

That’s the binary for:

Hello World
````
So a .txt file with "Hello World" in it looks like that in raw bytes.

### What is an io.Reader in Go?

In Go, io.Reader is an interface that lets you read a stream of bytes one chunk at a time. It’s very simple:

````go
type Reader interface {
    Read(p []byte) (n int, err error)
}
````
You give it a byte slice (p []byte), and it fills it with data, telling you how many bytes it wrote.


###  io.Seeker lets you move around in a file


When reading from a file, you might not always want to read from the start. You might want to "jump" (seek) to a different position.

io.Seeker lets you move the "cursor" to any byte offset in the file:
````go
type Seeker interface {
    Seek(offset int64, whence int) (int64, error)
}
````
offset: how far to move
whence: where to move from:

0 = from start
1 = from current
2 = from end

### Sample example

```go

package main

import (
    "fmt"
    "io"
    "os"
)

func main() {
    // Create a sample file with some data
    os.WriteFile("sample.txt", []byte("Hello GoLang"), 0644)

    // Open the file
    f, _ := os.Open("sample.txt")
    defer f.Close()

    // Read first 5 bytes
    buf := make([]byte, 5)
    n, _ := f.Read(buf)
    fmt.Println("First 5 bytes:", string(buf[:n])) // Output: Hello

    // Seek to byte 6 from beginning
    f.Seek(6, io.SeekStart) // skip "Hello "

    // Read next 6 bytes
    buf2 := make([]byte, 6)
    n2, _ := f.Read(buf2)
    fmt.Println("Next 6 bytes:", string(buf2[:n2])) // Output: GoLang
}

```


### Setup 

- use same go version
- go mod tidy
- go mod vendor
- go run main.go
- access localhost:8080


### High level

- storage directory -> is kind of database for demo purpose.
- chunking each file and saving. 
- while updating, only saving the changed chunks.
- at time of downloading merging all chunks. 