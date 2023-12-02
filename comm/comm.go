package comm

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

const SOCK_PATH = "/tmp/systemr_sock"

type Monitor struct {
	conn net.Conn
}

const (
	CONNECTED    = "Connected"
	DISCONNECTED = "Disconnected"
)

type Prot struct {
	Command string            `json:"command"`
	Payload map[string]string `json:"payload,omitempty"`
}

var mon Monitor

func init() {
	var err error
	mon, err = newUnixSocket()

	if err != nil {
		panic(err)
	}
}

func newUnixSocket() (Monitor, error) {

	conn, err := net.Dial("unix", SOCK_PATH)
	return Monitor{conn: conn}, err

}

func Close() error {
	return mon.conn.Close()
}

func Send(p Prot) (int, error) {
	var bytes []byte
	var sent int
	var err error

	bytes, err = json.Marshal(p)
	if err != nil {
		fmt.Printf("json marshal failed: %s\n", err.Error())
		return 0, err
	}
	fmt.Printf("marshall length: %d\n", len(bytes))
	fmt.Printf("cmd: %s\n", bytes)

	sent, err = mon.conn.Write(bytes)
	if err != nil {
		fmt.Printf("send failed: %s\n", err.Error())
		return 0, err
	} else if sent != len(bytes) {
		fmt.Printf("sent length wrong, buffer size: %d, sent size: %d\n", len(bytes), sent)
		return 0, errors.New("Wrong size")
	}

	return sent, nil
}
