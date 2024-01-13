package main

import (
	"fmt"
	"log"
	"net"
	"net-cat/pkg"
	"strconv"
	"os"
	"io/ioutil"
)

func main() {
	arg := os.Args
	var port string
	if len(arg) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	} else {
		if len(arg) == 1 {
			port = "8989"
		} else {
			port = arg[1]
			if !isValidPort(port) {
				log.Fatal("Invalid port number")
			}
		}
	}

	listener, err := net.Listen("tcp","localhost:" + port)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Listening on the port :" + port)

	//To clear the log when program excute
	err = ioutil.WriteFile("log.txt", []byte{}, 0644)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go pkg.HandleClient(conn)
	}

}


func isValidPort(port string) bool {
	portNumber, err := strconv.Atoi(port)
	if err != nil {
		return false
	}

	if portNumber >= 0 && portNumber <= 65535 {
		return true
	}

	return false
}
