package main

import (
	"log"

	"maelstrom-golang/lib"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()
	n.Handle("generate", func(msg maelstrom.Message) error {
		repBody := lib.GenerateUniqueID(n)
		return n.Reply(msg, repBody)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
