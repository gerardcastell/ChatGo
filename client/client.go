package main

import (
	"os"
 	"bufio"
 	"fmt"
	"net"
	"sync"
)

var wg sync.WaitGroup
// Tell the 'wg' WaitGroup how many threads/goroutines
//   that are about to run concurrently.

func main(){
	wg.Add(2)
	nick := newUser()

	// connect to this socket
	connection, err := net.Dial("tcp", "127.0.0.1:1234")
	checkError(err, "")
	user := []byte(nick)
	_, _ = connection.Write(user)

	go sender(connection, nick)
	go receiver(connection)
	wg.Wait()
}

func newUser ()(nick string) {
	fmt.Println("Introduce tu nick:")
	reader := bufio.NewReader(os.Stdin)
	input, err_input := reader.ReadString('\n')
	checkError(err_input, "")
	nick = input[0:len(input)-1]
	return nick
}

func sender(connection net.Conn, nick string){
	reader := bufio.NewReader(os.Stdin)
	message, _, _ := reader.ReadLine()

	for message != nil{
		_, err := connection.Write(message)
		checkError(err, "")
		message, _, _ = reader.ReadLine()
	}

	wg.Done()
}

func receiver(connection net.Conn){
	buf := make([]byte, 1000)
	input, err := connection.Read(buf)
	message := buf[:input]

	for message != nil {
		if err != nil {
			fmt.Println("Connection closed...")
			return
		}
		fmt.Print(string(message))
		input, err = connection.Read(buf)
		message = buf[:input]
	} 
	wg.Done()
}

func checkError(err error, txt string){
	message := "Fatal error: "
	if txt == ""{
		message = txt
	}
	if err != nil {
			fmt.Println(message, err.Error())
	}
}