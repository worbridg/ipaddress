package ipaddress

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidIPAddress = errors.New("it isn't valid textual representation of an IP address")
	ErrInvalidPrefix    = errors.New("prefix must be 0-32")

	ClassA    = &IPv4Address{n: net.ParseIP("10.0.0.0"), prefix: 8}
	ClassB    = &IPv4Address{n: net.ParseIP("172.16.0.0"), prefix: 12}
	ClassC    = &IPv4Address{n: net.ParseIP("192.168.0.0"), prefix: 16}
	Multicast = &IPv4Address{n: net.ParseIP("224.0.0.0"), prefix: 4}
	Loopback  = &IPv4Address{n: net.ParseIP("127.0.0.0"), prefix: 8}
	LinkLocal = &IPv4Address{n: net.ParseIP("169.254.0.0"), prefix: 16}
)

// IPv4Address represents an IP address in version 4.
type IPv4Address struct {
	n      net.IP
	prefix int
}

// NewIPv4Address returns a new IPv4Address. addr formatted in IPv4 address is
// required to create an IPv4Address object. if addr formatted in CIDR is
// given, it is split by "/" to addr and prefix and used when IPv4Address is
// created. otherwise prefix is always 32. if invalid strings are given,
// returns an error.
func NewIPv4Address(addr string) (*IPv4Address, error) {
	addr, prefix, err := splitIPv4AddressToAddressAndPreix(addr)
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(addr)
	if ip == nil {
		return nil, ErrInvalidIPAddress
	}

	return &IPv4Address{
		n:      ip.To4(),
		prefix: prefix,
	}, nil
}

func splitIPv4AddressToAddressAndPreix(addr string) (string, int, error) {
	prefix := 32
	n := strings.Index(addr, "/")
	if n == -1 {
		return addr, prefix, nil
	}

	prefix, err := strconv.Atoi(addr[n+1:])
	if err != nil || !validatePrefix(prefix) {
		return "", 0, ErrInvalidPrefix
	}

	return addr[:n], prefix, nil
}

func validatePrefix(prefix int) bool {
	return prefix >= 0 && prefix <= 32
}

// String returns an IPv4 address string.
func (ipv4 *IPv4Address) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", ipv4.n[0], ipv4.n[1], ipv4.n[2], ipv4.n[3])
}

// Bits returns a string in binary format.
func (ipv4 *IPv4Address) Bits() string {
	bits := ""

	for _, b := range ipv4.n.To4() {
		bits = fmt.Sprintf("%s%08b", bits, b)
	}

	return bits
}

// Uint32 returns an unsigned number converted from the IP address.
func (ipv4 *IPv4Address) Uint32() uint32 {
	var n uint32 = 0
	i := 0
	for _, b := range ipv4.n.To4() {
		n += uint32(b) << (24 - (8 * i))
		i += 1
	}

	return n
}

// ToIPv4Address creates a IPv4 address with a unsigned number.
// prefix is always 32.
func ToIPv4Address(n uint32) *IPv4Address {
	return ToIPv4AddressWithPrefix(n, 32)
}

// ToIPv4AddressWithPrefix creates an IPv4 address with a unsigned number and prefix.
func ToIPv4AddressWithPrefix(n uint32, prefix int) *IPv4Address {
	if !validatePrefix(prefix) {
		return nil
	}

	b := make([]int, 4)
	b[0] = int(n) >> 24 & 0xff
	b[1] = int(n) >> 16 & 0xff
	b[2] = int(n) >> 8 & 0xff
	b[3] = int(n) & 0xff
	addr := fmt.Sprintf("%d.%d.%d.%d", b[0], b[1], b[2], b[3])

	return &IPv4Address{n: net.ParseIP(addr).To4(), prefix: prefix}
}

// Network returns a network IP address of version 4.
func (ipv4 *IPv4Address) Network() *IPv4Address {
	return ToIPv4AddressWithPrefix(
		ipv4.Uint32()&(0xffffffff<<(32-ipv4.prefix)),
		ipv4.prefix,
	)
}

// Network returns a broadcast IP address of version 4.
func (ipv4 *IPv4Address) Broadcast() *IPv4Address {
	return ToIPv4AddressWithPrefix(
		ipv4.Uint32()|(1<<(32-ipv4.prefix)-1),
		ipv4.prefix,
	)
}

// Equals compares the IP address to one of the other.
func (ipv4 *IPv4Address) Equals(other *IPv4Address) bool {
	return ipv4.Uint32() == other.Uint32()
}

