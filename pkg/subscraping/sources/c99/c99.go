package c99

import (
	"context"
	"fmt"
	"github.com/Lengso/iplookup/pkg/subscraping"
	jsoniter "github.com/json-iterator/go"
)

type c99Response struct {
	Success bool     `json:"success"`
	IP      string   `json:"ip"`
	Domains []string `json:"domains"`
	Error   string   `json:"error"`
}

type Source struct{}

// Run function returns all subdomains found with the service
func (s *Source) Run(ctx context.Context, ip string, session *subscraping.Session) <-chan subscraping.Result {
	results := make(chan subscraping.Result)

	go func() {
		defer close(results)

		if session.Keys.C99 == "" {
			return
		}

		resp, err := session.SimpleGet(ctx, fmt.Sprintf("https://api.c99.nl/ip2domains?key=%s&ip=%s&json", session.Keys.C99, ip))
		if err != nil {
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: err}
			session.DiscardHTTPResponse(resp)
			return
		}

		var response c99Response
		err = jsoniter.NewDecoder(resp.Body).Decode(&response)
		//fmt.Printf("%v",response)
		if err != nil {
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: err}
			resp.Body.Close()
			return
		}

		defer resp.Body.Close()

		if response.Success {
			for _, subdomain := range response.Domains {
				if subdomain == "    " {
					continue
				}
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
	return "c99"
}
