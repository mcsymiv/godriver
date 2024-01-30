package driver

import (
	"log"
	"net/http"
)

type findStrategy struct {
	CommandHandler
}

func (f *findStrategy) Execute(req *http.Request) (*http.Response, error) {
	log.Printf("find strategy request: %s", req.URL.Path)

	response, err := f.CommandHandler(req)

	log.Printf("find response strategy status code: %v", response.StatusCode)

	return response, err
}

// openWrapper
// wraps W3C Navigate To command
func findWrapper(handler CommandHandler) CommandHandler {
	return func(req *http.Request) (*http.Response, error) {
		log.Printf("find request: %s", req.URL.Path)

		response, err := handler(req)

		log.Printf("find response status code: %v", response.StatusCode)

		return response, err
	}
}

func (d Driver) FindElement(selector string) (*Element, error) {
	jFind := &JsonFindUsing{
		Using: ByCssSelector,
		Value: selector,
	}

	op := &Command{
		Path:   "/element",
		Method: http.MethodPost,
		Data:   marshalData(jFind),
	}

	st := &findStrategy{
		CommandHandler: findWrapper(d.Client.Execute),
	}

	res, err := d.Client.ExecuteCommandStrategy(op, st)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	el := new(struct{ Value map[string]string })
	unmarshalData(res, &el)
	eId := elementID(el.Value)

	return &Element{
		Id:      eId,
		Client:  d.Client,
		Session: d.Session,
	}, nil

}
