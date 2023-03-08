package lib

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractBroadcastValue(t *testing.T) {
	t.Parallel()
	bcastNoMsg := json.RawMessage(
		`{
		"type": "broadcast"
	}`)
	bcastBadVal := json.RawMessage(
		`{
		"type": "broadcast",
		"message": "lolno"
	}`)
	bcastGood := json.RawMessage(
		`{
		"type": "broadcast",
		"message": 1024
	}`)

	val, err := extractBroadcastValue(bcastNoMsg)
	require.Error(t, err)
	require.Equal(t, errNoBcastValue, err)

	val, err = extractBroadcastValue(bcastBadVal)
	require.Error(t, err)
	require.Equal(t, errNotNumber, err)

	val, err = extractBroadcastValue(bcastGood)
	require.NoError(t, err)
	require.Equal(t, 1024, val)
}

func TestExtractPeersFromTopology(t *testing.T) {
	t.Parallel()
	noTopObj := json.RawMessage(
		`{
		"type": "topology"
	}`)
	badTopObj := json.RawMessage(
		`{
			"type": "topology",
			"topology": "lolno"
		}`)
	badPeerObj := json.RawMessage(
		`{
			"type": "topology",
			"topology": {
				"theNode": "lolno"
			}
		}`)
	noPeerObj := json.RawMessage(
		`{
			"type": "topology",
			"topology": {
			}
		}`)
	noPeers := json.RawMessage(
		`{
			"type": "topology",
			"topology": {
				"theNode": []
			}
		}`)
	good := json.RawMessage(
		`{
		"type": "topology",
		"topology": {
			"theNode": [ "theOtherNode", "theThirdNode" ]
		}
	}`)

	val, err := extractPeersFromTopology("theNode", noTopObj)
	require.Error(t, err)
	require.Equal(t, errNoTopObj, err)
	val, err = extractPeersFromTopology("theNode", badTopObj)
	require.Error(t, err)
	require.Equal(t, errBadTopObj, err)
	val, err = extractPeersFromTopology("theNode", badPeerObj)
	require.Error(t, err)
	require.Equal(t, errBadPeers, err)
	val, err = extractPeersFromTopology("theNode", noPeerObj)
	require.Error(t, err)
	require.Equal(t, errNoPeerObj, err)
	val, err = extractPeersFromTopology("theNode", noPeers)
	require.NoError(t, err)
	require.Equal(t, 0, len(val))
	val, err = extractPeersFromTopology("theNode", good)
	require.NoError(t, err)
	require.Equal(t, 2, len(val))
}
