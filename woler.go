package woler

import (
	"bytes"
	"encoding/hex"
	"net"
	"strings"
)

// Do: devices physical address that want to wake up
// macaddr: AA-BB-CC-DD-EE-FF, AA:BB:CC:DD:EE:FF, AABBCCDDEEFF
func Do(macaddr string) error {
	switch {
	case strings.Contains(macaddr, "-"):
		macaddr = strings.ReplaceAll(macaddr, "-", "")
	case strings.Contains(macaddr, ":"):
		macaddr = strings.ReplaceAll(macaddr, ":", "")
	}
	machex, err := hex.DecodeString(macaddr)
	if err != nil {
		return err
	}
	// Establish UDP connection with broadcast address
	conn, err := net.DialUDP("udp", &net.UDPAddr{}, &net.UDPAddr{IP: net.IPv4bcast})
	if err != nil {
		return err
	}
	defer conn.Close()
	var buffer bytes.Buffer
	buffer.Write([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})
	for i := 0; i < 16; i++ {
		buffer.Write(machex)
	}
	if _, err = conn.Write(buffer.Bytes()); err != nil {
		return err
	}
	return nil
}
