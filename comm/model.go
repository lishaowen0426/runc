package comm

import (
	"encoding/json"

	specs "github.com/opencontainers/runc/specs-go"
)

type Action string

type Payload interface {
}

const (
	CONNECTED    Action = "Connected"
	DISCONNECTED Action = "Disconnected"
	SPEC         Action = "Spec"
	NEWPARENT    Action = "Newparent"
)

func ToMap(p Payload) (map[string]any, error) {
	var err error
	var jsonbody []byte
	var result map[string]any
	if jsonbody, err = json.Marshal(&p); err != nil {
		return nil, err
	} else {
		if err = json.Unmarshal(jsonbody, &result); err != nil {
			return nil, err
		}

	}
	return result, nil

}
func FromMap(dict map[string]any, result *Payload) error {
	var err error
	var jsonbody []byte
	if jsonbody, err = json.Marshal(dict); err != nil {
		return err
	} else {
		if err = json.Unmarshal(jsonbody, result); err != nil {
			return err
		}

	}
	return nil
}

type ConnectedPayload struct {
	Root string `json:"root"`
	Log  string `json:"log"`
}

func SendConnected(root, log string) error {
	p := &Prot{
		Act: CONNECTED,
		Payload: ConnectedPayload{
			Root: root,
			Log:  log,
		},
	}

	return p.Send()

}

func SendDisconnected() error {
	p := &Prot{
		Act: DISCONNECTED,
	}
	return p.Send()
}

type SpecPayload struct {
	Spec specs.Spec `json:"spec"`
}

func SendSpec(s *specs.Spec) error {
	p := &Prot{
		Act:     SPEC,
		Payload: SpecPayload{Spec: *s},
	}

	return p.Send()
}

type file struct {
	Pid  uint   `json:"pid"`
	Name string `json:"name"`
}

type processComm struct {
	InitParent file `json:"init_parent"`
	InitChild  file `json:"init_child"`

	SyncParent file `json:"sync_parent"`
	SyncChild  file `json:"sync_child"`

	LogParent file `json:"log_parent"`
	LogÇhild  file `json:"log_çhild"`
}

type ConfigPaylog struct {
}
