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

package common

import (
	"net"
)

const (
	DefaultIP      = "127.0.0.1"
	DefaultAPIPort = 9388
	DefaultRPCPort = 9389
)

type (
	Config struct {
		IP      string `json:"ip" description:"ip address of the server" default:"127.0.0.1"`
		APIPort int    `json:"apiport" description:"port of the API Service" default:"9388"`
		RPCPort int    `json:"rpcport" description:"port of the RPC listener" default:"9389"`
	}

	// A node represents a potential peer on the network
	Node struct {
		RPCAddr net.TCPAddr
	}

	server struct {
		APIAddr net.TCPAddr
		Node    Node
	}
)

var (
	ServerConfig Config
	Server       server
)

func init() {
	ServerConfig = Config{IP: DefaultIP, APIPort: DefaultAPIPort, RPCPort: DefaultRPCPort}
	ConfigServer(&ServerConfig)
}

func ConfigServer(config *Config) {
	ServerConfig = *config

	apiAddr := net.TCPAddr{IP: net.ParseIP(config.IP), Port: config.APIPort}
	rpcAddr := net.TCPAddr{IP: net.ParseIP(config.IP), Port: config.RPCPort}
	Server = server{APIAddr: apiAddr, Node: Node{RPCAddr: rpcAddr}}
}