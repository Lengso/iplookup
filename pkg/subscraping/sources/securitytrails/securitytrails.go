package securitytrails

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Lengso/iplookup/pkg/subscraping"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"regexp"
	"strings"
)

type securityTrailsResponse struct {
	Props struct {
		Pageprops struct {
			IP               string `json:"ip"`
			Islist           bool   `json:"isList"`
			Locationsearch   string `json:"locationSearch"`
			Locationpathname string `json:"locationPathname"`
			Page             int    `json:"page"`
			Searchvalue      string `json:"searchValue"`
			Serverresponse   struct {
				Success    bool   `json:"success"`
				Status     int    `json:"status"`
				Statustext string `json:"statusText"`
				Error      string `json:"error"`
				Data       struct {
					Meta struct {
						LimitReached bool `json:"limit_reached"`
						MaxPage      int  `json:"max_page"`
						Page         int  `json:"page"`
						TotalPages   int  `json:"total_pages"`
					} `json:"meta"`
					Records []struct {
						HostProvider []string      `json:"host_provider"`
						Hostname     string        `json:"hostname"`
						MailProvider []interface{} `json:"mail_provider"`
						OpenPageRank interface{}   `json:"open_page_rank"`
					} `json:"records"`
					Total    int `json:"total"`
					Previews []struct {
						Hostname string `json:"hostname"`
						Rank     string `json:"rank"`
						Provider string `json:"provider"`
					} `json:"previews"`
				} `json:"data"`
				Asnrisklevel string `json:"asnRiskLevel"`
			} `json:"serverResponse"`
			Type string      `json:"type"`
			User interface{} `json:"user"`
		} `json:"pageProps"`
		NSsp bool `json:"__N_SSP"`
	} `json:"props"`
	Page  string `json:"page"`
	Query struct {
		IP string `json:"ip"`
	} `json:"query"`
	Buildid    string `json:"buildId"`
	Isfallback bool   `json:"isFallback"`
	Gssp       bool   `json:"gssp"`
}

// Source is the passive scraping agent
type Source struct{}

// Run function returns all subdomains found with the service
func (s *Source) Run(ctx context.Context, ip string, session *subscraping.Session) <-chan subscraping.Result {
	results := make(chan subscraping.Result)

	go func() {
		defer close(results)

		resp, err := session.SimpleGet(ctx, fmt.Sprintf("https://securitytrails.com/list/ip/%s/", ip))
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
		var securityTrails_Regexp = regexp.MustCompile(`/json">(.*?)</script>`)

		match := securityTrails_Regexp.FindString(string(body))
		rebody := strings.Replace(match, "/json\">", "", -1)
		rebody1 := strings.Replace(rebody, "</script>", "", -1)

		var response securityTrailsResponse
		err = jsoniter.NewDecoder(bytes.NewBufferString(rebody1)).Decode(&response)
		if err != nil {
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: err}
			resp.Body.Close()
			return
		}

		for _, subdomain := range response.Props.Pageprops.Serverresponse.Data.Records {
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Subdomain, Value: subdomain.Hostname}
		}
	}()

	return results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "securitytrails"
}
