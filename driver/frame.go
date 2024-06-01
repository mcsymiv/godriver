package driver

import (
	"net/http"
)

func newFrameCommand(el Element) Command {
	var elFrameId map[string]string

	if el == (Element{}) {
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

func (el Element) SwitchFrame() {
	st := defaultStrategy{
		Driver: el.Driver,
		Command: newFrameCommand(el),
	}

	st.execute()
}

func (el Element) SwitchFrameParent() {
	st := defaultStrategy{
		Driver: el.Driver,
		Command: newFrameCommand(Element{}),
	}

	st.execute()
}
