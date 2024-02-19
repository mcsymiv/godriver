package driver

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// reusableReader
// allows perform multiple requests in handlers
type reusableReader struct {
	io.Reader
	readBuf *bytes.Buffer
	backBuf *bytes.Buffer
}

func ReusableReader(r io.Reader) io.Reader {
	readBuf := bytes.Buffer{}
	_, err := readBuf.ReadFrom(r)
	if err != nil {
		log.Println("error on reusable reader buffer:", err)
		return nil
	}

	backBuf := bytes.Buffer{}

	return reusableReader{
		io.TeeReader(&readBuf, &backBuf),
		&readBuf,
		&backBuf,
	}
}

func (r reusableReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	if err == io.EOF {
		r.reset()
	}
	return n, err
}

func (r reusableReader) reset() {
	io.Copy(r.readBuf, r.backBuf)
}

func marshalData(body interface{}) []byte {
	b, err := json.Marshal(body)
	if err != nil {
		log.Println("error on marshal: ", err)
		return nil
	}

	return b
}

func unmarshalData(res *http.Response, any interface{}) []byte {
	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("error on reading response:", err)
		return nil
	}

	if err := json.Unmarshal(b, &any); err != nil {
		log.Println("error on unmarshal:", err)
		return nil
	}

	return b
}

// TODO: consider this unmarshal as oppose to above func
func unmarshalResponse(rData []byte, any interface{}) {
	if err := json.Unmarshal(rData, &any); err != nil {
		log.Println("error on unmarshal response:", err)
	}
}

func unmarshalResponses(buffRes []*buffResponse, any ...interface{}) {
	for i, res := range buffRes {
		if err := json.Unmarshal(res.buff, &any[i]); err != nil {
			log.Printf("error on unmarshal %d response: %v", i, err)
		}
	}
}
