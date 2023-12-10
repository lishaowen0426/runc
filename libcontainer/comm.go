package libcontainer

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"net"
	"os"
	"strconv"
	"time"
)

const SockPath = "/tmp/systemr_sock"
const ServerAddr = "127.0.0.1:8083"
const Server = "127.0.0.1"

type Monitor struct {
	conn net.Conn
}

type Prot struct {
	Act     Action  `json:"action"`
	Payload Payload `json:"payload,omitempty"`
}

var mon Monitor

func CommConnect() {
	var err error
	mon, err = newUDPConnection()

	if err != nil {
		if oerr, ok := err.(*net.OpError); ok {
			if serr, ook := oerr.Err.(*os.SyscallError); ook && serr.Err.Error() != "connection refused" {
				logrus.Warnf("Connect to systemr failed with %s", serr.Error())
			}

		} else {
			logrus.Warnf("Connect to systemr failed with %s", err.Error())
		}

	}
}

func newUDPConnection() (Monitor, error) {
	ack := make([]byte, 3)
	var conn net.Conn
	var err error
	if conn, err = net.Dial("udp", ServerAddr); err != nil {
		return Monitor{conn: nil}, err
	} else {

		if udp, ok := conn.(*net.UDPConn); !ok {
			return Monitor{}, errors.New("not an udp")
		} else {

			time.Sleep(1 * time.Second)

			udp.Write(ack)
			// wait port number
			if n, addr, uerr := udp.ReadFromUDP(ack); uerr != nil {

				return Monitor{}, uerr
			} else if addr.String() != ServerAddr || n != 2 {
				return Monitor{}, uerr
			} else {

				port := binary.BigEndian.Uint16(ack)
				commAddr := Server + ":" + strconv.Itoa(int(port))
				if conn, err = net.Dial("udp", commAddr); err != nil {
					return Monitor{}, err
				}

				//send ack
				ack = []byte("ack")
				if n, err = conn.Write(ack); err != nil || n != len(ack) {
					return Monitor{}, err
				}
				return Monitor{conn: conn}, nil

			}

		}

	}
}

func CommClose() error {
	var err error
	if mon.conn != nil {
		if err = SendDisconnected(); err != nil {
			return err
		}
		return mon.conn.Close()

	}
	return nil

}

func (p *Prot) Send() error {

	if mon.conn == nil {
		return nil
	}

	var bytes []byte
	var sent int
	var err error

	bytes, err = json.Marshal(p)
	if err != nil {
		logrus.Warnf("json marshal failed: %s", err.Error())
		return err
	}

	sent, err = mon.conn.Write(bytes)
	if err != nil {
		logrus.Warnf("send failed: %s", err.Error())
		return err
	} else if sent != len(bytes) {
		logrus.Warnf("sent length wrong, buffer size: %d, sent size: %d", len(bytes), sent)
		return errors.New("wrong size")
	}

	logrus.Infof("Send to monitor: %s", string(bytes))

	return nil
}

func GetFile() *os.File {

	if udp, ok := mon.conn.(*net.UDPConn); ok {
		if file, err := udp.File(); err != nil {
			return nil
		} else {
			return file
		}
	} else {
		return nil
	}
}
