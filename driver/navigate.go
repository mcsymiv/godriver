package driver

import (
	"net/http"
)

func (d Driver) Open(u string) {
	d.Client.ExecuteCmd(&Command{
		Path:   "/url",
		Method: http.MethodPost,
		Data:   marshalData(map[string]string{"url": u}),
	})
}

func (d Driver) Refresh() {
	d.Client.ExecuteCommand(&Command{
		Path:   "/refresh",
		Method: http.MethodPost,
		Data:   marshalData(&Empty{}),
	})
}
