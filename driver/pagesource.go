package driver

import "log"

func (d Driver) PageSource() {

	op := d.Commands["source"]

	res, err := d.Client.ExecuteCommandStrategy(op)
	if err != nil {
		log.Println("error on page sourse")
		return
	}
	s := new(struct{ Value string })
	unmarshalData(res, s)

	log.Println(string(s.Value))
}
