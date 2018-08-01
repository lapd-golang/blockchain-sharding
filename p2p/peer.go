/*
 * Copyright © 2018 Lynn <lynn9388@gmail.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package p2p

import (
	"net/rpc"
	"sync"
	"time"

	"github.com/lynn9388/blockchain-sharding/common"
)

type peer struct {
	node   *common.Node
	client *rpc.Client
}

const (
	maxPeerNum         = 4
	lackNodesSleepTime = 1
	fullNodesSleepTime = 2
)

var (
	peers          = make(map[string]peer)
	addPeerChan    = make(chan *common.Node)
	removePeerChan = make(chan *common.Node)
	peerMux        = sync.RWMutex{}
)

func NewPeerManager() {
	go connectPeers()

	for {
		select {
		case node := <-addPeerChan:
			addPeer(node)
		case node := <-removePeerChan:
			removePeer(node)
		}
	}
}

func addPeer(node *common.Node) {
	peerMux.Lock()
	defer peerMux.Unlock()
	if _, exists := peers[node.RPCAddr.String()]; !exists {
		client := ping(node)
		if client != nil {
			peers[node.RPCAddr.String()] = peer{node, client}
			common.Logger.Debugf("add new peer: %v", node.RPCAddr.String())
		}
	}
}

func removePeer(node *common.Node) {
	peerMux.Lock()
	defer peerMux.Unlock()
	peers[node.RPCAddr.String()].client.Close()
	delete(peers, node.RPCAddr.String())
}

// ping tests if a node is reachable and returns connected client
func ping(node *common.Node) *rpc.Client {
	ack := ""
	client, err := connectNode(node)
	if err != nil {
		return nil
	}
	err = client.Call("PingPongService.PingPong", pingMsg, &ack)
	if err != nil {
		common.Logger.Errorf("failed to call PingPong on %+v: %v", *node, err)
		return nil
	}
	if ack != pongMsg {
		common.Logger.Errorf("not a valid pong message: %v", ack)
		return nil
	}
	return client
}

func connectPeers() {
	for {
		peerMux.RLock()
		length := len(peers)
		peerMux.RUnlock()
		if length < maxPeerNum {
			shuffleNodes := getShuffleNodes()
			if len(*shuffleNodes) > maxPeerNum {
				*shuffleNodes = (*shuffleNodes)[:maxPeerNum]
			}
			for _, n := range *shuffleNodes {
				addPeerChan <- &n
			}
		}

		peerMux.RLock()
		length = len(peers)
		peerMux.RUnlock()
		if length < maxPeerNum {
			discoverSigChan <- true
			time.Sleep(lackNodesSleepTime * time.Second)
		} else {
			time.Sleep(fullNodesSleepTime * time.Second)
		}
	}
}