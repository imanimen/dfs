package p2p

import "errors"

// ErrInvalidHandshake is returned if the handshake between the local and remote peer fails.
var ErrInvalidHandshake = errors.New("invalid handshake")

// HandShakeFunc is a function that is called when a new peer connects to the network.
type HandShakeFunc func(Peer) error

// NOPHandshakeFunc is a no-op handshake function that always returns nil.
func NOPHandshakeFunc(Peer) error { return nil }