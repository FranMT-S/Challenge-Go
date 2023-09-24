package mysocket

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

func Client() {
	var msg string
	var msgServer string
	var sw sync.WaitGroup
	connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		panic(err)
	}
	defer connection.Close() ///send some data
	fmt.Println("client connected")
	buffer := make([]byte, 1024)

	sw.Add(1)
	go func() {
		for {
			// fmt.Println("waiting server message:")
			mLen, err := connection.Read(buffer)
			if err != nil {
				fmt.Println("Error reading:", err.Error())
				sw.Done()
				return
			}

			msgServer = strings.ToLower(string(buffer[:mLen]))
			fmt.Println("Received: ", msgServer)
			if msgServer == "close" {
				fmt.Println("Saliendo")
				sw.Done()
				return
			}
		}

	}()

	go func() {
		for {
			fmt.Print("send message:")
			fmt.Scanln(&msg)
			_, err = connection.Write([]byte(msg))
			if err != nil {
				fmt.Println("message cannot read")
			}

		}
	}()

	sw.Wait()
}
