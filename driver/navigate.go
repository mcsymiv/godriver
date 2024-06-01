package driver

import (
	"net/http"
)

func (d Driver) Url(u string) Driver {
	st := defaultStrategy{
		Driver: d,
		Command: Command{
			Path:   PathDriverUrl,
			Method: http.MethodPost,
			Data:   marshalData(map[string]string{"url": u}),
		},
	}

	st.execute()
	return d
}

func (d Driver) OpenInNewTab(u string) {
	d.NewTab()
	tNum := getTabs(d)
	d.Tab(len(tNum) - 1)
	d.Url(u)
}

func (d Driver) Refresh() {
	st := defaultStrategy{
		Driver: d,
		Command: Command{
			Path:   PathDriverRefresh,
			Method: http.MethodPost,
			Data:   marshalData(Empty{}),
		},
	}

	st.execute()
}

func (d Driver) NewTab() {
	st := defaultStrategy{
		Driver: d,
		Command: Command{
			Path:   PathDriverWindowNew,
			Method: http.MethodPost,
			Data:   marshalData(Empty{}),
		},
	}

	st.execute()
}

// SwitchToTab
// Swiches context to N tab in browser, where 0 is first tab
func (d Driver) Tab(n int) Driver {
	h := getTabs(d)

	st := defaultStrategy{
		Driver: d,
		Command: Command{
			Path:   PathDriverWindow,
			Method: http.MethodPost,
			Data:   marshalData(map[string]string{"handle": h[n]}),
		},
	}

	st.execute()
	return d
}

func getTabs(d Driver) []string {
	h := new(struct{ Value []string })

	st := defaultStrategy{
		Driver: d,
		Command: Command{
			Path:         PathDriverWindowHandles,
			Method:       http.MethodGet,
			ResponseData: h,
		},
	}

	st.execute()
	return h.Value
}
