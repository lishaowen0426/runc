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
	p :=
		Prot{
			Command: "connected",
			Payload: map[string]string{
				"firstname": "a",
				"lastname":  "b",
			},
		}
	time.Sleep(10 * time.Second)

	if sent, err := Send(p); err != nil {
		t.Errorf("Send wrong: %s\n", err.Error())
	} else {
		t.Logf("Sent: %d\n", sent)
	}

}
