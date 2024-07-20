# go-lanzou

Go 语言实现的蓝奏云直链解析程序。

## 快速开始

```shell
# 运行项目
go run .

# 生成可执行文件
go build
```

🍎 你也可以[下载可执行文件]()进行使用。

## 作为模块

```shell
go get https://github.com/iuroc/go-lanzou
```

```go
package main

import (
    goLanzou "https://github.com/iuroc/go-lanzou"
)

func main() {
    shareURL := "https://www.lanzoui.com/imcSy2340ssb"
    downloadURL, err := goLanzou.GetDownloadURL(shareURL)
    if err != nil {
        fmt.Println("解析失败")
    } else {
        fmt.Println("解析成功：" + downloadURL)
    }
}
```
