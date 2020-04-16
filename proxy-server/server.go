// http://www.codingcookies.com/2013/09/21/creating-a-proxy-server-with-go/

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

var fromHost = flag.String("from", "localhost:80", "The proxy server's host.")
var toHost = flag.String("to", "localhost:8000", "The server that the proxy forwards to.")
var maxConnections = flag.Int("c", 25, "The maximum number of active connections.")
var maxWaitingConnections = flag.Int("cw", 10000, "The maximum numebr of waiting connections.")

func main() {
	flag.Parse()
	fmt.Printf("Proxying %s->%s\r\n", *fromHost, *toHost)

	// Create a TCP server
	server, err := net.Listen("tcp", *fromHost)

	if err != nil {
		log.Fatal(err)
	}

	// Create channels and fill
	waiting, spaces := make(chan net.Conn, *maxWaitingConnections), make(chan bool, *maxConnections)

	// Fill the channel to keep the function looping
	for i := 0; i < *maxConnections; i++ {
		spaces <- true
	}

	// Run matchConnections
	go matchConnections(waiting, spaces)

	for {
		// Take a connection
		connection, err := server.Accept()

		if err != nil {
			log.Print(err)
			continue
		}

		log.Printf("Received connection from %s.\r\n", connection.RemoteAddr())

		// Put the connection in the channel
		waiting <- connection
	}
}

func matchConnections(waiting chan net.Conn, spaces chan bool) {
	for connection := range waiting {
		// Wait until a space is available
		<-spaces

		// Run handle connection and refill space after
		go func(connection net.Conn) {
			handleConnection(connection)
			spaces <- true

			log.Printf("Closed connection from %s.\r\n", connection.RemoteAddr())
		}(connection)
	}
}

func handleConnection(connection net.Conn) {
	defer connection.Close()

	// Connect to forwarding address via tcp
	remote, err := net.Dial("tcp", *toHost)

	if err != nil {
		log.Fatal(err)
		return
	}

	defer remote.Close()

	complete, fromComplete, toComplete := make(chan bool, 2), make(chan bool, 1), make(chan bool, 1)

	// Send bits between forwarding address and host
	go copyContent(connection, remote, complete, fromComplete, toComplete)
	go copyContent(remote, connection, complete, fromComplete, toComplete)

	// Block until done
	<-complete
}

func copyContent(from net.Conn, to net.Conn, complete chan bool, fromComplete chan bool, toComplete chan bool) {
	var err error = nil
	var bytes []byte = make([]byte, 256)
	var read int = 0

	for {
		select {
		// External connection finished first
		case <-toComplete:
			complete <- true
			return
		// Transfer bytes from one stream to the other
		default:
			read, err = from.Read(bytes)

			if err != nil {
				complete <- true
				fromComplete <- true
				return
			}

			_, err = to.Write(bytes[:read])

			if err != nil {
				complete <- true
				fromComplete <- true
				return
			}
		}
	}
}

// func copyContent(from net.Conn, to net.Conn, complete chan bool) {
// 	var err error = nil
// 	var bytes []byte = make([]byte, 256)
// 	var read int = 0

// 	// Transfer bytes from one stream to another
// 	for {
// 		read, err = from.Read(bytes)

// 		if err != nil {
// 			complete <- true
// 			break
// 		}

// 		_, err = to.Write(bytes[:read])

// 		if err != nil {
// 			complete <- true
// 			break
// 		}
// 	}
// }
