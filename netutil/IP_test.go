package netutil

import (
	"testing"
	"net"
	"github.com/stretchr/testify/assert"
)

func TestIP(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testCases := []struct {
		desc	string
		ipStr string
		ip net.IP
	}{
		{
			desc: "Valid IPv4",
			ipStr: "127.0.0.1",
			ip: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff, 127, 0, 0, 1},
		},
		{
			desc: "Valid IPv6 without braces",
			ipStr: "fe80:db8::68",
			ip: []byte{0xfe, 0x80, 0x0d, 0xb8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x68},
		},
		{
			desc: "Invalid IPv4 too long",
			ipStr: "127.0.0.0.1",
			ip: []byte{},
		},
		{
			desc: "Invalid IPv4 bad character",
			ipStr: "127.0.0,1",
			ip: []byte{},
		},
		{
			desc: "Invalid IPv6",
			ipStr: "fe80:db8:::68",
			ip: []byte{},
		},
		{
			desc: "Valid IPv6 with braces",
			ipStr: "[2001:db8::68]",
			ip: []byte{},
		},
		{
			desc: "Not an IP",
			ipStr: "foobar",
			ip: []byte{},
		},
	}
	for i, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ip, err := IP(tC.ipStr)
			if i < 2 {
				assert.Nil(err, "Failed because err was not nil.")
				assert.Equal(tC.ip, ip, "Failed because the ip was not correct.")
			} else {
				assert.Nil(ip, "Failed because ip was not nil.")
				assert.EqualError(err, "Not an IP", "Failed because the error was not correct.")
			}
		})
	}
}

func TestIsZeros(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testCases := []struct {
		desc	string
		b []byte
	}{
		{
			desc: "Zero byte slice",
			b: []byte{0,0,0,0,0,0,0,0,0,0},
		},
		{
			desc: "Non zero byte slice",
			b: []byte{0,0,0,0,0,1,0,0,0,0},
		},
	}
	for i, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			b := isZeros(tC.b)
			if i == 0 {
				assert.True(b, "Failed because the slice was not recognized as zeroed even though it is.")
			} else {
				assert.False(b, "Failed because the slice was not recognized as not zeroed even though it is.")
			}
		})
	}
}

func TestIsIPv4(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testCases := []struct {
		desc	string
		ip net.IP
	}{
		{
			desc: "Valid IPv4",
			ip: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff, 127, 0, 0, 1},
		},
		{
			desc: "Valid IPv6",
			ip: []byte{0xfe, 0x80, 0x0d, 0xb8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x68},
		},
		{
			desc: "Invalid IPv4 too long",
			ip: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff, 127, 0, 0, 1},
		},
		{
			desc: "Invalid IPv4 bad pattern",
			ip: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0, 0xff, 127, 0, 0, 1},
		},
		{
			desc: "Invalid IPv6",
			ip: []byte{},
		},
		{
			desc: "Valid IPv6 with braces",
			ip: []byte{},
		},
		{
			desc: "Not an IP",
			ip: []byte{},
		},
	}
	for i, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			b := IsIPv4(tC.ip)
			if i == 0 {
				assert.True(b, "Failed because a valid ipv4 was not recognized as one.")
			} else {
				assert.False(b, "Failed because an invalid ipv4 was recognized as one.")
			}
		})
	}
}

func TestIsIPv6(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testCases := []struct {
		desc	string
		ip net.IP
	}{
		{
			desc: "Valid IPv6",
			ip: []byte{0xfe, 0x80, 0x0d, 0xb8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x68},
		},
		{
			desc: "Valid IPv4",
			ip: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff, 127, 0, 0, 1},
		},
		{
			desc: "Invalid IPv4 too long",
			ip: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff, 127, 0, 0, 1},
		},
		{
			desc: "Invalid IPv4 bad pattern",
			ip: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0, 0xff, 127, 0, 0, 1},
		},
		{
			desc: "Invalid IPv6",
			ip: []byte{},
		},
		{
			desc: "Not an IP",
			ip: []byte{},
		},
	}
	for i, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			b := IsIPv6(tC.ip)
			if i == 0 {
				assert.True(b, "Failed because a valid ipv6 was not recognized as one.")
			} else {
				assert.False(b, "Failed because an invalid ipv6 was recognized as one.")
			}
		})
	}
}

func TestBuildIPAddressString(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testCases := []struct {
		desc	string
		ip net.IP
		port int
		ippString string
	}{
		{
			desc: "Valid IPv4",
			ip: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff, 127, 0, 0, 1},
			port: 50000,
			ippString: "127.0.0.1:50000",
		},
		{
			desc: "Valid IPv6",
			ip: []byte{0xfe, 0x80, 0x0d, 0xb8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x68},
			port: 30012,
			ippString: "[fe80:db8::68]:30012",
		},
		{
			desc: "Invalid IPv4 too long",
			ip: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff, 127, 0, 0, 1},
			port: 30123,
			ippString: "",
		},
		{
			desc: "Invalid IPv4 bad pattern",
			ip: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0, 0xff, 127, 0, 0, 1},
			port: 12315,
			ippString: "",
		},
		{
			desc: "Invalid IPv6",
			ip: []byte{},
			port: 5123,
			ippString: "",
		},
		{
			desc: "Valid IPv6 with braces",
			ip: []byte{},
			port: 53142,
			ippString: "",
		},
		{
			desc: "Not an IP",
			ip: []byte{},
			port: 52312,
			ippString: "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			addr := BuildIPAddressString(tC.ip, tC.port)
			assert.Equal(tC.ippString, addr, "Failed because the returned IPAddress String is not correct.")
		})
	}
}