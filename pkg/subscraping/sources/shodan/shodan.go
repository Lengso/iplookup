package shodan

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/Lengso/iplookup/pkg/subscraping"
)

// Source is the passive scraping agent
type Source struct{}

var ShodanReponse map[string]interface{}

// Run function returns all subdomains found with the service
func (s *Source) Run(ctx context.Context, ip string, session *subscraping.Session) <-chan subscraping.Result {
	results := make(chan subscraping.Result)

	go func() {
		defer close(results)

		if session.Keys.Shodan == "" {
			return
		}

		searchURL := fmt.Sprintf("https://api.shodan.io/dns/reverse?ips=%s&key=%s", ip, session.Keys.Shodan)
		resp, err := session.SimpleGet(ctx, searchURL)
		if err != nil {
			session.DiscardHTTPResponse(resp)
			return
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		json.Unmarshal(body, &ShodanReponse)
		if domain, ok := ShodanReponse[ip]; ok {
			if domain != nil {
				value := strings.ReplaceAll(fmt.Sprintf("%s", domain), "[", "")
				value = strings.ReplaceAll(fmt.Sprintf("%s", value), "]", "")
				results <- subscraping.Result{Source: s.Name(), Type: subscraping.Subdomain, Value: value}
			}
		}
	}()

	return results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "shodan"
}
