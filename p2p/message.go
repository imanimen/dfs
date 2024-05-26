package p2p

// Message represents any arbitrary data that is being sent over each transports
// between two nodes in the network.
type Message struct {
	Payload 	[]byte
}