package network

import "fmt"

type IPAddr [4]byte

func (ip IPAddr) String() string {
	return fmt.Sprintf("%d.%d.%d.%d",
		ip[0], ip[1], ip[2], ip[3])
}
