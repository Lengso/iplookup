package chinaz

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/Lengso/iplookup/pkg/subscraping"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"strings"
)

type chinaResponse struct {
	Statecode int    `json:"StateCode"`
	Message   string `json:"Message"`
	Result    []struct {
		Index         int           `json:"index"`
		Host          string        `json:"host"`
		Title         string        `json:"title"`
		Rank          int           `json:"rank"`
		Shoulu        string        `json:"shoulu"`
		Ipaddresslist []interface{} `json:"IPAddressList"`
	} `json:"Result"`
	Total      int `json:"Total"`
	Totalpages int `json:"TotalPages"`
}

type Source struct{}

type agent struct {
	Results   chan subscraping.Result
	Session   *subscraping.Session
	PageCount int
}

func (a *agent) enumerate(ctx context.Context, ip string, page int) {
	select {
	case <-ctx.Done():
		return
	default:
	}
	var body string
	b := []byte(ip)

	base64ip := base64.StdEncoding.EncodeToString(b)

	if page == 0 {
		body = fmt.Sprintf("page=1&iplist=%s", base64ip)
	} else {
		body = fmt.Sprintf("page=%d&iplist=%s", page, base64ip)
	}

	resp, err := a.Session.SimplePost(ctx, "https://s.tool.chinaz.com/AjaxTool.ashx?action=getsameip&callback=", "application/x-www-form-urlencoded", bytes.NewBufferString(body))
	if err != nil {
		a.Results <- subscraping.Result{Source: "chinaz", Type: subscraping.Error, Error: err}
		a.Session.DiscardHTTPResponse(resp)
		return
	}
	body1, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 404 {

		var response chinaResponse

		rebody := strings.Replace(string(body1), "(", "", -1)
		rebody1 := strings.Replace(rebody, ")", "", -1)

		err = jsoniter.NewDecoder(bytes.NewBufferString(rebody1)).Decode(&response)

		if err != nil {
			a.Results <- subscraping.Result{Source: "chinaz", Type: subscraping.Error, Error: err}
			resp.Body.Close()
			return
		}

		if response.Statecode < 0 {
			resp.Body.Close()
			return
		}
		defer resp.Body.Close()

		for _, subdomains := range response.Result {
			a.Results <- subscraping.Result{Source: "chinaz", Type: subscraping.Subdomain, Value: subdomains.Host}
		}

		if response.Total > 20 {
			a.PageCount = response.Total / 20
		}
	}
}

// Run function returns all subdomains found with the service
func (s *Source) Run(ctx context.Context, ip string, session *subscraping.Session) <-chan subscraping.Result {
	results := make(chan subscraping.Result)

	a := agent{
		Session: session,
		Results: results,
	}

	go func() {

		a.enumerate(ctx, ip, 0)

		if a.PageCount >= 1 {
			for i := 0; i < a.PageCount; i++ {
				a.enumerate(ctx, ip, a.PageCount+1)
			}
		}

		close(a.Results)
	}()

	return a.Results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "chianz"
}
