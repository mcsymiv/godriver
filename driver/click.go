package driver

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type clickStrategy struct {
	http.Client
}

func (cl clickStrategy) Execute(req *http.Request) (*http.Response, error) {
	log.Printf("click on: %s", req.URL.Path)
	return cl.Client.Do(req)
}

func (el *Element) Click() *Element {
	el.Client.ExecuteCmd(&Command{
		Path:           "/element/%s/click",
		PathFormatArgs: []any{el.Id},
		Method:         http.MethodPost,
		Data:           marshalData(&Empty{}),
		Strategies: []CommandExecutor{
			&clickStrategy{},
		},
	})

	return el
}

func (el *Element) ClickV2() *Element {

	el.Driver.Script(
		fmt.Sprintf(
			`
			function ev() {
				document.querySelector("%s").addEventListener("click", function(e) { 
					console.log("clicked");
					return "any value";
				});
			}
			return ev();
			`,
			el.Selector.Value,
		),
	)

	time.Sleep(2 * time.Second)

	el.Client.ExecuteCmd(&Command{
		Path:           "/element/%s/click",
		PathFormatArgs: []any{el.Id},
		Method:         http.MethodPost,
		Data:           marshalData(&Empty{}),
		Strategies: []CommandExecutor{
			&clickStrategy{},
		},
	})

	return el
}
