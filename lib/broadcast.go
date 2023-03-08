package lib

import (
	"encoding/json"
	"errors"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type broadcastHistory map[int]struct{}

var (
	memMtx          sync.Mutex
	memory          = broadcastHistory{}
	errNoBcastValue = errors.New("no broadcast-value field")
	errNotNumber    = errors.New("broadcast value is not a number")
	bcastOK         = ReplyBody{
		"type": "broadcast_ok",
	}
	peersMtx     sync.Mutex
	peers        = []string{}
	errNoTopObj  = errors.New("no topology object in message")
	errBadTopObj = errors.New("malformed topology object in message")
	errNoPeerObj = errors.New("no peer object in topology message")
	errNoPeers   = errors.New("no peers specified in topology message")
	errBadPeers  = errors.New("malformed peer object in topology message")
)

func extractBroadcastValue(byt json.RawMessage) (int, error) {
	var bcast = ReplyBody{}
	if err := json.Unmarshal(byt, &bcast); err != nil {
		return 0, err
	}
	raw, ok := bcast["message"]
	if !ok {
		return 0, errNoBcastValue
	}
	val, ok := raw.(float64)
	if !ok {
		log.Printf("expected JSON number, got %v (%T)", raw, raw)
		return 0, errNotNumber
	}
	// The API is for integral values
	return int(val), nil
}

func NewBroadcastHandler(n *maelstrom.Node) maelstrom.HandlerFunc {
	return func(msg maelstrom.Message) error {
		val, err := extractBroadcastValue(msg.Body)
		if err != nil {
			log.Println("invalid broadcast value")
			return err
		}
		log.Printf("broadcast %v to peers", val)
		memMtx.Lock()
		memory[val] = struct{}{}
		memMtx.Unlock()
		peersMtx.Lock()
		// This protects us in case we got a broadcast message before a topology update
		// If the broadcast message comes in before the init message...welp.
		if len(peers) == 0 {
			peers = n.NodeIDs()
		}
		recipients := peers
		peersMtx.Unlock()
		for _, sendTo := range recipients {
			if sendTo == n.ID() {
				log.Printf("don't send to yourself, that's silly.")
				continue
			}
			// Send() does return errors but only on json marshaling and if we're here
			// we've already successfully unmarshaled it.
			log.Printf("sending to %s", sendTo)
			n.Send(sendTo, msg)
		}
		n.Reply(msg, bcastOK)
		return nil
	}
}

func NewReadHandler(n *maelstrom.Node) maelstrom.HandlerFunc {
	return func(msg maelstrom.Message) error {
		repBody := ReplyBody{
			"type": "read_ok",
		}
		memMtx.Lock()
		vals := make([]int, 0, len(memory))
		for v, _ := range memory {
			vals = append(vals, v)
		}
		memMtx.Unlock()
		repBody["messages"] = vals
		n.Reply(msg, repBody)
		return nil
	}
}

func extractPeersFromTopology(id string, byt json.RawMessage) ([]string, error) {
	var top = ReplyBody{}
	if err := json.Unmarshal(byt, &top); err != nil {
		return nil, err
	}
	topObj, ok := top["topology"]
	if !ok {
		return nil, errNoTopObj
	}
	topMap, ok := topObj.(map[string]any)
	if !ok {
		return nil, errBadTopObj
	}
	peerObj, ok := topMap[id]
	if !ok {
		return nil, errNoPeerObj
	}
	peerArray, ok := peerObj.([]any)
	if !ok {
		return nil, errBadPeers
	}
	if len(peerArray) == 0 {
		return nil, errNoPeers
	}
	newPeers := make([]string, 0, len(peerArray))
	for _, val := range peerArray {
		// skip anything that isn't a string
		if str, ok := val.(string); ok {
			newPeers = append(newPeers, str)
		}
	}
	return newPeers, nil
}

func NewTopologyHandler(n *maelstrom.Node) maelstrom.HandlerFunc {
	return func(msg maelstrom.Message) error {
		repBody := ReplyBody{
			"type": "topology_ok",
		}
		newPeers, err := extractPeersFromTopology(n.ID(), msg.Body)
		if err != nil {
			return err
		}
		peersMtx.Lock()
		peers = newPeers
		peersMtx.Unlock()
		n.Reply(msg, repBody)
		return nil
	}
}
