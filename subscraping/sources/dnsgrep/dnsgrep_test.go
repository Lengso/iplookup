package dnsgrep_test

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"iplookup/subscraping"
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

		if session.Keys.Dnsgrep == "" {
			return
		}

		resp, err := session.SimpleGet(ctx, fmt.Sprintf("https://www.dnsgrep.cn/api/query?q=%s&token=%s", ip, session.Keys.Dnsgrep))
		if err != nil {
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: err}
			session.DiscardHTTPResponse(resp)
			return
		}

		var response dsngrepResponse
		err = jsoniter.NewDecoder(resp.Body).Decode(&response)
		//fmt.Printf("%v",response)
		if err != nil {
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: err}
			resp.Body.Close()
			return
		}

		defer resp.Body.Close()

		//bufforesponse
		//fmt.Printf("%v",response.Data.Data)
		//
		if response.Data.Count > 0 {
			for _, subdomain := range response.Data.Data {
				results <- subscraping.Result{Source: s.Name(), Type: subscraping.Subdomain, Value: subdomain.Domain}
			}
		}

	}()

	return results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "dnsgrep_api"
}
