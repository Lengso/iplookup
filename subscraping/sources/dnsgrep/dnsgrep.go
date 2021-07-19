// Package dnsgrep logic
package dnsgrep

import (
	"context"
	"fmt"
	"io/ioutil"
	"iplookup/subscraping"
	"regexp"
)

type dsngrepResponse struct {
	Status int `json:"status"`
	Data   struct {
		Data []struct {
			Domain string `json:"domain"`
			Value  string `json:"value"`
			Type   string `json:"type"`
			Time   string `json:"time"`
		} `json:"data"`
		Count int `json:"count"`
		Type  int `json:"type"`
	} `json:"data"`
}

// Source is the passive scraping agent
type Source struct{}

// Run function returns all subdomains found with the service
func (s *Source) Run(ctx context.Context, ip string, session *subscraping.Session) <-chan subscraping.Result {
	results := make(chan subscraping.Result)

	go func() {
		defer close(results)

		resp, err := session.SimpleGet(ctx, fmt.Sprintf("https://www.dnsgrep.cn/ip/%s", ip))
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

		var Dnsgrep_regexp = regexp.MustCompile(`<tr>\s+<td data="(.*?)">`)
		match := Dnsgrep_regexp.FindAllStringSubmatch(string(body), -1)

		out := make([]string, len(match))

		for i := range out { // match result  return Subdomain struct
			//fmt.Printf("%s", match[i][1])
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Subdomain, Value: match[i][1]}
		}
	}()

	return results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "dnsgrep"
}
