package network // main.goと同じパッケージ名

import (
	"bytes"
	"fmt"
	"net"
)

// SendMagicPacket は指定したMACアドレスに対してUDPブロードキャストでマジックパケットを送信します。
func SendMagicPacket(macAddr string) error {
	// 1. MACアドレスのパース
	hwAddr, err := net.ParseMAC(macAddr)
	if err != nil {
		return fmt.Errorf("invalid MAC address: %w", err)
	}

	// 2. マジックパケットの構築 (6*0xFF + 16*MAC)
	var packet bytes.Buffer
	packet.Write(bytes.Repeat([]byte{0xFF}, 6))
	for i := 0; i < 16; i++ {
		packet.Write(hwAddr)
	}

	// 3. UDP接続の設定 (ブロードキャストアドレスを指定)
	// 一般的なポートは 7 または 9 です
	addr, err := net.ResolveUDPAddr("udp", "255.255.255.255:9")
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// 4. パケットの送信
	_, err = conn.Write(packet.Bytes())
	if err != nil {
		return err
	}

	fmt.Printf("Magic packet sent to %s\n", macAddr)
	return nil
}