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
	"net"
	"net/rpc"
	"sync"
	"time"

	"github.com/lynn9388/blockchain-sharding/common"
	"google.golang.org/grpc"
)

const (
	maxPeerNum         = 4
	lackNodesSleepTime = 1
	fullNodesSleepTime = 2
)

type peer struct {
	common.Node
	Conn *grpc.ClientConn
}

var (
	peers    = make(map[string]peer)
	peersMux = sync.RWMutex{}
)

func NewPeerManager() {
	go connectPeers()
}

func addPeer(p *peer) {
	peersMux.Lock()
	defer peersMux.Unlock()
	if _, exists := peers[p.RPCAddr.String()]; !exists {
		peers[p.RPCAddr.String()] = *p
		common.Logger.Debug("add new peer: ", p.RPCAddr.String())
	}
}

func removePeer(addr *net.TCPAddr) {
	peersMux.Lock()
	defer peersMux.Unlock()
	if _, exists := peers[addr.String()]; exists {
		peers[addr.String()].Conn.Close()
		delete(peers, addr.String())
		common.Logger.Debug("remove peer: ", addr.String())
	}
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
		peersMux.RLock()
		length := len(peers)
		peersMux.RUnlock()
		if length < maxPeerNum {
			shuffleNodes := getShuffleNodes()
			if len(shuffleNodes) > maxPeerNum {
				shuffleNodes = shuffleNodes[:maxPeerNum]
			}
			//for _, n := range shuffleNodes {
			//	addPeerChan <- &n
			//}
		}

		peersMux.RLock()
		length = len(peers)
		peersMux.RUnlock()
		if length < maxPeerNum {
			discoverSigChan <- true
			time.Sleep(lackNodesSleepTime * time.Second)
		} else {
			time.Sleep(fullNodesSleepTime * time.Second)
		}
	}
}
