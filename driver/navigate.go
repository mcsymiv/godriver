package driver

import (
	"fmt"
	"net/http"
)

func (d *Driver) Url(u string) error {
	_, err := d.Client.ExecuteCmd(&Command{
		Path:   "/url",
		Method: http.MethodPost,
		Data:   marshalData(map[string]string{"url": u}),
	})

	if err != nil {
		return fmt.Errorf("error on Url: %v", err)
	}

	return nil
}

func (d Driver) Open(u string) {
	d.Client.ExecuteCmd(&Command{
		Path:   "/url",
		Method: http.MethodPost,
		Data:   marshalData(map[string]string{"url": u}),
	})
}

func (d *Driver) OpenInNewTab(u string) {
	d.NewTab()
	tNum := getTabs(d)
	d.SwitchToTab(len(tNum) - 1)
	d.Open(u)
}

func (d Driver) Refresh() {
	d.Client.ExecuteCmd(&Command{
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

// SwitchToTab
// Swiches context to N tab in browser, where 0 is first tab
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
