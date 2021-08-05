package subscraping

import (
	"context"
	"net/http"
)

// BasicAuth request's Authorization header
type BasicAuth struct {
	Username string
	Password string
}

// Source is an interface inherited by each passive source
type Source interface {
	// Run takes a domain as argument and a session object
	// which contains the extractor for subdomains, http client
	// and other stuff.
	Run(context.Context, string, *Session) <-chan Result
	// Name returns the name of the source
	Name() string
}

// Session is the option passed to the source, an option is created
// uniquely for eac source.
type Session struct {
	// Keys is the API keys for the application
	Keys *Keys
	// Client is the current http client
	Client *http.Client
}

// Keys contains the current API Keys we have in store
type Keys struct {
	Dnsgrep      string `json:"dnsgrep"`
	C99          string `json:"c99"`
	FofaUsername string `json:"fofa_username"`
	FofaSecret   string `json:"fofa_secret"`
	Shodan       string `json:"shodan"`
}

type Proxy string

// Result is a result structure returned by a source
type Result struct {
	Type   ResultType
	Source string //  API接口
	Value  string //  返回的数据
	Error  error
	//Count string
}

// ResultType is the type of result returned by the source
type ResultType int

// Types of results returned by the source
const (
	Subdomain ResultType = iota
	Error
)

// HostEntry defines a host with the source
type HostEntry struct {
	Host   string `json:"host"`
	Source string `json:"source"`
}
