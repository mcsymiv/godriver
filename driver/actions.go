package driver

import "net/http"

type KeyAction struct {
	Type string `json:"type"`
	Key  string `json:"value"`
}

type KeyActionPayload struct {
	Actions []interface{} `json:"actions"`
}

const (
	NullInput    string = "null"
	KeyInput     string = "key"
	PointerInput string = "pointer"
	WheelInput   string = "wheel"
)

func (d *Driver) Action(act, keys string) {
	actions := make([]KeyAction, 0, len(keys))

	for _, key := range keys {
		actions = append(actions, KeyAction{
			Type: act,
			Key:  string(key),
		})
	}

	d.Client.ExecuteCmd(&Command{
		Path:   "/actions",
		Method: http.MethodPost,
		Data: marshalData(&KeyActionPayload{
			Actions: []interface{}{
				map[string]interface{}{
					"type":    "key",
					"id":      "default keyboard",
					"actions": actions,
				}},
		}),
	})
}
