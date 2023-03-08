package lib

import (
	"encoding/json"
	"testing"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"github.com/stretchr/testify/require"
)

func TestGenerateUniqueID(t *testing.T) {
	t.Parallel()
	n := maelstrom.NewNode()
	n.Init("nodeID", []string{"peer1", "peer2"})
	repBody := GenerateUniqueID(n)
	require.Equal(t, "nodeID-1", repBody["id"], "bad id string")
	require.Equal(t, "generate_ok", repBody["type"], "bad message type")
	_, err := json.Marshal(repBody)
	require.NoError(t, err, "bad json")
}
