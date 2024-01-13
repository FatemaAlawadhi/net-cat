package pkg
import (
	"bufio"
	"net"
	"sync"
)

type Client struct {
	Name   string
	Conn   net.Conn
	Writer *bufio.Writer
	Reader *bufio.Reader
}

var clientsMux sync.Mutex
var clients []*Client