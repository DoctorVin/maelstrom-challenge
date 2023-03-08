package main

import (
	"log"

	"maelstrom-golang/lib"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

var noopHandler = func(_ maelstrom.Message) error {
	return nil
}

func main() {
	n := maelstrom.NewNode()
	bcast := lib.NewBroadcastHandler(n)
	top := lib.NewTopologyHandler(n)
	read := lib.NewReadHandler(n)
	n.Handle("broadcast", bcast)
	n.Handle("broadcast_ok", noopHandler)
	n.Handle("topology", top)
	n.Handle("read", read)

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
