# Surge Mac 连接百度 CDN 网络代理的小玩具

刘同志在几年前，为 Surge Mac 引入了名为 “[External Proxy Provider](https://community.nssurge.com/d/3-external-proxy-provider)” 的新特性，处于 “免流” 或加速中国骨干网访问，特此用 Go 开发了一个小玩具。

- 你也可以参考：“[奶昔论坛文档](https://forum.naixi.net/thread-9195-1-1.html)”


## 配置在 Surge Mac

新增如下的配置即可：

```yaml
[Proxy]
BaiduCDN = external, exec = "/path/to/your/BaiduCDNProxySurgeMac", args = "-p", args = "18964", local-port = 18964
```
