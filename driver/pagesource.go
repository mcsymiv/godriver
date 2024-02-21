package driver

import (
	"log"
	"net/http"
)

func (d Driver) PageSource() {
	op := &Command{
		Path:   "/source",
		Method: http.MethodGet,
	}

	bRes := d.Client.ExecuteCommand(op)
	s := new(struct{ Value string })
	unmarshalResponses(bRes, s)

	log.Println(string(s.Value))
}
