# Surge Mac 连接百度 CDN 网络代理的小玩具

刘同志在几年前，为 Surge Mac 引入了名为 “[External Proxy Provider](https://community.nssurge.com/d/3-external-proxy-provider)” 的新特性，处于 “免流” 或加速中国骨干网访问，特此用 Go 开发了一个小玩具。

- 你也可以参考：“[奶昔论坛文档](https://forum.naixi.net/thread-9195-1-1.html)”


## 配置在 Surge Mac

新增如下的配置即可：

``` yaml
[Proxy]
BaiduCDN = external, exec = "/path/to/your/BaiduCDNProxySurgeMac", args = "-p", args = "18964", local-port = 18964
```


## 完整教程

请前往本仓库的 “[Releases](https://github.com/imPrk0/BaiduCDNProxySurgeMac/releases)”，找到标记为 “Latest” (最新版) 的版本，在 “Assets” (资产) 的文件列表中，会有两种类型的文件：

> 因为我写这篇文档的时候的最新版本是 `v1.0.1`，但是前面都一样，就后面的版本号不一样

- `BaiduCDNProxySurgeMac_drawin_arm64_v1.0.1.zip`
- `BaiduCDNProxySurgeMac_drawin_amd64_v1.0.1.zip`

出了固定的文件名：`BaiduCDNProxySurgeMac`，和作业系统 “`drawin`” (毕竟 Surge Mac 只出 macOS 版) 以及版本号外，最大的区别就是 “`arm64`” 和 “`amd64`” 了。

如果你的 Mac 是 Apple silicon 的 (搭载 M* / A* 晶片的设备)，那么你要选择的是 “`arm64`”，如果你的 Mac 是搭载的 Intel 芯片的，那么你要选择的是 “`amd64`”。

下载好后解压缩，放入目录，比如 `/usr/local/bin`，然后在配置你的 Surge：

``` yaml
[Proxy]
BaiduCDN = external, exec = "/usr/local/bin/BaiduCDNProxySurgeMac", args = "-p", args = "18964", local-port = 18964
```

完成。


## 自行编译

Go 语言的标准程序，你也可以参考我的命令：

``` sh
GOOS=darwin GOARCH=arm64 go build -trimpath -o release/BaiduCDNProxySurgeMac_arm64 .
GOOS=darwin GOARCH=amd64 go build -trimpath -o release/BaiduCDNProxySurgeMac_amd64 .
```
