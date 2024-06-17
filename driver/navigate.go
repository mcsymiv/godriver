package driver

import (
	"net/http"
)

func (d *Driver) Url(u string) *Driver {
	d.execute(defaultStrategy{Command{
		Path:   PathDriverUrl,
		Method: http.MethodPost,
		Data:   marshalData(map[string]string{"url": u}),
	}})

	return d
}

func (d *Driver) OpenInNewTab(u string) {
	d.NewTab()
	tNum := getTabs(d)
	d.Tab(len(tNum) - 1)
	d.Url(u)
}

func (d *Driver) Refresh() {
	d.execute(defaultStrategy{Command{
		Path:   PathDriverRefresh,
		Method: http.MethodPost,
		Data:   marshalData(Empty{}),
	}})
}

func (d *Driver) NewTab() {
	d.execute(defaultStrategy{Command{
		Path:   PathDriverWindowNew,
		Method: http.MethodPost,
		Data:   marshalData(Empty{}),
	}})
}

// SwitchToTab
// Swiches context to N tab in browser, where 0 is first tab
func (d *Driver) Tab(n int) *Driver {
	h := getTabs(d)

	d.execute(defaultStrategy{Command{
		Path:   PathDriverWindow,
		Method: http.MethodPost,
		Data:   marshalData(map[string]string{"handle": h[n]}),
	}})

	return d
}

func getTabs(d *Driver) []string {
	h := new(struct{ Value []string })

	d.execute(defaultStrategy{Command{
		Path:         PathDriverWindowHandles,
		Method:       http.MethodGet,
		ResponseData: h,
	}})

	return h.Value
}
