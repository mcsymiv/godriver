package driver

import (
	"net/http"
)

func newFrameCommand(el *Element) *Command {
	var elFrameId map[string]string

	if el != nil {
		elFrameId = el.ElementIdentifier()
	}

	return &Command{
		Path:   PathDriverFrame, 
		Method: http.MethodPost,
		Data: marshalData(map[string]interface{}{
			"id": elFrameId,
		}),
	}
}

func (el *Element) SwitchFrame() {
	op := newFrameCommand(el)
	el.Client.ExecuteCmd(op)
}

func (el *Element) SwitchFrameParent() {
	op := newFrameCommand(nil)
	el.Client.ExecuteCmd(op)
}
