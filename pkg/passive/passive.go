package passive

import (
	"context"
	"fmt"
	"github.com/Lengso/iplookup/pkg/subscraping"
	"sync"
	"time"

	"github.com/projectdiscovery/gologger"
)

// EnumerateSubdomains enumerates all the subdomains for a given domain
func (a *Agent) EnumerateIp(ip string, keys *subscraping.Keys, proxy *subscraping.Proxy, timeout int, maxEnumTime time.Duration) chan subscraping.Result {
	results := make(chan subscraping.Result)

	go func() {
		session := subscraping.NewSession(keys, proxy, timeout) // switch keys struct list
		//if err != nil {
		//	results <- subscraping.Result{Type: subscraping.Error, Error: fmt.Errorf("could not init passive session for %s: %s", ip, err)}
		//}
		//fmt.Printf("session %+v", session)

		ctx, cancel := context.WithTimeout(context.Background(), maxEnumTime)

		timeTaken := make(map[string]string)
		timeTakenMutex := &sync.Mutex{}

		wg := &sync.WaitGroup{}
		// Run each source in parallel on the target domain
		for source, runner := range a.sources {
			wg.Add(1) // 遍历存在的接口 每个接口添加1个wg workder

			now := time.Now()
			go func(source string, runner subscraping.Source) { //匿名函数，执行source中的接口
				//fmt.Println(runner.Name())
				for resp := range runner.Run(ctx, ip, session) {
					results <- resp
				}

				//fmt.Printf("%v",results)
				duration := time.Since(now)
				timeTakenMutex.Lock()
				timeTaken[source] = fmt.Sprintf("Source took %s for enumeration\n", duration) // api time
				timeTakenMutex.Unlock()

				wg.Done()
			}(source, runner)

		}
		wg.Wait()

		for source, data := range timeTaken {
			gologger.Verbose().Label(source).Msg(data)
		}

		close(results)
		cancel()
	}()

	return results
}
