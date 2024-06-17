package driver

import (
	"net/http"
)

func newFrameCommand(el *Element) Command {
	var elFrameId map[string]string

	if el.Id == "" {
		elFrameId = el.ElementIdentifier()
	}

	return Command{
		Path:   PathDriverFrame,
		Method: http.MethodPost,
		Data: marshalData(map[string]interface{}{
			"id": elFrameId,
		}),
	}
}

func (el *Element) SwitchFrame() {
	el.Driver.execute(defaultStrategy{newFrameCommand(el)})
}

func (el *Element) SwitchFrameParent() {
	el.Driver.execute(defaultStrategy{newFrameCommand(nil)})
}
