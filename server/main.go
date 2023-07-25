package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

var ERROR_EXIT_CODE int = 99

type server_state struct {
    clients []net.Conn
}

type clients_connection_states = map[net.Conn]bool

func main() {
    var address string = ":6969"
    if len(os.Args) > 1 {
        address = os.Args[1]
    }
    listener, err := net.Listen("tcp", address)
    defer listener.Close()

    if err != nil {
        log_error(err.Error())
        os.Exit(ERROR_EXIT_CODE)
    }

    var state *server_state = new(server_state)
    state.clients = make([]net.Conn, 0)

    fmt.Printf("Server listening on %s\n", listener.Addr())
    for {
        fmt.Printf("Waiting for client connection...\n")
        conn, err := listener.Accept()
        if err != nil {
            log_error(err.Error())
            os.Exit(ERROR_EXIT_CODE)
        }
        fmt.Printf("Client connected: %s\n", conn.RemoteAddr())
        state.clients = append(state.clients, conn)
        go handle_client(conn, state)
    }
}

func remove_disconnected_clients(clients []net.Conn, connection_states clients_connection_states) []net.Conn {
    connected_clients := make([]net.Conn, 0)
    for _, client := range clients {
        if connection_states[client] {
            connected_clients  = append(connected_clients, client)
        } else {
            fmt.Printf("Client [%s] disconnected.\n", client.RemoteAddr())
        }
    }
    return connected_clients
}

func send_to_all_except(conn net.Conn, state *server_state, message string) {
    clients_connection_states := make(clients_connection_states)
    clients_connection_states[conn] = true
    for _, other_conn := range state.clients {
        if conn == other_conn {
            continue
        }
        writer := bufio.NewWriter(other_conn)
        _, err := writer.WriteString(message)
        err = writer.Flush()
        if err != nil {
            log_error(err.Error())
            clients_connection_states[other_conn] = false
            continue
        }
        clients_connection_states[other_conn] = true
    }
    state.clients = remove_disconnected_clients(state.clients, clients_connection_states)
}

func handle_client(conn net.Conn, state *server_state) {
    reader := bufio.NewReader(conn)
    fmt.Printf("Started handling client [%s]\n", conn.RemoteAddr())
    for {
        message, err := reader.ReadString('\n')
        if err != nil {
            if !errors.Is(err, io.EOF) {
                log_error(err.Error())
            }

            conn.Close()
            return
        }
        fmt.Printf("Received message [%s]: %s", conn.RemoteAddr(), message)
        send_to_all_except(conn, state, message)
    }
}

func log_error(message string) {
    fmt.Printf("Error: %s\n", message)
    panic("123")
}
