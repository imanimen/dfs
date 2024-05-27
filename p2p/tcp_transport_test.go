package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	// listenAddr := ":4000"
	opts := TCPTransportOptions{
		ListenAddr: ":4000",
		HandShakeFunc: NOPHandshakeFunc,
		Decoder: DefaultDecoder{},
	}
	tr := NewTCPTransport(opts)
	assert.Equal(t, tr.ListenAddr, ":4000")

	// Server
	assert.Nil(t, tr.ListenAndAccept())
	
}