// Next returns the next IP Address.
// if the address is greater than broadcast, returns nil.
func (ipv4 *IPv4Address) Next() *IPv4Address {
	if ipv4.Uint32() >= ipv4.Broadcast().Uint32() {
		return nil
	}

	return ToIPv4AddressWithPrefix(ipv4.Uint32()+1, ipv4.prefix)
}

// Prev returns the previous IP address.
// if the address is less than network, returns nil.
func (ipv4 *IPv4Address) Prev() *IPv4Address {
	if ipv4.Uint32() <= ipv4.Network().Uint32() {
		return nil
	}

	return ToIPv4AddressWithPrefix(ipv4.Uint32()-1, ipv4.prefix)
}

// ToIP converts the IP address to net.IP.
func (ipv4 *IPv4Address) ToIP() net.IP {
	return net.ParseIP(ipv4.String())
}

// Prefix returns the prefix of the IPv4 address.
func (ipv4 *IPv4Address) Prefix() int {
	return ipv4.prefix
}

// Contains checks if the other IPv4 address is in the IPv4 address range or not.
func (ipv4 *IPv4Address) Contains(other *IPv4Address) bool {
	return ipv4.Network().Uint32() <= other.Uint32() && other.Uint32() <= ipv4.Broadcast().Uint32()
}

// IsPrivate checks if the IPv4 address is private class(A/B/C) or not.
func (ipv4 *IPv4Address) IsPrivate() bool {
	classes := []*IPv4Address{ClassA, ClassB, ClassC}

	for _, c := range classes {
		if c.Contains(ipv4) {
			return true
		}
	}

	return false
}

// IsA checks if the IPv4 address is class A or not.
func (ipv4 *IPv4Address) IsA() bool {
	return ClassA.Contains(ipv4)
}

// IsB checks if the IPv4 address is class B or not.
func (ipv4 *IPv4Address) IsB() bool {
	return ClassB.Contains(ipv4)
}

// IsC checks if the IPv4 address is class C or not.
func (ipv4 *IPv4Address) IsC() bool {
	return ClassC.Contains(ipv4)
}

// IsMulticast checks if the IPv4 address is multicast or not.
func (ipv4 *IPv4Address) IsMulticast() bool {
	return Multicast.Contains(ipv4)
}

// IsLoopback checks if the IPv4 address is loopback or not.
func (ipv4 *IPv4Address) IsLoopback() bool {
	return Loopback.Contains(ipv4)
}

// IsLinkLocal checks if the IPv4 address is link local or not.
func (ipv4 *IPv4Address) IsLinkLocal() bool {
	return LinkLocal.Contains(ipv4)
}

// Bytes returns a byte slice of the IPv4 address.
func (ipv4 *IPv4Address) Bytes() []byte {
	n := ipv4.Uint32()
	return []byte{
		byte(n >> 24 & 0xff),
		byte(n >> 16 & 0xff),
		byte(n >> 8 & 0xff),
		byte(n & 0xff),
	}
}

// Netmask returns a netmask of the IPv4 address with prefix of that.
func (ipv4 *IPv4Address) Netmask() string {
	return ToIPv4Address(0xffffffff << (32 - ipv4.prefix)).String()
}

// Class returns a class name of the IPv4 address.
// if the class is neither A, B nor C, returns an empty string.
func (ipv4 *IPv4Address) Class() string {
	classes := map[string]func() bool{
		"A": ipv4.IsA,
		"B": ipv4.IsB,
		"C": ipv4.IsC,
	}

	for name, fn := range classes {
		if fn() {
			return name
		}
	}

	return ""
}

// Size returns the number of IPv4 addresses in the IPv4 address range.
func (ipv4 *IPv4Address) Size() int {
	return 1 << (32 - ipv4.prefix)
}

// Sample randomly picked an IPv4 address in the IPv4 address range up and returns it.
func (ipv4 *IPv4Address) Sample() *IPv4Address {
	rand.Seed(time.Now().UnixNano())
	n := uint32(rand.Intn(ipv4.Size() - 1))

	return ToIPv4AddressWithPrefix(
		ipv4.Network().Uint32()+n,
		ipv4.prefix,
	)
}

// ARPA returns a decimal string separated by dot with "in-addr.arpa" suffix.
func (ipv4 *IPv4Address) ARPA() string {
	b := ipv4.Bytes()
	return fmt.Sprintf("%d.%d.%d.%d.in-addr.arpa", b[3], b[2], b[1], b[0])
}
