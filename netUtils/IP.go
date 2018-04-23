package netUtils

import (
	"errors"
	"net"
)

//IP takes a string and returns the net.IP representation of it if ip is a valid ip.
//Returns error otherwise.
func IP(ip string) (net.IP, error) {
	ret := net.ParseIP(ip)
	if ret == nil {
		err := errors.New("Not an IP")
		return ret, err
	}

	return ret, nil
}

func isZeros(slice []byte) bool {
	for _, elem := range slice {
		if elem != 0 {
			return false
		}
	}
	return true
}

//IsIPv4 returns true if either the passed IP has len 4 or is a v4 addr in v6 representation.
func IsIPv4(ip net.IP) bool {
	if len(ip) == 4 ||
		len(ip) == net.IPv6len && isZeros(ip[0:10]) && ip[10] == 0xff && ip[11] == 0xff {
		return true
	}
	return false
}

//IsIPv6 returns true if the passed IP has len(16) and is not a v4 addr in v6 representation
func IsIPv6(ip net.IP) bool {
	if len(ip) == 16 && !IsIPv4(ip) {
		return true
	}
	return false
}
