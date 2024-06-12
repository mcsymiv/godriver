package driver

import (
	"encoding/json"
	"log"
	"net/http"
)

func marshalData(body interface{}) []byte {
	b, err := json.Marshal(body)
	if err != nil {
		log.Println("error on marshal: ", err)
		return nil
	}

	return b
}

func unmarshalRes(res *http.Response, any interface{}) error {
	return json.NewDecoder(res.Body).Decode(any)
}
