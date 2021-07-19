package omnisint

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"iplookup/subscraping"
)

type omnisintResponse []string

type Source struct{}

// Run function returns all subdomains found with the service
func (s *Source) Run(ctx context.Context, ip string, session *subscraping.Session) <-chan subscraping.Result {
	results := make(chan subscraping.Result)

	go func() {
		defer close(results)

		resp, err := session.SimpleGet(ctx, fmt.Sprintf("https://sonar.omnisint.io/reverse/%s", ip))
		if err != nil {
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: err}
			session.DiscardHTTPResponse(resp)
			return
		}

		var response omnisintResponse
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
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Subdomain, Value: subdomain}
			//if threshold == i  {
			//	break
			//}
		}
	}()

	return results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "omnisint"
}
