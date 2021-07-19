package yougetsignal

import (
	"bytes"
	"context"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"iplookup/subscraping"
	"strings"
)

type yougetsignalResponse struct {
	Status          string     `json:"status"`
	Resultsmethod   string     `json:"resultsMethod"`
	Lastscrape      string     `json:"lastScrape"`
	Domaincount     string     `json:"domainCount"`
	Remoteaddress   string     `json:"remoteAddress"`
	Remoteipaddress string     `json:"remoteIpAddress"`
	Domainarray     [][]string `json:"domainArray"`
}

type Source struct{}

// Run function returns all subdomains found with the service
func (s *Source) Run(ctx context.Context, ip string, session *subscraping.Session) <-chan subscraping.Result {
	results := make(chan subscraping.Result)

	go func() {
		defer close(results)

		body := "remoteAddress=" + ip
		resp, err := session.SimplePost(ctx, "https://domains.yougetsignal.com/domains.php", "application/x-www-form-urlencoded", bytes.NewBufferString(body))
		if err != nil {
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: err}
			session.DiscardHTTPResponse(resp)
			return
		}
		body1, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		//println(resp.StatusCode)
		if resp.StatusCode != 403 {

			if strings.Contains(string(body1), "No web sites found.") {
				resp.Body.Close()
				return
			} else if strings.Contains(string(body1), "Service unavailable.") {
				resp.Body.Close()
				return
			}
			//fmt.Printf("%s", body1)

			var response yougetsignalResponse

			err = jsoniter.NewDecoder(bytes.NewBufferString(string(body1))).Decode(&response)
			//fmt.Printf("%v",response)
			if err != nil {
				results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: err}
				resp.Body.Close()
				return
			}

			for _, subdomains := range response.Domainarray {
				//fmt.Printf("%v", subdomain)
				results <- subscraping.Result{Source: s.Name(), Type: subscraping.Subdomain, Value: subdomains[0]}
				//if i == threshold {
				//	break
				//}
			}
		}

	}()

	return results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "yougetsignal"
}
