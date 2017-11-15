package utils

import (
	"errors"
	"net"
)

//IP takes a string and returns the net.IP representation of it if ip is a valid ip
//returns error otherwise
func IP(ip string) (net.IP, error) {
	ret := net.ParseIP(ip)
	if ret == nil {
		err := errors.New("Not an IP")
		return ret, err
	}
	return ret, nil
}
