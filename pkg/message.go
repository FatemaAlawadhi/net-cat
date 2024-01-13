package pkg


import (
	"io/ioutil"
	"log"
	"os"
	"fmt"
	"strings"
	"time"
)

var err error
func SendMessageToAllClients(message string, name string) {
	clientsMux.Lock()
	defer clientsMux.Unlock()
	//true means that the message is Valid
	for _, c := range clients {
		if c.Name == name {
			continue
		} else {
			c.Writer.WriteString("\n" + message)
			c.Writer.Flush()
			prompt := fmt.Sprintf("[%s][%s]: ", time.Now().Format("2006-01-02 15:04:05"), c.Name)
			c.Conn.Write([]byte(prompt))
		}
	}
	AddToLog(message)
}

func CheckLogMessage(message string) bool {
	if strings.Contains(message, " has joined our chat...") || strings.Contains(message, " has left our chat...") || strings.Contains(message, " has changed the Username to "){
		return false
	} else {
		return true
	}
}

func AddToLog(message string) {
	if CheckLogMessage(message) == true {
		filePath := "log.txt"
		// Open the file in append mode
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		message = strings.TrimRight(message, "\n")
		// Append new message to the file
		_, err = fmt.Fprintln(file, message)
		if err != nil {
			log.Fatal(err)
		}	
	}
}

func ReadFromLog() string {
	content, err := ioutil.ReadFile("log.txt")
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

func SendPreviousMessages(client *Client) {
	clientsMux.Lock()
	defer clientsMux.Unlock()

	for _, c := range clients {
		if c == client {
			c.Writer.WriteString(ReadFromLog())
			c.Writer.Flush()
		}
	}
}

func HandleClientDisconnect(client *Client) {
	clientsMux.Lock()
	defer clientsMux.Unlock()
	message := "\033[1;31m" + client.Name + " has left our chat..." + "\033[0m\n"

	// Remove client from the list
	for i, c := range clients {
		if c == client {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}

	// Inform other clients about the client's exit
	SendDisconnectionMessage(message)
	client.Conn.Close()
}

func SendDisconnectionMessage(message string) {
	for _, c := range clients {
		c.Writer.WriteString("\n" + message)
		c.Writer.Flush()
		prompt := fmt.Sprintf("[%s][%s]: ", time.Now().Format("2006-01-02 15:04:05"), c.Name)
		c.Conn.Write([]byte(prompt))
	}
}