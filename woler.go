package woler

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
)

// WOLPort 标准Wake-on-LAN端口
const (
	WOLPort = 9
	// magicPacketMACRepetitions Magic Packet中MAC地址重复次数
	magicPacketMACRepetitions = 16
	// magicPacketPrefix Magic Packet前缀字节数
	magicPacketPrefix = 6
)

// parseMACAddress 解析并验证MAC地址
func parseMACAddress(macaddr string) ([]byte, error) {
	// 转换为大写并移除分隔符
	macaddr = strings.ToUpper(macaddr)
	macaddr = strings.NewReplacer("-", "", ":", "").Replace(macaddr)

	// 验证MAC地址长度
	if len(macaddr) != 12 {
		return nil, fmt.Errorf("invalid MAC address length: expected 12 characters, got %d", len(macaddr))
	}

	// 解码MAC地址
	machex, err := hex.DecodeString(macaddr)
	if err != nil {
		return nil, fmt.Errorf("invalid MAC address format: %v", err)
	}

	return machex, nil
}

// createMagicPacket 创建Wake-on-LAN魔术包
func createMagicPacket(machex []byte) []byte {
	var buffer bytes.Buffer
	// 写入6字节的0xFF
	buffer.Write(bytes.Repeat([]byte{0xFF}, magicPacketPrefix))
	// 写入16次MAC地址
	for i := 0; i < magicPacketMACRepetitions; i++ {
		buffer.Write(machex)
	}
	return buffer.Bytes()
}

// Do 使用标准广播地址发送Wake-on-LAN包
// macaddr: 支持的格式 AA-BB-CC-DD-EE-FF, AA:BB:CC:DD:EE:FF, AABBCCDDEEFF
func Do(macaddr string) error {
	return DoWithBroadcast(macaddr, net.IPv4bcast.String(), WOLPort)
}

// DoWithBroadcast 允许指定广播地址和端口的高级Wake-on-LAN实现
// macaddr: MAC地址
// bcastAddr: 广播地址，为空时使用默认广播地址
// port: 端口号，为0时使用标准WOL端口
func DoWithBroadcast(macaddr string, bcastAddr string, port int) error {
	// 解析MAC地址
	machex, err := parseMACAddress(macaddr)
	if err != nil {
		return err
	}

	// 设置默认端口
	if port == 0 {
		port = WOLPort
	}

	// 解析广播地址
	bcastIP := net.IPv4bcast
	if bcastAddr != "" {
		if ip := net.ParseIP(bcastAddr); ip != nil {
			bcastIP = ip
		} else {
			return fmt.Errorf("invalid broadcast address: %s", bcastAddr)
		}
	}

	// 建立UDP连接
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   bcastIP,
		Port: port,
	})
	if err != nil {
		return fmt.Errorf("failed to create UDP connection: %v", err)
	}
	defer conn.Close()

	// 创建并发送Magic Packet
	packet := createMagicPacket(machex)
	if _, err = conn.Write(packet); err != nil {
		return fmt.Errorf("failed to send magic packet: %v", err)
	}

	return nil
}
