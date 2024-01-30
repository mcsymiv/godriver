package capabilities

type Capabilities struct {
	DriverSetupCapabilities
	Capabilities BrowserCapabilities `json:"capabilities"`
}

type DriverSetupCapabilities struct {
	Port     string `json:"-"`
	Host     string `json:"-"`
	Protocol string `json:"-"`
}

type BrowserCapabilities struct {
	AlwaysMatch `json:"alwaysMatch"`
}

type AlwaysMatch struct {
	AcceptInsecureCerts bool   `json:"acceptInsecureCerts"`
	BrowserName         string `json:"browserName"`
	Timeouts            `json:"timeouts,omitempty"`
	MozOptions          `json:"moz:firefoxOptions,omitempty"`
}

type Timeouts struct {
	Implicit float32 `json:"implicit,omitempty"`
	PageLoad float32 `json:"pageLoad,omitempty"`
	Script   float32 `json:"script,omitempty"`
}

type MozOptions struct {
	Profile string   `json:"profile,omitempty"`
	Binary  string   `json:"binary,omitempty"`
	Args    []string `json:"args,omitempty"`
	Log     `json:"log,omitempty"`
}

type Log struct {
	Level string `json:"level,omitempty"`
}

type CapabilitiesFunc func(*Capabilities)

// DefaultCapabilities
// Sets default firefox browser with local dev url
// With defined in service port, i.e. :4444
// Port and Host fields are used and passed to the WebDriver instance
// To reference and build current driver url
func DefaultCapabilities() Capabilities {
	return Capabilities{

		DriverSetupCapabilities: DriverSetupCapabilities{
			Port:     "4444",
			Host:     "localhost",
			Protocol: "http",
		},

		Capabilities: BrowserCapabilities{
			AlwaysMatch{
				AcceptInsecureCerts: true,
				BrowserName:         "firefox",
			},
		},
	}
}

func ImplicitWait(w float32) CapabilitiesFunc {
	return func(cap *Capabilities) {
		cap.Capabilities.AlwaysMatch.Timeouts.Implicit = w
	}
}

func Firefox(moz *MozOptions) CapabilitiesFunc {
	return func(cap *Capabilities) {
		cap.Capabilities.AlwaysMatch.MozOptions = *moz
	}
}

func BrowserName(b string) CapabilitiesFunc {
	return func(cap *Capabilities) {
		cap.Capabilities.AlwaysMatch.BrowserName = b
	}
}

func Port(p string) CapabilitiesFunc {
	return func(caps *Capabilities) {
		caps.DriverSetupCapabilities.Port = p
	}
}

func Host(h string) CapabilitiesFunc {
	return func(caps *Capabilities) {
		caps.DriverSetupCapabilities.Host = h
	}
}
