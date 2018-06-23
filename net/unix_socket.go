package net

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"time"
)

// Socket stores the connection to a Unix Domain Socket and allows the communication through JSON messages between them
type Socket interface {
	Listen() error
	Accept() error
	Dial() error
	Read(msg interface{}) error
	Send(msg interface{}) error
	More() bool
	Close() error
	CloseListener() error
}

type unixSocket struct {
	address  string
	listener net.Listener
	conn     net.Conn
	decoder  *json.Decoder
	encoder  *json.Encoder
}

// NewSocket create a new socket object
func NewSocket(address string) Socket {
	return &unixSocket{address: address}
}

// Listen start a socket connection
func (socket *unixSocket) Listen() error {
	os.Remove(socket.address)
	listener, err := net.Listen("unix", socket.address)
	socket.listener = listener
	return err
}

// Accept connection on the listener
func (socket *unixSocket) Accept() error {
	conn, err := socket.listener.Accept()
	if err != nil {
		return err
	}
	socket.conn = conn
	socket.decoder = json.NewDecoder(conn)
	socket.encoder = json.NewEncoder(conn)
	return nil
}

// Dial connects a client to a server via Unix Sockets
func (socket *unixSocket) Dial() error {
	conn, err := net.Dial("unix", socket.address)
	if err != nil {
		return err
	}
	log.Println("Connected to socket: ", socket.address)
	socket.conn = conn
	socket.decoder = json.NewDecoder(conn)
	socket.encoder = json.NewEncoder(conn)
	return nil
}

func (socket *unixSocket) Read(msg interface{}) error {
	socket.conn.SetReadDeadline(time.Now().Add(time.Second))
	err := socket.decoder.Decode(msg)
	return err
}

// Send a message as json to the socket
func (socket *unixSocket) Send(msg interface{}) error {
	socket.conn.SetWriteDeadline(time.Now().Add(time.Second))
	err := socket.encoder.Encode(msg)
	return err
}

// More checks if there is more to read from the connection
func (socket unixSocket) More() bool {
	return socket.decoder.More()
}

// Close the socket connection
func (socket *unixSocket) Close() error {
	if socket.conn != nil {
		return socket.conn.Close()
	}
	return nil
}

// CloseListener closes the server listener
func (socket *unixSocket) CloseListener() error {
	if socket.listener != nil {
		return socket.listener.Close()
	}
	return nil
}
