package driver

import (
	"log"
	"net/http"
)

func (d Driver) PageSource() {
	s := new(struct{ Value string })

	st := defaultStrategy{
		Driver: d,
		Command: Command{
			Path:         "/source",
			Method:       http.MethodGet,
			ResponseData: s,
		},
	}

	st.execute()

	log.Println(string(s.Value))
}
