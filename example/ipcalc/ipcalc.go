package main

import (
	"fmt"
	"os"

	"github.com/worbridg/ipaddress"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: ipcalc <ipaddress|cidr>\n")
	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}

	ipv4, err := ipaddress.NewIPv4Address(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("Address = %s\n", ipv4)
	fmt.Printf("CIDR = %s/%d\n", ipv4, ipv4.Prefix())
	fmt.Printf("Netmask = %s\n", ipv4.Netmask())
	fmt.Printf("Network = %s\n", ipv4.Network())
	fmt.Printf("Broadcast = %s\n", ipv4.Broadcast())
	fmt.Printf("Range = %s to %s\n", ipv4.Network().Next(), ipv4.Broadcast().Prev())
	fmt.Printf("Class %s\n", ipv4.Class())
}
