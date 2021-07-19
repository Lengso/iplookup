package bugscaner

import (
	"context"
	"fmt"
	"io/ioutil"
	"iplookup/subscraping"
	"regexp"
)

// Source is the passive scraping agent
type Source struct{}

// Run function returns all subdomains found with the service
func (s *Source) Run(ctx context.Context, ip string, session *subscraping.Session) <-chan subscraping.Result {
	results := make(chan subscraping.Result)

	go func() {
		defer close(results)

		resp, err := session.SimpleGet(ctx, fmt.Sprintf("http://dns.bugscaner.com/%s.html", ip))
		if err != nil {
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: err}
			session.DiscardHTTPResponse(resp)
			return
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		//fmt.Println(string(body))
		var Viewdns_regexp = regexp.MustCompile(`<a href="http://dns.bugscaner.com/(.+?).html"`)

		match := Viewdns_regexp.FindAllStringSubmatch(string(body), -1)
		//fmt.Println(string(match))
		out := make([]string, len(match))
		//fmt.Printf(" %v",body)
		for i := range out { // match result  return Subdomain struct
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Subdomain, Value: match[i][1]}
			//if threshold == i  {
			//	break
			//}
		}

	}()

	return results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "bugscaner"
}
