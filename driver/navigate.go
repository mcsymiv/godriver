package driver

import (
	"net/http"
)

func (d *Driver) Url(u string) *Driver {
	_, err := d.Client.ExecuteCmd(&Command{
		Path:   PathDriverUrl,
		Method: http.MethodPost,
		Data:   marshalData(map[string]string{"url": u}),
	})

	if err != nil {
		return nil
	}

	return d
}

func (d Driver) Open(u string) {
	d.Client.ExecuteCmd(&Command{
		Path:   PathDriverUrl,
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
		Path:   PathDriverRefresh,
		Method: http.MethodPost,
		Data:   marshalData(&Empty{}),
	})
}

func (d *Driver) NewTab() {
	d.Client.ExecuteCmd(&Command{
		Path:   PathDriverWindowNew,
		Method: http.MethodPost,
		Data:   marshalData(&Empty{}),
	})
}

// SwitchToTab
// Swiches context to N tab in browser, where 0 is first tab
func (d *Driver) SwitchToTab(n int) {
	h := getTabs(d)

	d.Client.ExecuteCmd(&Command{
		Path:   PathDriverWindow,
		Method: http.MethodPost,
		Data:   marshalData(map[string]string{"handle": h[n]}),
	})
}

func (d *Driver) Tab(n int) *Driver {
	h := getTabs(d)

	_, err := d.Client.ExecuteCmd(&Command{
		Path:   PathDriverWindow,
		Method: http.MethodPost,
		Data:   marshalData(map[string]string{"handle": h[n]}),
	})

	if err != nil {
		return nil
	}

	return d
}

func getTabs(d *Driver) []string {
	h := new(struct{ Value []string })

	_, err := d.Client.ExecuteCmd(&Command{
		Path:  PathDriverWindowHandles,
		Method: http.MethodGet,
	}, h)

	if err != nil {
		return nil
	}

	return h.Value
}
