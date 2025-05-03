# Go Wake-on-LAN (WOL) Package

A simple Go package for sending Wake-on-LAN (WOL) magic packets over the network. This package allows you to construct and broadcast a WOL packet to wake up devices using their MAC address. This was essentially forked from https://github.com/sabhiram/go-wol with a few changes I wanted to include.

## Features

- Generate valid WOL magic packets from a MAC address
- Send packets to a specified broadcast IP or auto-detect via local interfaces
- Default to port 9, or use a custom port

## Installation

```bash
go get github.com/ahiggins0/go-wol
```

## Usage
### Sending a Magic Packet
```
package main

import (
	"log"

	"github.com/ahiggins0/go-wol"
)

func main() {
	err := wol.SendMagicPacket("00:11:22:33:44:55", "", "") // send to all available broadcast addresses over port 9
	if err != nil {
		log.Fatal(err)
	}
}
```
### Send to a specific IP/Port
```
package main

import (
	"log"

	"github.com/ahiggins0/go-wol"
)

func main() {
	err := wol.SendMagicPacket("00:11:22:33:44:55", "255.255.255.255", "7") // send to specific broadcast address over port 7
	if err != nil {
		log.Fatal(err)
	}
}
```

## Parameters
`mac`: MAC address of the target device. Format: 00:11:22:33:44:55.  
`broadcastIP`: (Optional) IP to send the packet to. Leave blank to auto-detect.  
`port`: (Optional) UDP port to use. Defaults to 9.

## Package Overview
### Package `bmp`

Responsible for building a magic packet.
- `New(mac string) (*MagicPacket, error)`: Constructs a magic packet from a MAC address.
- `(*MagicPacket).Marshal() ([]byte, error)`: Serializes the packet to raw bytes.

### Package `wol`
Responsible for sending the packet over the network.
- `SendMagicPacket(mac string, broadcastIP string, port string) error`: Sends the packet over UDP.

## Example Output
```
Magic packet sent via eth0 to 192.168.1.255 on port 9
Magic packet sent via eth1 to 172.25.1.255 on port 9
```