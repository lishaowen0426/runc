package comm

import (
	"testing"
	"time"
)

type testPayload struct {
	a string
	b string
}

func TestNewUnixSocket(t *testing.T) {
	defer Close()
	p := ConnectedPayload{
		Root: "/run/runc",
		Log:  "",
	}
	c := &Prot{
		Act:     CONNECTED,
		Payload: p,
	}

	if err := c.Send(); err != nil {
		t.Errorf("Send wrong: %s\n", err.Error())
	}

	time.Sleep(2 * time.Second)

}

func TestPayload(t *testing.T) {
	p := ConnectedPayload{
		Root: "/run/runc",
		Log:  "",
	}

	if m, err := ToMap(p); err != nil {
		t.Errorf("to map failed: %s", err.Error())
	} else if m["root"] != p.Root || m["log"] != p.Log {
		t.Errorf("to map wrong value, %s: %s, %s: %s", "root", m["root"], "log", m["log"])
	}

}
