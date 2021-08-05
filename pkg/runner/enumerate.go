package runner

import (
	"context"
	"github.com/Lengso/iplookup/pkg/subscraping"
	"io"
	"sync"
	"time"

	"github.com/hako/durafmt"
	"github.com/projectdiscovery/gologger"
)

const maxNumCount = 2

// EnumerateSingleDomain performs subdomain enumeration against a single domain
func (r *Runner) EnumerateSingleDomain(ctx context.Context, ip string, outputs []io.Writer) error { //处理器函数处理各个接口返回结果， 统一处理结果
	gologger.Info().Msgf("Enumerating IP for %s\n", ip)

	// Get the API keys for sources from the configuration
	// and also create the active resolving engine for the domain.
	keys := r.options.YAMLConfig.GetKeys() //获取key

	proxy := r.options.YAMLConfig.GetProxy() //获取proxy

	// Run the passive subdomain enumeration
	now := time.Now()
	passiveResults := r.passiveAgent.EnumerateIp(ip, &keys, &proxy, r.options.Timeout, time.Duration(r.options.MaxEnumerationTime)*time.Minute) //

	wg := &sync.WaitGroup{}
	wg.Add(1)
	// Create a unique map for filtering duplicate subdomains out
	uniqueMap := make(map[string]subscraping.HostEntry)
	// Create a map to track sources for each host
	sourceMap := make(map[string]map[string]struct{})
	// Process the results in a separate goroutine
	go func() {
		for result := range passiveResults { // 遍历查询结果
			switch result.Type {
			case subscraping.Error: // error for api
				gologger.Warning().Msgf("Could not run source %s: %s\n", result.Source, result.Error)
			case subscraping.Subdomain:
				// Validate the subdomain found and remove wildcards from
				//if !strings.HasSuffix(result.Value, "."+ip) {
				//	continue
				//}

				subdomain := result.Value
				//subdomain := strings.ReplaceAll(strings.ToLower(result.Value), "*.", "")

				if _, ok := uniqueMap[subdomain]; !ok {
					sourceMap[subdomain] = make(map[string]struct{})
				}
				//fmt.Printf("%v",result)

				// Log the verbose message about the found subdomain per source
				if _, ok := sourceMap[subdomain][result.Source]; !ok {
					gologger.Verbose().Label(result.Source).Msg(subdomain)
				}

				sourceMap[subdomain][result.Source] = struct{}{}

				// Check if the subdomain is a duplicate. If not,
				// send the subdomain for resolution.
				if _, ok := uniqueMap[subdomain]; ok { //去除重复
					continue
				}

				hostEntry := subscraping.HostEntry{Host: subdomain, Source: result.Source}
				//fmt.Println(hostEntry.Host,hostEntry.Source)

				if len(uniqueMap) >= r.options.Threshold {
					break
				}
				uniqueMap[subdomain] = hostEntry

				//gologger.Verbose().Label(result.Source).Msg(result.Count)
			}
		}
		wg.Done()
	}()

	// If the user asked to remove wildcards, listen from the results
	// queue and write to the map. At the end, print the found results to the screen
	wg.Wait()

	//	uniqueMap := make(map[string]subscraping.HostEntry) copy

	outputter := NewOutputter(r.options.JSON)

	// Now output all results in output writers
	var err error
	for _, w := range outputs {
		err = outputter.WriteHost(uniqueMap, w)
		if err != nil {
			gologger.Error().Msgf("Could not verbose results for %s: %s\n", ip, err)
			return err
		}
	}

	// Show found subdomain count in any case.
	duration := durafmt.Parse(time.Since(now)).LimitFirstN(maxNumCount).String()

	gologger.Info().Msgf("Found %d subdomains for %s in %s\n", len(uniqueMap), ip, duration) // Todo   run End

	return nil
}
