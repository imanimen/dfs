package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCPPeer represents the remote node over tcp established connection.
type TCPPeer struct {
	// conn is the underlying connection of the peer
	conn	net.Conn
	// if we dial and retrieve a connection 	-> outbound == true
	// if we accept and retrieve a connection 	-> outbound == false
	outbound	bool
}

// NewTCPPeer creates a new TCPPeer instance from the provided net.Conn and outbound flag.
// The TCPPeer represents a peer connected over a TCP transport.
func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransportOptions struct {
	ListenAddr string
	HandShakeFunc HandShakeFunc
	Decoder    	  Decoder	
}

type TCPTransport struct {
	TCPTransportOptions
	Listener       net.Listener
	mu 			   sync.RWMutex
	peers 		   map[net.Addr]Peer

}


// NewTCPTransport creates a new TCPTransport instance with the provided listen address.
// The TCPTransport is responsible for managing the TCP connection for the p2p network.
func NewTCPTransport(opts TCPTransportOptions) *TCPTransport {
	return &TCPTransport{
		TCPTransportOptions: opts,
	}
}

// ListenAndAccept starts listening on the configured address and accepts incoming TCP connections.
// It returns an error if the listener cannot be created or if there is an error accepting connections.
func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.Listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()
	return nil

}


// startAcceptLoop listens for incoming TCP connections and handles them in separate goroutines.
// It runs in a loop, accepting connections and passing them to handleConn to be processed.
// If there is an error accepting a connection, it is logged but the loop continues.
func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.Listener.Accept()
		if err != nil {
			fmt.Printf("dfs: TCP accept error: %s\n", err)
		}

		fmt.Printf("dfs: new incoming connection %+v\n", conn)

		go t.handleConn(conn)
	}
}

type Temp struct {}

func (t *TCPTransport) handleConn(conn net.Conn) {
	// create a new tcp peer
	peer := NewTCPPeer(conn, true)

	if err := t.HandShakeFunc(peer); err != nil {
		conn.Close()
		fmt.Printf("dfs: handshake error: %s\n", err)
		return
	}


	// Read Loop
	msg := &Temp{}
	for {
		if err := t.Decoder.Decode(conn, msg); err != nil {
			fmt.Printf("dfs: TCP read error: %s\n", err)
			continue
		}
	}

}