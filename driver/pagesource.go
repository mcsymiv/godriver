package driver

import (
	"log"
	"net/http"
)

func (d Driver) PageSource() {
	s := new(struct{ Value string })
	d.Client.ExecuteCommand(&Command{
		Path:   "/source",
		Method: http.MethodGet,
	}, s)

	log.Println(string(s.Value))
}
