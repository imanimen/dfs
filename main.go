package main

import (
	"log"

	"github.com/imanimen/dfs/p2p"
)

func main() {
	tcpOpts := p2p.TCPTransportOptions{
		ListenAddr: ":3000",
		HandShakeFunc: p2p.NOPHandshakeFunc,
		Decoder: p2p.GOBDecoder{},
	}
	tr := p2p.NewTCPTransport(tcpOpts)
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	select{}
}