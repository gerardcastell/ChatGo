package main
import "fmt"
import "net"

var (
        clients = make(map[string]net.Conn)
)

func main(){
        fmt.Println()
        fmt.Println("Launching server...")
        // listen on all interfaces
        server, err := net.Listen("tcp", "localhost:1234")
        if err == nil {
                fmt.Println("Server listening on port 1234")
                fmt.Println()
        } else {
                checkError(err)
        }

        defer server.Close()

        for {
                // accept connection on port
                conn, err := server.Accept()
                checkError(err)
                if (conn != nil) {
                        go openConnection(conn)
                }
        }
}

func isNickInUse(nick string) (isUsed bool) {
	for foreignNick := range clients {
		if foreignNick == nick {
			return true
		}
	}
	return false
}

func openConnection(connection net.Conn) {

        nickBuffer := make([]byte, 256)
        nickRawDataSize, err := connection.Read(nickBuffer)
        checkErrorAndCloseConn(err, connection)
        if err != nil {
		return
	}
        nickRawData := nickBuffer[:nickRawDataSize]
        nick := string(nickRawData)

        isNickUsed := isNickInUse(nick)
        nickUsed := "Nickname in use. Please enter a new nick: "

        for isNickUsed {
		connection.Write([]byte(nickUsed))
		nickRawDataSize, _ := connection.Read(nickBuffer)
		nickRawData := nickBuffer[:nickRawDataSize]
		nick = string(nickRawData)
		isNickUsed = isNickInUse(nick)
	}
  

        for foreignNick := range clients {
                if foreignNick == nick {
                        connection.Write([]byte(nickUsed))
                        nickRawDataSize, _ := connection.Read(nickBuffer)
                        nickRawData := nickBuffer[:nickRawDataSize]
                        nick = string(nickRawData)
                        break
                }
        }
        newUserTxt := fmt.Sprintln("New user: ", nick)
        fmt.Print(newUserTxt)
        infoTxt := fmt.Sprintln(nick, "has joined the room")
        for _, clientConn := range clients {
                infoTxt = fmt.Sprintln(nick, "has joined the room")
                clientConn.Write([]byte(infoTxt))
        }

        clients[nick] = connection
        messageBuffer := make([]byte, 1000)
        disconnectCommand := "exit"

        for {
                messageRawDataSize, err := connection.Read(messageBuffer)
                checkErrorAndCloseConn(err, connection)
		if err != nil {
			delete(clients, nick)
			fmt.Println(nick + " disconnected")
			return
		}
                messageRawData := messageBuffer[:messageRawDataSize]
                message := string(messageRawData)
                
                dis_msg := fmt.Sprintln(nick, " disconnected")
                if message == disconnectCommand {
                        dis_msg = fmt.Sprintln(nick, " disconnected")
                        fmt.Println(dis_msg)
                        delete(clients, nick)
                        connection.Close()
                        for _, client := range clients {
                               _, err := client.Write([]byte(dis_msg))
                               checkErrorAndCloseConn(err, client)
                        }
                        return
                } else {
                        response := fmt.Sprintln(nick, ": ", message)
                        for username, client := range clients {
                                if username != nick {
                                        _, err := client.Write([]byte(response))
                                        checkErrorAndCloseConn(err, client)
                                }
                         }      
                }
        }
}

func checkError(err error) {
        if err != nil {
                fmt.Println("Fatal error ", err.Error())
        }
}

func checkErrorAndCloseConn(err error, conn net.Conn) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		if conn != nil {
			conn.Close()
		}
	}
}