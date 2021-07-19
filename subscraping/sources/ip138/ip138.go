// Package hackertarget logic
package ip138

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

		resp, err := session.SimpleGet(ctx, fmt.Sprintf("https://site.ip138.com/%s/", ip))
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
		var Ip138_regexp = regexp.MustCompile(`</span><a href="/(.*?)/" target="_blank">`)
		match := Ip138_regexp.FindAllStringSubmatch(string(body), -1)
		//fmt.Println(string(match))
		out := make([]string, len(match))
		for i := range out { // match result  return Subdomain struct
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Subdomain, Value: match[i][1]}
			//if threshold == i  {
			//	break
			//}
		}
		//
		//for _, subdomain := range subdomains {
		//	for _, value := range session.Extractor.FindAllString(subdomain, -1) {
		//		results <- subscraping.Result{Source: s.Name(), Type: subscraping.Subdomain, Value: value}
		//	}
		//}
	}()
	//fmt.Printf("test %+v",results)
	return results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "ip138"
}
