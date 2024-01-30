package driver

import (
	"log"
	"net/http"
)

func (d Driver) Open(u string) error {
	op := &Command{
		Path:   "/url",
		Method: http.MethodPost,
		Data:   marshalData(map[string]string{"url": u}),
	}

	_, err := d.Client.ExecuteCommandStrategy(op)
	if err != nil {
		log.Println("error on open:", err)
		return err
	}

	return nil
}
