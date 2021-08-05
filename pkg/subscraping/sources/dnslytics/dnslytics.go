package dnslytics

import (
	"bytes"
	"context"
	"github.com/Lengso/iplookup/pkg/subscraping"
	"io/ioutil"
	"regexp"
	"strings"
)

type Source struct{}

// Run function returns all subdomains found with the service
func (s *Source) Run(ctx context.Context, ip string, session *subscraping.Session) <-chan subscraping.Result {
	results := make(chan subscraping.Result)

	go func() {
		defer close(results)

		body := "reverseip=" + ip
		resp, err := session.SimplePost(ctx, "https://dnslytics.com/reverse-ip", "application/x-www-form-urlencoded", bytes.NewBufferString(body))
		if err != nil {
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: err}
			session.DiscardHTTPResponse(resp)
			return
		}
		body1, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}

		if strings.Contains(string(body1), "No domains found hosted") {
			resp.Body.Close()
			return
		}
		//fmt.Printf("%s", body1)

		var Dnslytics_regexp = regexp.MustCompile(`</td><td><b>(.*?)</b></td><td><form`)
		match := Dnslytics_regexp.FindAllStringSubmatch(string(body1), -1)
		//fmt.Println(string(match))
		//fmt.Printf("%v", match)
		out := make([]string, len(match))
		for i := range out { // match result  return Subdomain struct
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Subdomain, Value: match[i][1]}
		}

	}()

	return results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "dnslytics"
}
