# iplookup

## 简介
IP反查域名工具简称：IPgo，模仿(抄袭)subfinder实现

## 主要接口

- [x] webscan  
- [x] rapiddns  
- [x] ip138  
- [x] yougetsignal  
- [x] aizhan  
- [x] c99  
- [x] chinaz  
- [x] viewdns  
- [x] bugscaner  
- [x] hackertarget  
- [x] dnslytics  
- [x] omnisint  
- [x] dnsgrep  
- [x] domaintools  
- [x] securitytrails
- [ ] ipip
- [ ] fofa
- [ ] 360quake

## 使用说明
有些接口会有历史绑定域名,默认提取50个域名防止查询CDN域名数量过多

可自行设置阈值`-count 9999`

查询过多会导致IP被封，建议搭配代理使用

可以搭配httpx,nuclei等工具食用效果更佳

## usege  
```
echo 1.1.1.1 | ipgo 
ipgo.exe -i 1.1.1.1 -silent | httpx -title -ip -content-length -status-code -tech-detect -random-agent


cat ips.txt | ipgo -oD out
#设置阈值输出目录
ipgo.exe -count 9999 -iL ips.txt -oD out  
```
 
## 参考
https://github.com/projectdiscovery/subfinder

