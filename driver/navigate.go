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

func (d *Driver) NewTab() {
	d.Client.ExecuteCmd(&Command{
		Path:   "/window/new",
		Method: http.MethodPost,
		Data:   marshalData(&Empty{}),
	})
}

func (d *Driver) SwitchToTab(n int) {
	h := getTabs(d)

	d.Client.ExecuteCmd(&Command{
		Path:   "/window",
		Method: http.MethodPost,
		Data:   marshalData(map[string]string{"handle": h[n]}),
	})
}

func getTabs(d *Driver) []string {
	h := new(struct{ Value []string })
	d.Client.ExecuteCmd(&Command{
		Path:   "/window/handles",
		Method: http.MethodGet,
	}, h)

	return h.Value
}
