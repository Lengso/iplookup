package passive

import (
	"github.com/Lengso/iplookup/pkg/subscraping"
	"github.com/Lengso/iplookup/pkg/subscraping/sources/aizhan"
	"github.com/Lengso/iplookup/pkg/subscraping/sources/bugscaner"
	"github.com/Lengso/iplookup/pkg/subscraping/sources/c99"
	"github.com/Lengso/iplookup/pkg/subscraping/sources/chinaz"
	"github.com/Lengso/iplookup/pkg/subscraping/sources/dnsgrep"
	"github.com/Lengso/iplookup/pkg/subscraping/sources/dnslytics"
	"github.com/Lengso/iplookup/pkg/subscraping/sources/domaintools"
	"github.com/Lengso/iplookup/pkg/subscraping/sources/fofa"
	"github.com/Lengso/iplookup/pkg/subscraping/sources/hackertarget"
	"github.com/Lengso/iplookup/pkg/subscraping/sources/ip138"
	"github.com/Lengso/iplookup/pkg/subscraping/sources/omnisint"
	"github.com/Lengso/iplookup/pkg/subscraping/sources/rapiddns"
	"github.com/Lengso/iplookup/pkg/subscraping/sources/securitytrails"
	"github.com/Lengso/iplookup/pkg/subscraping/sources/shodan"
	"github.com/Lengso/iplookup/pkg/subscraping/sources/viewdns"
	"github.com/Lengso/iplookup/pkg/subscraping/sources/webscan"
	"github.com/Lengso/iplookup/pkg/subscraping/sources/yougetsignal"
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
	"fofa",
	"shodan",
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
		case "fofa":
			a.sources[source] = &fofa.Source{}
		case "shodan":
			a.sources[source] = &shodan.Source{}
		}
	}
}

// removeSources deletes the given sources from the source map
func (a *Agent) removeSources(sources []string) {
	for _, source := range sources {
		delete(a.sources, source)
	}
}
