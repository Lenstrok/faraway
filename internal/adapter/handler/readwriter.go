package handler

import (
	"encoding/binary"
	"io"
	"net"
)

// ReadMessage reads a message from the connection.
func ReadMessage(conn net.Conn) (msg []byte, err error) {
	var length uint64
	if err = binary.Read(conn, binary.BigEndian, &length); err != nil {
		return
	}

	msg = make([]byte, length)
	_, err = io.ReadFull(conn, msg)
	return
}

// WriteMessage writes a message to the connection.
func WriteMessage(conn net.Conn, msg []byte) (err error) {
	if err = binary.Write(conn, binary.BigEndian, uint64(len(msg))); err != nil {
		return
	}
	_, err = conn.Write(msg)
	return
}
