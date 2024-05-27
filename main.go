package main

import (
	"fmt"
	"log"

	"github.com/imanimen/dfs/p2p"
)

func main() {
	tcpOpts := p2p.TCPTransportOptions{
		ListenAddr: ":3000",
		HandShakeFunc: p2p.NOPHandshakeFunc,
		Decoder: p2p.DefaultDecoder{},
		OnPeer: func(p2p.Peer) error {return fmt.Errorf("dfs: failed the OnPeer func")},
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