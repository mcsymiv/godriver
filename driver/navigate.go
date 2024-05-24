package driver

import (
	"net/http"
)

func (d *Driver) Url(u string) *Driver {
	d.Client.ExecuteCommand(&Command{
		Path:   PathDriverUrl,
		Method: http.MethodPost,
		Data:   marshalData(map[string]string{"url": u}),
	}, nil)

	return d
}

func (d Driver) Open(u string) {
	d.Client.ExecuteCommand(&Command{
		Path:   PathDriverUrl,
		Method: http.MethodPost,
		Data:   marshalData(map[string]string{"url": u}),
	}, nil)
}

func (d *Driver) OpenInNewTab(u string) {
	d.NewTab()
	tNum := getTabs(d)
	d.Tab(len(tNum) - 1)
	d.Open(u)
}

func (d Driver) Refresh() {
	d.Client.ExecuteCommand(&Command{
		Path:   PathDriverRefresh,
		Method: http.MethodPost,
		Data:   marshalData(&Empty{}),
	}, nil)
}

func (d *Driver) NewTab() {
	d.Client.ExecuteCommand(&Command{
		Path:   PathDriverWindowNew,
		Method: http.MethodPost,
		Data:   marshalData(&Empty{}),
	}, nil)
}

// SwitchToTab
// Swiches context to N tab in browser, where 0 is first tab
func (d *Driver) Tab(n int) *Driver {
	h := getTabs(d)

	d.Client.ExecuteCommand(&Command{
		Path:   PathDriverWindow,
		Method: http.MethodPost,
		Data:   marshalData(map[string]string{"handle": h[n]}),
	}, nil)

	return d
}

func getTabs(d *Driver) []string {
	h := new(struct{ Value []string })

	d.Client.ExecuteCommand(&Command{
		Path:   PathDriverWindowHandles,
		Method: http.MethodGet,
	}, h)

	return h.Value
}
