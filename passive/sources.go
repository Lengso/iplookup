package passive

import (
	"iplookup/subscraping"
	"iplookup/subscraping/sources/aizhan"
	"iplookup/subscraping/sources/bugscaner"
	"iplookup/subscraping/sources/c99"
	"iplookup/subscraping/sources/chinaz"
	"iplookup/subscraping/sources/dnsgrep"
	"iplookup/subscraping/sources/dnslytics"
	"iplookup/subscraping/sources/domaintools"
	"iplookup/subscraping/sources/hackertarget"
	"iplookup/subscraping/sources/ip138"
	"iplookup/subscraping/sources/omnisint"
	"iplookup/subscraping/sources/rapiddns"
	"iplookup/subscraping/sources/securitytrails"
	"iplookup/subscraping/sources/viewdns"
	"iplookup/subscraping/sources/webscan"
	"iplookup/subscraping/sources/yougetsignal"
)

// DefaultSources contains the list of fast sources used by default.
var DefaultSources = []string{
	"webscan",
	"rapiddns",
	"ip138",
	"yougetsignal",
	"aizhan",
	"chinaz",
	"viewdns",
	"bugscaner",
	"hackertarget",
	"dnslytics",
	"omnisint",
	"dnsgrep",
	"domaintools",
	"securitytrails",
}

// DefaultAllSources contains list of all sources
var DefaultAllSources = []string{
	"webscan",
	"rapiddns",
	"ip138",
	"yougetsignal",
	"aizhan",
	"c99",
	"chinaz",
	"viewdns",
	"bugscaner",
	"hackertarget",
	"dnslytics",
	"omnisint",
	"dnsgrep",
	"domaintools",
	"securitytrails",
}

// Agent is a struct for running passive subdomain enumeration
// against a given host. It wraps subscraping package and provides
// a layer to build upon.
type Agent struct {
	sources map[string]subscraping.Source
}

// New creates a new agent for passive subdomain discovery
func New(sources, exclusions []string) *Agent {
	// Create the agent, insert the sources and remove the excluded sources
	agent := &Agent{sources: make(map[string]subscraping.Source)}

	agent.addSources(sources)
	agent.removeSources(exclusions)

	return agent
}

// addSources adds the given list of sources to the source array
func (a *Agent) addSources(sources []string) {
	for _, source := range sources {
		switch source {
		case "webscan":
			a.sources[source] = &webscan.Source{}
		case "hackertarget":
			a.sources[source] = &hackertarget.Source{}
		case "dnsgrep":
			a.sources[source] = &dnsgrep.Source{}
		case "rapiddns":
			a.sources[source] = &rapiddns.Source{}
		case "c99":
			a.sources[source] = &c99.Source{}
		case "ip138":
			a.sources[source] = &ip138.Source{}
		case "aizhan":
			a.sources[source] = &aizhan.Source{}
		case "omnisint":
			a.sources[source] = &omnisint.Source{}
		case "viewdns":
			a.sources[source] = &viewdns.Source{}
		case "bugscaner":
			a.sources[source] = &bugscaner.Source{}
		case "dnslytics":
			a.sources[source] = &dnslytics.Source{}
		case "domaintools":
			a.sources[source] = &domaintools.Source{}
		case "yougetsignal":
			a.sources[source] = &yougetsignal.Source{}
		case "chinaz":
			a.sources[source] = &chinaz.Source{}
		case "securitytrails":
			a.sources[source] = &securitytrails.Source{}
		}
	}
}

// removeSources deletes the given sources from the source map
func (a *Agent) removeSources(sources []string) {
	for _, source := range sources {
		delete(a.sources, source)
	}
}
