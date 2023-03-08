package lib

import (
	"fmt"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type ReplyBody map[string]any

var idSeed int64

func GenerateUniqueID(n *maelstrom.Node) ReplyBody {
	idSeed++
	b := make(ReplyBody)
	b["type"] = "generate_ok"
	b["id"] = fmt.Sprintf("%s-%d", n.ID(), idSeed)
	return b
}
