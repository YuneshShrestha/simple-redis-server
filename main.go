package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

const addr = "0.0.0.0:6379"

func main() {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("Failed to bind to %s\n", addr)
		os.Exit(1)
	}

	fmt.Printf("Starting Redis server at %s\n", addr)

	// For inputs
	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		go func(conn net.Conn) {
			defer conn.Close()

			buf := make([]byte, 2014)

			for {
				len, err := conn.Read(buf)
				if err != nil {
					if err != io.EOF {
						fmt.Printf("Error reading: %#v\n", err)
					}
					break
				}

				// Buffer has length of 2014 so for command take only necessary part
				command := buf[:len]
				response := Parse(command)

				_, responseErr := conn.Write([]byte(response + "\r\n"))
				if responseErr != nil {
					fmt.Printf("Error writing: %#v\n", responseErr)
					break
				}
			}
		}(conn)
	}
}
