package p2p

import "net"

// RPC represents any arbitrary data that is being sent over each transports
// between two nodes in the network.
type RPC struct {
	From 		net.Addr
	Payload 	[]byte
}