package driver

import (
	"log"
	"net/http"
)

func (el *Element) SwitchFrame() error {
	op := &Command{
		Path:   "/frame",
		Method: http.MethodPost,
		Data: marshalData(map[string]interface{}{
			"id": el.ElementIdentifier(),
		}),
	}

	_, err := el.Client.ExecuteCommandStrategy(op)
	if err != nil {
		log.Println("Switch frame request error", err)
		return err
	}

	return nil
}

func (el *Element) SwitchFrameParent() error {
	op := &Command{
		Path:   "/frame",
		Method: http.MethodPost,
		Data: marshalData(map[string]interface{}{
			"id": nil,
		}),
	}

	_, err := el.Client.ExecuteCommandStrategy(op)
	if err != nil {
		log.Println("Switch frame request error", err)
		return err
	}

	return nil
}
