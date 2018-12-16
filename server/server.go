package main
import "fmt"
import "net"
//import "os"
//import "bufio"


var (
        clients = make(map[string]net.Conn)
)

func main(){
        fmt.Println("Launching server...")
        // listen on all interfaces
        server, err := net.Listen("tcp", "localhost:1234")
        if err == nil {
                fmt.Println("Server listening on port 1234")
                fmt.Println()
        } else {
                checkError(err)
        }

        for {
                // accept connection on port
                conn, err := server.Accept()
                checkError(err)
                if (conn != nil) {
                        go openConnection(conn)
                }
        }

        defer server.Close()
}

func openConnection(connection net.Conn) {

        nickBuffer := make([]byte, 256)
        nickRawDataSize, _ := connection.Read(nickBuffer)
        nickRawData := nickBuffer[:nickRawDataSize]
        nick := string(nickRawData)

        nickUsed := "Nickname in use. Please enter a new nick: "
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
                messageRawDataSize, _ := connection.Read(messageBuffer)
                messageRawData := messageBuffer[:messageRawDataSize]
                message := string(messageRawData)
                
                dis_msg := fmt.Sprintln(nick, "disconnected")
                if message == disconnectCommand {
                        //dis_msg := nick + "disconnected"
                        dis_msg = fmt.Sprintln(nick, "disconnected")
                        fmt.Println(dis_msg)
                        delete(clients, nick)
                        connection.Close()
                        for _, client := range clients {
                               _, err := client.Write([]byte(dis_msg))
                                checkError(err)
                        }
                        return
                } else {
                        //response := nick + ": " + message
                        response := fmt.Sprintln(string(nick), ":", message)
                        for username, client := range clients {
                                if username !=nick {
                                        _, err := client.Write([]byte(response))
                                        checkError(err)
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