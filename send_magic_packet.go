package wol

import (
	"fmt"
	"net"
	"strconv"

	"github.com/ahiggins0/go-wol/wol"
)

func sendPacket(addr *net.UDPAddr, data []byte) error {
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write(data)
	return err
}

func SendMagicPacket(mac string, broadcastIP string, portStr string) error {
	packet, err := bmp.New(mac)
	if err != nil {
		fmt.Println("Error creating magic packet:", err)
		return fmt.Errorf("error creating magic packet: %w", err)
	}

	data, err := packet.Marshal()
	if err != nil {
		fmt.Println("Error marshaling magic packet:", err)
		return fmt.Errorf("error marshaling magic packet: %w", err)
	}

	// Default port to 9
	if portStr == "" {
		portStr = "9"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 || port > 65535 {
		return fmt.Errorf("invalid port: %v", portStr)
	}

	// If a broadcast IP is provided, send to it directly
	if broadcastIP != "" {
		udpAddr := &net.UDPAddr{
			IP:   net.ParseIP(broadcastIP),
			Port: port,
		}
		if udpAddr.IP == nil {
			return fmt.Errorf("invalid broadcast IP: %s", broadcastIP)
		}
		fmt.Printf("Magic packet sent to %s\n", broadcastIP)
		return sendPacket(udpAddr, data)
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error getting interfaces:", err)
		return fmt.Errorf("error getting interfaces: %w", err)
	}

	for _, iface := range interfaces {
		// Skip interfaces that are down or loopback
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok || ipnet.IP == nil || ipnet.IP.To4() == nil {
				continue
			}

			ip := ipnet.IP.To4()
			mask := ipnet.Mask

			// Calculate broadcast address
			broadcast := make(net.IP, 4)
			for i := 0; i < 4; i++ {
				broadcast[i] = ip[i] | ^mask[i]
			}

			udpAddr := &net.UDPAddr{
				IP:   broadcast,
				Port: port,
			}

			err = sendPacket(udpAddr, data)

			if err != nil {
				fmt.Printf("Interface %s: Send error: %v\n", iface.Name, err)
			} else {
				fmt.Printf("Magic packet sent via %s to %s on port %s\n", iface.Name, broadcast.String(), portStr)
			}
		}
	}

	return nil
}