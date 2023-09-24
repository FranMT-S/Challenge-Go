package mysocket

import (
	"fmt"
	"net"
	"os"
	"sync"
)

func Server() {
	fmt.Println("Server Running...")
	server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer server.Close()
	fmt.Println("Listening on " + SERVER_HOST + ":" + SERVER_PORT)
	fmt.Println("Waiting for client...")
	for {
		connection, err := server.Accept()

		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("client connected")
		go processClient(connection)
	}
}

func processClient(connection net.Conn) {
	var msgClient string
	var sw sync.WaitGroup
	buffer := make([]byte, 1024)

	go func() {
		for {
			// fmt.Println("waiting client message")
			mLen, err := connection.Read(buffer)
			if err != nil {
				fmt.Println("Error reading:", err.Error())
				connection.Close()
				return
			}
			msgClient = string(buffer[:mLen])
			fmt.Println("\rReceived: ", msgClient)
			fmt.Println("command to client:")
		}
	}()
	sw.Add(1)
	go func() {
		for {
			fmt.Println("command to client:")
			fmt.Scan(&msgClient)
			_, err := connection.Write([]byte(msgClient))
			if err != nil {
				fmt.Println("cannot send message to client")
				connection.Close()
				sw.Done()
				return
			}

		}
	}()
	sw.Wait()

}
