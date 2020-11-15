// Copyright 2021 The OCGI Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	sdk "github.com/ocgi/carrier-sdk/sdks/sdkgo"
	sdkapi "github.com/ocgi/carrier-sdk/sdks/sdkgo/api/v1alpha1"
)

var (
	Version = "default"
)

func main() {
	go doSignal()

	port := flag.String("port", "7654", "The port to listen to tcp traffic on")
	flag.Parse()
	if ep := os.Getenv("PORT"); ep != "" {
		port = &ep
	}

	log.Printf("Starting TCP server, listening on port %s", *port)
	ln, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("Could not start tcp server: %v", err)
	}
	defer ln.Close()

	log.Print("Creating SDK instance")
	s, err := sdk.NewSDK()
	if err != nil {
		log.Fatalf("Could not connect to sdk: %v", err)
	}

	log.Print("Marking this server as ready")
	if err := s.SetCondition("carrier.ocgi.dev/ready", "True"); err != nil {
		log.Fatalf("Could not set ready condition: %v", err)
	}

	log.Print("Starting watch gameserver")
	f := func(gs *sdkapi.GameServer) {
		// We just dump the GameServer fields.
		status, err := json.Marshal(gs)
		if err != nil {
			log.Printf("Decode gameserver error: %v", err)
			return
		}
		log.Printf("Dump GameServer: %s", string(status))
		if gs.Spec.Constraints == nil {
			return
		}
		for _, constraint := range gs.Spec.Constraints {
			if constraint == nil {
				continue
			}
			if constraint.Type == "NotInService" &&
				constraint.Effective {
				log.Printf("Fake wating for closing connnetion for %v", gs.ObjectMeta.Name)
				time.Sleep(30 * time.Second)
				if err = s.SetCondition("carrier.ocgi.dev/retired", "True"); err != nil {
					log.Printf("Failed to set retired: %v", err)
				}
				if err = s.SetCondition("carrier.ocgi.dev/has-no-player", "True"); err != nil {
					log.Printf("Failed to set hasplayer: %v", err)
				}
				os.Exit(0)
			}
		}
	}
	err = s.WatchGameServer(f)
	if err != nil {
		log.Fatalf("Failed to watch gameserver: %v", err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Unable to accept incoming tcp connection: %v", err)
		}
		go handleConnection(conn, s)
	}
}

// doSignal shutsdown on SIGTERM/SIGKILL
func doSignal() {
	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		close(stop)
	}()
	<-stop
	log.Println("Exit signal received. Shutting down.")
	os.Exit(0)
}

// handleConnection services a single tcp connection to the server
func handleConnection(conn net.Conn, s *sdk.SDK) {
	log.Printf("Client %s connected", conn.RemoteAddr().String())
	scanner := bufio.NewScanner(conn)
	for {
		if ok := scanner.Scan(); !ok {
			log.Printf("Client %s disconnected", conn.RemoteAddr().String())
			return
		}
		handleCommand(conn, scanner.Text(), s)
	}
}

// send response to client
func respond(conn net.Conn, txt string) {
	log.Printf("Responding with %q", txt)
	if _, err := conn.Write([]byte(txt + "\n")); err != nil {
		log.Fatalf("Could not write to tcp stream: %v", err)
	}
}

// handle client command request
func handleCommand(conn net.Conn, txt string, s *sdk.SDK) {
	parts := strings.Split(strings.TrimSpace(txt), " ")

	log.Printf("parts: %v", parts)
	cmd := parts[0]
	switch cmd {
	case "FILLED", "RETIRED", "HASPLAYER": // set command
		if len(parts) != 2 {
			log.Fatal("Bool value can not be empty")
			return
		}
		doSetRequest(conn, cmd, parts[1], s)
	case "EXIT":
		log.Printf("Receive EXIT command, exiting...")
		os.Exit(0)
	case "VERSION":
		respond(conn, "Version: "+Version)
	default:
		log.Printf("Invalid command: %s", cmd)
		respond(conn, "Invalid command: "+cmd)
	}
}

func doSetRequest(conn net.Conn, cmd, val string, sdk *sdk.SDK) {
	var err error

	v := func(v string) string {
		if v == "TRUE" {
			return "True"
		} else {
			return "False"
		}
	}(val)

	switch cmd {
	case "FILLED":
		if err = sdk.SetCondition("carrier.ocgi.dev/filled", v); err != nil {
			log.Printf("Failed to set filled: %v", err)
		}
	case "RETIRED":
		if err = sdk.SetCondition("carrier.ocgi.dev/retired", v); err != nil {
			log.Printf("Failed to set retired: %v", err)
		}
	case "HASNOPLAYER":
		if err = sdk.SetCondition("carrier.ocgi.dev/has-no-player", v); err != nil {
			log.Printf("Failed to set hasplayer: %v", err)
		}
	}
	if err != nil {
		respond(conn, "Failed to run command: "+cmd)
	} else {
		// send ACK to client
		respond(conn, "ACK: "+cmd)
	}
}
