package driver

import (
	"net/http"
)

type KeyAction struct {
	Type string `json:"type"`
	Key  string `json:"value"`
}

type WheelAction struct {
	Type     string `json:"type"`
	Duration uint   `json:"duration"` // uint(time.Duration duration / time.Millisecond)

	// PointerMoveOrigin controls how the offset for
	// the pointer move action is calculated.
	// type PointerMoveOrigin string

	// const (
	// FromViewport calculates the offset from the viewport at 0,0.
	// FromViewport PointerMoveOrigin = "viewport"
	// FromPointer calculates the offset from the current pointer position.
	// FromPointer = "pointer"
	// )
	Origin string `json:"origin"` // If origin is undefined let origin equal "viewport". src: https://w3c.github.io/webdriver/#wheel-input-source
	X      int    `json:"x"`
	Y      int    `json:"y"`
	DeltaX int    `json:"deltaX"` // deltaX":null,"deltaY":null,
	DeltaY int    `json:"deltaY"`
}

type ActionPayload struct {
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
		Data: marshalData(&ActionPayload{
			Actions: []interface{}{
				map[string]interface{}{
					"type":    "key",
					"id":      "default keyboard",
					"actions": actions,
				}},
		}),
	})
}

func (d *Driver) WheelAction() {
	actions := make([]WheelAction, 0)

	actions = append(actions, WheelAction{
		Type:     "scroll",
		Duration: uint(500),
		Origin:   "viewport",
		X:        100,
		Y:        100,
		DeltaX:   0,
		DeltaY:   500,
	})

	d.Client.ExecuteCmd(&Command{
		Path:   "/actions",
		Method: http.MethodPost,
		Data: marshalData(&ActionPayload{
			Actions: []interface{}{
				map[string]interface{}{
					"type":    "wheel",
					"id":      "default wheel",
					"actions": actions,
				}},
		}),
	})
}
