package main

import (
	"bytes"
	"log"
	"net"
	"os"
	"time"

	dhcp "github.com/krolaw/dhcp4"
	dhcpConn "github.com/krolaw/dhcp4/conn"
)

func main() {
	clientMAC, _ := net.ParseMAC(os.Args[1])
	clientIP, clientIPNet, _ := net.ParseCIDR(os.Args[2])
	clientMask := clientIPNet.Mask
	serverIface := os.Args[3]
	serverIP := net.ParseIP(os.Args[4])
	routerIP := net.ParseIP(os.Args[5])
	dnsIP := net.ParseIP(os.Args[6])
	SingleClientDHCPServer(clientMAC, clientIP, clientMask, serverIface, serverIP, routerIP, dnsIP)
}

func SingleClientDHCPServer(
	clientMAC net.HardwareAddr,
	clientIP net.IP,
	clientMask net.IPMask,
	serverIface string,
	serverIP net.IP,
	routerIP net.IP,
	dnsIP net.IP) {

	handler := &DHCPHandler{
		clientIP:      clientIP,
		clientMAC:     clientMAC,
		serverIP:      serverIP,
		leaseDuration: 999 * 24 * time.Hour,
		options: dhcp.Options{
			dhcp.OptionSubnetMask:       []byte(clientMask),
			dhcp.OptionRouter:           []byte(routerIP),
			dhcp.OptionDomainNameServer: []byte(dnsIP),
		},
	}

	l, err := dhcpConn.NewUDP4BoundListener(serverIface, ":67")
	if err != nil {
		return
	}
	defer l.Close()
	log.Fatal(dhcp.Serve(l, handler))
}

type DHCPHandler struct {
	serverIP      net.IP
	clientIP      net.IP
	clientMAC     net.HardwareAddr
	leaseDuration time.Duration
	options       dhcp.Options
}

func (h *DHCPHandler) ServeDHCP(p dhcp.Packet, msgType dhcp.MessageType, options dhcp.Options) (d dhcp.Packet) {
	switch msgType {

	case dhcp.Discover:
		if mac := p.CHAddr(); !bytes.Equal(mac, h.clientMAC) {
			return nil // Is not our client
		}
		return dhcp.ReplyPacket(p, dhcp.Offer, h.serverIP, h.clientIP, h.leaseDuration,
			h.options.SelectOrderOrAll(options[dhcp.OptionParameterRequestList]))

	case dhcp.Request:
		if mac := p.CHAddr(); !bytes.Equal(mac, h.clientMAC) {
			return nil // Is not our client
		}
		if server, ok := options[dhcp.OptionServerIdentifier]; ok && !net.IP(server).Equal(h.serverIP) {
			return nil // Message not for this DHCP server
		}
		return dhcp.ReplyPacket(p, dhcp.ACK, h.serverIP, h.clientIP, h.leaseDuration,
			h.options.SelectOrderOrAll(options[dhcp.OptionParameterRequestList]))

	default:
		return nil
	}
}
