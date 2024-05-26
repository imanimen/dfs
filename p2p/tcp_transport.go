package p2p

import (
	"fmt"
	"net"
	"sync"
)

type TCPTransport struct {
	ListenAddress  string
	Listener 	   net.Listener

	mu 			   sync.RWMutex
	peers 		   map[net.Addr]Peer

}

// NewTCPTransport creates a new TCPTransport instance with the provided listen address.
// The TCPTransport is responsible for managing the TCP connection for the p2p network.
func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		ListenAddress: listenAddr,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.Listener, err = net.Listen("tcp", t.ListenAddress)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()
	return nil

}


func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.Listener.Accept()
		if err != nil {
			fmt.Printf("dfs: TCP accept error: %s\n", err)
		}

		go t.handleConn(conn)
	}
}


func (t *TCPTransport) handleConn(conn net.Conn) {
	fmt.Printf("dfs: new incoming connection %+v\n", conn)
}