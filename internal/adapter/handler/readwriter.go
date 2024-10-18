package handler

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// ReadMessage reads a message from the connection.
func ReadMessage(conn net.Conn) ([]byte, error) {
	var length uint64
	if err := binary.Read(conn, binary.BigEndian, &length); err != nil {
		return nil, fmt.Errorf("failed to binary.Read: %w", err)
	}

	msg := make([]byte, length)
	if _, err := io.ReadFull(conn, msg); err != nil {
		return nil, fmt.Errorf("failed to io.ReadFull: %w", err)
	}

	return msg, nil
}

// WriteMessage writes a message to the connection.
func WriteMessage(conn net.Conn, msg []byte) error {
	if err := binary.Write(conn, binary.BigEndian, uint64(len(msg))); err != nil {
		return fmt.Errorf("binary.Write failed: %w", err)
	}

	if _, err := conn.Write(msg); err != nil {
		return fmt.Errorf("conn.Write failed: %w", err)
	}

	return nil
}
