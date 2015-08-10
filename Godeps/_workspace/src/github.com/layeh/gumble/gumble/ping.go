package gumble

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"io"
	"net"
	"time"
)

// PingResponse contains information about a server that responded to a UDP
// ping packet.
type PingResponse struct {
	// The address of the pinged server.
	Address *net.UDPAddr
	// The round-trip time from the client to the server.
	Ping time.Duration
	// The server's version. Only the Version field and SemanticVersion method of
	// the value will be valid.
	Version Version
	// The number users currently connected to the server.
	ConnectedUsers int
	// The maximum number of users that can connect to the server.
	MaximumUsers int
	// The maximum audio bitrate per user for the server.
	MaximumBitrate int
}

// Ping sends a UDP ping packet to the given server. Returns a PingResponse and
// nil on success. The function will return nil and an error if a valid
// response is not received after the given timeout.
func Ping(address string, timeout time.Duration) (*PingResponse, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	var packet [12]byte
	if _, err := rand.Read(packet[4:]); err != nil {
		return nil, err
	}
	start := time.Now()
	if _, err := conn.Write(packet[:]); err != nil {
		return nil, err
	}

	conn.SetReadDeadline(time.Now().Add(timeout))
	for {
		var incoming [24]byte
		if _, err := io.ReadFull(conn, incoming[:]); err != nil {
			return nil, err
		}
		if !bytes.Equal(incoming[4:12], packet[4:]) {
			continue
		}

		return &PingResponse{
			Address: addr,
			Ping:    time.Since(start),
			Version: Version{
				Version: binary.BigEndian.Uint32(incoming[0:]),
			},
			ConnectedUsers: int(binary.BigEndian.Uint32(incoming[12:])),
			MaximumUsers:   int(binary.BigEndian.Uint32(incoming[16:])),
			MaximumBitrate: int(binary.BigEndian.Uint32(incoming[20:])),
		}, nil
	}
}
