package p2p

import (
	"fmt"
	"net"
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

// Close implements the Peer interface
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

type TCPTransportOptions struct {
	ListenAddr string
	HandShakeFunc HandShakeFunc
	Decoder    	  Decoder	
	OnPeer		  func(Peer) error
}

type TCPTransport struct {
	TCPTransportOptions
	Listener       net.Listener
	rpcCh			   chan RPC

}


// NewTCPTransport creates a new TCPTransport instance with the provided listen address.
// The TCPTransport is responsible for managing the TCP connection for the p2p network.
func NewTCPTransport(opts TCPTransportOptions) *TCPTransport {
	return &TCPTransport{
		TCPTransportOptions: opts,
		rpcCh: make(chan RPC),
	}
}

// Consume implements the Transport interface, which will return read-only channel
// for reading the incoming messages received from other peer in the network
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcCh
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


func (t *TCPTransport) handleConn(conn net.Conn) {

	var err error
	
	defer func() {
		fmt.Printf("dfs: dropping peer connection: %s", err)
		conn.Close()
	}()

	// create a new tcp peer
	peer := NewTCPPeer(conn, true)

	// HandShakeFunc is called to perform the handshake with the given peer.
	// If the handshake fails, the function returns an error.
	if err = t.HandShakeFunc(peer); err != nil {
		return
	}

	// OnPeer is called when a new peer is connected. If an error is returned,
	// the connection will be closed.
	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	// Read Loop
	rpc := RPC{}
	for {
		err := t.Decoder.Decode(conn, &rpc);
		if err == net.ErrClosed {
			return
		}
		if err != nil {
			fmt.Printf("dfs: TCP read error: %s\n", err)
			continue
		}
		rpc.From = conn.RemoteAddr()
		t.rpcCh <- rpc
	}

}