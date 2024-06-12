package driver

import (
	"log"
	"net/http"
)

func (d Driver) PageSource() {
	s := new(struct{ Value string })

	d.execute(defaultStrategy{Command{
		Path:         "/source",
		Method:       http.MethodGet,
		ResponseData: s,
	}})

	log.Println(string(s.Value))
}
