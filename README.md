# Single Client DHCP Server

Simple DHCP server serving configuration for a single client. This
server doesn't require IP in the same subnet as the offered IP.

## Build

```shell
go get github.com/krolaw/dhcp4
go build scdhcpd.go
```

## Usage

```shell
./main \
  CLIENT_MAC \
  CLIENT_IP \
  SERVER_IFACE \
  SERVER_IP \
  ROUTER_IP \
  DNS_IP
```

## Test

This will probably break your connectivity. Run as root.

```shell
# prepare namespaces for client and server
./prepare_netns

# run dhcp server on background
ip netns exec scdhcpd-server ./scdhcpd C4:4D:71:8D:AF:F8 192.168.1.2/24 veth_server 10.0.0.1 192.168.1.1 8.8.8.8 &

# run dhcp client
ip netns exec scdhcpd-client dhclient veth_client

# verify that client obtained its configuration
ip netns exec scdhcpd-client ip address show
```
