package driver

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
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

func unmarshalResponses(buffRes []*buffResponse, any ...interface{}) {
	if len(buffRes) > 0 {
		for i, res := range buffRes {
			err := json.Unmarshal(res.buff, &any[i])

			if err != nil {
				log.Printf("error on unmarshal %d response: %v", i, err)
			}
		}
	}
}
