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

package server

import (
	"errors"
	"net"
	"net/rpc"
)

type PingPongService int

const (
	pingMsg = "PING"
	pongMsg = "PONG"
)

func newRPCListener(addr *net.TCPAddr) {
	rpc.Register(new(PingPongService))
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		logger.Fatalf("failed to start RPC listener: %v", err)
	}
	logger.Infof("start RPC listener on %v", addr.String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("failed to accept a RPC connection")
			continue
		}
		go rpc.ServeConn(conn)
	}
}

// PingPong send pong ack message for ping message
func (t *PingPongService) PingPong(msg *string, ack *string) error {
	if *msg != pingMsg {
		return errors.New("not a valid ping message: " + *msg)
	}
	*ack = pongMsg
	return nil
}