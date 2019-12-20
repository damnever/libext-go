package net

import (
	"errors"
	"net"
)

var (
	ErrNoAvailableIPAddress     = errors.New("libext-go/net: no available IP address")
	ErrNetworkInterfaceNotFound = errors.New("libext-go/net: network interface not found or not up")
)

// ResolveHostIP returns the non-loopback IP of a given network
// interface, the empty interfaceName means all network interfaces,
// only up interfaces will be checked.
func ResolveHostIP(interfaceName string) (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	interfaceFound := false
	for _, iface := range ifaces {
		if interfaceName != "" && iface.Name != interfaceName {
			continue
		}
		interfaceFound = true
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip.IsLoopback() || (ip.To4() == nil && ip.To16() == nil) {
				continue
			}
			return ip, nil
		}
	}

	if !interfaceFound {
		return nil, ErrNetworkInterfaceNotFound
	}
	return nil, ErrNoAvailableIPAddress
}
