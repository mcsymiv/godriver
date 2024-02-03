package driver

import (
	"log"
)

func (d Driver) Open(u string) error {
	op := d.Commands["open"]
	op.Data = marshalData(map[string]string{"url": u})

	_, err := d.Client.ExecuteCommandStrategy(op)
	if err != nil {
		log.Println("error on open:", err)
		return err
	}

	return nil
}
