package webscan

import (
	"context"
	"fmt"
	"github.com/Lengso/iplookup/pkg/subscraping"
	jsoniter "github.com/json-iterator/go"
)

type webscanResponse []struct {
	Domain string `json:"domain"`
	Title  string `json:"title"`
}

type Source struct{}

// Run function returns all subdomains found with the service
func (s *Source) Run(ctx context.Context, ip string, session *subscraping.Session) <-chan subscraping.Result {
	results := make(chan subscraping.Result)

	go func() {
		defer close(results)

		resp, err := session.SimpleGet(ctx, fmt.Sprintf("http://api.webscan.cc/?action=query&ip=%s", ip))
		if err != nil {
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: err}
			session.DiscardHTTPResponse(resp)
			return
		}

		var response webscanResponse
		err = jsoniter.NewDecoder(resp.Body).Decode(&response)
		//fmt.Printf("%v",response)
		if err != nil {
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: err}
			resp.Body.Close()
			return
		}

		defer resp.Body.Close()

		for _, subdomain := range response {
			//i++
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Subdomain, Value: subdomain.Domain}
			//if i ==  threshold {
			//	break
			//}
		}

	}()

	return results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "webscan"
}
