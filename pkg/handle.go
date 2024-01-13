package pkg

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

func HandleClient(conn net.Conn) {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	
	//Checking Connections Limit
	CheckLimit(writer,conn)

	// Get client name
	writer.WriteString(Logo())
	writer.Flush()
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}
	name = strings.TrimSpace(name)

	name, err = CheckUserNameLength(writer, reader, conn,name)
	name, err = IsNameTaken(writer, reader, conn,name)

	//Checking again in case somone entered his name faster than this user
	CheckLimit(writer,conn)

	client := &Client{
		Name:   name,
		Conn:   conn,
		Writer: writer,
		Reader: reader,
	}

	clientsMux.Lock()
	clients = append(clients, client)
	clientsMux.Unlock()
	// Send previous messages to the new client
	SendPreviousMessages(client)

	// Inform other clients about new client
	SendMessageToAllClients("\033[1;32m" + client.Name + " has joined our chat...\n" + "\033[0m",client.Name )
	go handleClientMessages(client)
}

func handleClientMessages(client *Client) {
	reader := bufio.NewReader(client.Conn)
	UserDisconnected := false
	for !UserDisconnected {
		prompt := fmt.Sprintf("[%s][%s]: ", time.Now().Format("2006-01-02 15:04:05"), client.Name)
		client.Conn.Write([]byte(prompt))

		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Printf("Client '%s' disconnected\n", client.Name)
				HandleClientDisconnect(client)
				UserDisconnected = true
			} else {
				log.Println(err)
			}
			continue
		}

		message = CheckFlag(client,message)
		if CheckeMessage(message) == true {
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			formattedMessage := fmt.Sprintf("[%s][%s]: %s", timestamp, client.Name, message)
			SendMessageToAllClients(formattedMessage, client.Name)
		}
	}
}

func CheckName(name string) bool {
	for _ , client := range clients {
		if client.Name == name {
			return false
		}
	}
	return true
}


func CheckFlag( client *Client , message string) string {
	if strings.Contains(message, "--") {
		message = strings.Replace(message, "--", "", 1)
		message = strings.TrimRight(message, "\n")
		SplitMessage := strings.Split(message," ")
		if SplitMessage[0] == "cu" {
			ActualUserName := client.Name
			UserName := strings.Join(SplitMessage[1:], " ")
			client.Name = ""
			UserName,err  = CheckUserNameLength(client.Writer , client.Reader , client.Conn, UserName)
			UserName,err = IsNameTaken(client.Writer , client.Reader , client.Conn, UserName)
			UserName = strings.TrimRight(UserName, "\n")
			client.Name = UserName
			message := "\033[1;34m" + ActualUserName + " has changed the Username to " + UserName + "\n" + "\033[0m"
			SendMessageToAllClients(message, client.Name)
			return ""

		} else if SplitMessage[0]== "h"{

			message := "\033[1;36m" + "Welcome " + client.Name + " to the net-cat ChatGroup!\n" +
			"Here are some commands that can enhance your chatting experience:\n\n" +
			"--cu NewUserName : Change Username\n--h : Help\n--b : Broadcast Important Message\n\n" +
			"Feel free to explore these commands and enjoy your time in the chat group!" + "\033[0m"
			client.Writer.WriteString("\n" + message + "\n")
			client.Writer.Flush()
			return ""

		} else if SplitMessage[0]== "b" {
			message = strings.Join(SplitMessage[1:], " ")
			return  "\033[1;35m" + message + "\033[0m" + "\n"

		}else {
			return message
		}
	} 
	return message
}

func CheckeMessage(message string) bool {
	//Don't Broadcast the message if it was empty
	trimmedMessage := strings.TrimSpace(message)
	if len(trimmedMessage) == 0 {
		return false
	}
	return true
}

func IsNameTaken( writer *bufio.Writer ,reader *bufio.Reader ,conn net.Conn, name string) (string,error) {
	for CheckName(name) == false {
		writer.WriteString("\033[1;31mUsername already taken. Please type a different Username: \033[0m")
		writer.Flush()
		name, err = reader.ReadString('\n')
		if err != nil {
			log.Println(err)
			conn.Close()
			return name,err
		}
		name = strings.TrimSpace(name)
	}
	return name, err
}

func CheckUserNameLength(writer *bufio.Writer ,reader *bufio.Reader ,conn net.Conn, name string) (string,error) {
	NameChar := []byte(name)
	for len(NameChar) < 1 || len(NameChar) > 17 {
		writer.WriteString("\033[1;31mUsername length should be in the range of 1-17 characters. Please type a different Username: \033[0m")
		writer.Flush()
		name, err = reader.ReadString('\n')
		if err != nil {
			log.Println(err)
			conn.Close()
			return name,err
		}
		name = strings.TrimSpace(name)
		NameChar = []byte(name)
	}
	return name,err
}

func CheckLimit(writer *bufio.Writer , conn net.Conn) {
	if len(clients) >= 10 {
		writer.WriteString("\033[1;31mConnection limit reached. Please try joining later...\033[0m")
		writer.Flush()
		conn.Close()
	}
}