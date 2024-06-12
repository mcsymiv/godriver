package driver

// Command
// represents a request (command) to webdriver
// with added Strategies to execute in CommandExecutor
type Command struct {
	Path           string
	Method         string
	PathFormatArgs []any
	Data           []byte

	ResponseData interface{}
}
