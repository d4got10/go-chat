package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

var ERROR_EXIT_CODE int = 99

func main() {
    var address string = "127.0.0.1:6969"
    if len(os.Args) > 1 {
        address = os.Args[1]
    }

    conn, err := net.Dial("tcp", address)
    if err != nil {
        log_error(err.Error())
        os.Exit(ERROR_EXIT_CODE)
    }
    fmt.Printf("Connected to server!\n")
    go receive(conn)

    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(conn)
    for {
        input, err := reader.ReadString('\n')
        if err != nil {
            log_error("[Stdin read error] " + err.Error())
            break
        }
        _, err = writer.WriteString(input)
        writer.Flush()

        if err != nil {
            log_error(err.Error())
            break
        }
    }
}

func receive(conn net.Conn) {
    reader := bufio.NewReader(conn)
    for {
        message, err := reader.ReadString('\n')
        if err != nil {
            if !errors.Is(err, io.EOF) {
                log_error(err.Error())
            }
            conn.Close()
            return
        }
        message = strings.Trim(message, "\n")
        fmt.Printf("%s\n", message)
    }
}

func log_error(message string) {
    fmt.Printf("Error: %s\n", message)
}

