package main

import (
	"fmt"
	"log"

	"github.com/imanimen/dfs/p2p"
)

func OnPeer(peer p2p.Peer) error {
	// fmt.Println("dfs: doing some logic with the peer outside of TCPTransport")
	peer.Close() // close the peer connection and enters the loop <- FIX #4 pr
	return nil
}


func main() {
	tcpOpts := p2p.TCPTransportOptions{
		ListenAddr: ":3000",
		HandShakeFunc: p2p.NOPHandshakeFunc,
		Decoder: p2p.DefaultDecoder{},
		OnPeer: OnPeer,
	}
	tr := p2p.NewTCPTransport(tcpOpts)

	go func() {
		for {
			msg :=  <-tr.Consume()
			fmt.Printf("dfs: message %+v\n", msg)
		}
	}()
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	select{}
}