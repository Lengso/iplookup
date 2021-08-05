// Package hackertarget logic
package hackertarget

import (
	"bufio"
	"context"
	"fmt"
	"github.com/Lengso/iplookup/pkg/subscraping"
	"regexp"
)

// Source is the passive scraping agent
type Source struct{}

var Hackertarget_regexp = regexp.MustCompile(`</span><a href="/(.*?)/" target="_blank">`)

// Run function returns all subdomains found with the service
func (s *Source) Run(ctx context.Context, ip string, session *subscraping.Session) <-chan subscraping.Result {
	results := make(chan subscraping.Result)

	go func() {
		defer close(results)

		resp, err := session.SimpleGet(ctx, fmt.Sprintf("http://api.hackertarget.com/reverseiplookup/?q=%s", ip))
		if err != nil {
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: err}
			session.DiscardHTTPResponse(resp)
			return
		}

		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() { // 先把body转换成输入流，再从里面查找域名
			line := scanner.Text()
			if line == "" {
				continue
			}
			match := Hackertarget_regexp.FindAllString(line, -1)
			for _, subdomain := range match { // match result  return Subdomain struct
				results <- subscraping.Result{Source: s.Name(), Type: subscraping.Subdomain, Value: subdomain}
				//if threshold == i  {
				//	break
				//}
			}
		}
	}()

	return results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "hackertarget"
}
