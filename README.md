# iplookup

## 简介
IP反查域名工具,模仿(抄袭)subfinder实现

## 主要接口

- [x] [webscan](https://www.webscan.cc/)  
- [x] [rapiddns](https://rapiddns.io)  
- [x] [ip138](https://site.ip138.com/)  
- [x] [yougetsignal](https://www.yougetsignal.com/)  
- [x] [aizhan](https://dns.aizhan.com/)  
- [x] [c99](https://api.c99.nl/)
- [x] [chinaz](http://s.tool.chinaz.com)  
- [x] [viewdns](https://viewdns.info/)  
- [x] [bugscaner](http://dns.bugscaner.com/)  
- [x] [hackertarget](https://api.hackertarget.com/)
- [x] [dnslytics](https://dnslytics.com/reverse-ip)  
- [x] [omnisint](https://omnisint.io/)
- [x] [dnsgrep](https://www.dnsgrep.cn/)  
- [x] [domaintools](https://reverseip.domaintools.com/)  
- [x] [securitytrails](https://securitytrails.com/)
- [x] [fofa](https://fofa.so/)
- [x] [shodan](https://api.shodan.io/dns/reverse)
- [ ] [quake](https://quake.360.cn/)
- [ ] [ipip](https://tools.ipip.net/ipdomain.php)


## 使用说明
 
有些接口会有历史绑定域名,默认提取50个域名防止查询CDN域名数量过多

可自行设置阈值`-count 9999`

查询过多会导致IP被封，建议搭配代理使用

可以搭配httpx,nuclei等工具食用效果更佳

以下接口需要设置 API 密钥。

- [C99](https://api.c99.nl/)
- [Fofa](https://fofa.so/)
- [Quake](https://quake.360.cn/)
- [Shodan](https://www.shodan.io/)

  
## usege  

常规用法
```sh
echo 1.1.1.1 | ipgo 
cat ips.txt | ipgo -oD out
#搭配httpx使用
ipgo.exe -i 1.1.1.1 -silent | httpx -title -ip -content-length -status-code -tech-detect -random-agent
#设置阈值 
ipgo.exe -count 9999 -iL ips.txt -oD out  
```

编译命令
```sh
make build-all
```

完整参数
```sh
C:\Users\administrator
λ ipgo

  _       _             _
 (_)_ __ | | ___   ___ | | ___   _ _ __
 | | '_ \| |/ _ \ / _ \| |/ / | | | '_ \
 | | |_) | | (_) | (_) |   <| |_| | |_) |
 |_| .__/|_|\___/ \___/|_|\_\\__,_| .__/
   |_|                            |_|      v1.1

[FTL] Program exiting: no input list provided

C:\Users\ot
λ ipgo -h
Usage of ipgo:
  -all
        Use all sources (slow) for enumeration
  -config string
        Configuration file for API Keys, etc (default "C:\\Users\\administrator/.config/iplookup/config.yaml")
  -count int
          Number of domain name thresholds (default 50)
  -exclude-sources string
        List of sources to exclude from enumeration
  -i string
        ip to find domain for
  -iL string
        File containing list of ips to enumerate
  -json
        Write output in JSON lines Format
  -max-time int
        Minutes to wait for enumeration results (default 10)
  -nC
        Don't Use colors in output
  -o string
        File to write output to (optional)
  -oD string
        Directory to write enumeration results to (optional)
  -silent
        Show only subdomains in output
  -sources string
        Comma separated list of sources to use
  -t int
        Number of concurrent goroutines for resolving (default 10)
  -timeout int
        Seconds to wait before timing out (default 30)
  -v    Show Verbose output
  -version
```

一个示例配置文件 `$HOME/.config/iplookup/config.yaml`
```yaml
sources:
  - webscan
  - rapiddns
  - ip138
  - yougetsignal
  - aizhan
  - chinaz
  - viewdns
  - c99
  - bugscaner
  - hackertarget
  - dnslytics
  - omnisint
  - dnsgrep
  - domaintools
  - securitytrails
  - fofa
  - shodan
all-sources:
  - webscan
  - rapiddns
  - ip138
  - yougetsignal
  - aizhan
  - c99
  - chinaz
  - viewdns
  - bugscaner
  - hackertarget
  - dnslytics
  - omnisint
  - dnsgrep
  - domaintools
  - securitytrails
  - fofa
  - shodan
proxy: "http://127.0.0.1:8080/"
dnsgrep: []
c99:
  - XXXXX-XXXXX-XXXXX-XXXXX
shodan:
  - XXXXX-XXXXX-XXXXX-XXXXX
fofa:
  - XXXXX@gmail.com:xxx
iplookup-version: "1.1"
```

包引用,配置文件修改为```config/iplookup.yaml```
```golang
package main

import (
	"fmt"
	"github.com/Lengso/iplookup"
)

func main() {
	output := iplookup.GetDomain("1.1.1.1")

	for _,domain := range output{
		fmt.Println(domain)
	}

}

```

## 参考
https://github.com/projectdiscovery/subfinder

