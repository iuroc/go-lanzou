# go-lanzou

Go 语言实现的蓝奏云直链解析程序。

## 已实现功能

- 获取单个文件（可带密码）的直链
- 获取文件夹（可带密码）内最新一个文件的直链
- 获取文件夹（可带密码）内任意页码的文件列表

## 快速开始

```shell
# 运行项目
go run .

# 生成可执行文件
go build
```

🍎 你也可以[下载可执行文件](https://github.com/iuroc/go-lanzou/releases/download/1.1.2/go-lanzou.exe)进行使用。

## 作为模块

```shell
go get github.com/iuroc/go-lanzou
```

```go
package main

import (
	"fmt"
	"github.com/iuroc/go-lanzou/lanzou"
)

func main() {
	shareURL := "https://www.lanzoui.com/imcSy2340ssb"
	downloadURL, err := lanzou.GetDownloadURL(shareURL)
	if err != nil {
		fmt.Println("解析失败")
	} else {
		fmt.Println("解析成功：" + downloadURL)
	}
}
```

## API

```go
// 获取文件夹最新的一个文件的信息，包含直链
//
// urlOrId 是文件夹的分享链接或 ID
//
// password 是访问密码
func lanzou.GetLatestFile(shareURL string, password string) (*lanzou.DownloadInfo, error)
```

```go
// 获取单个文件的信息，包含直链
func lanzou.GetDownloadInfo(shareURL string, password string) (*lanzou.DownloadInfo, error)
```

```go
// 获取文件夹中指定页码的文件列表
//
// page 的值务必从 0 开始，每次只允许增长 1，不可以直接从 0 变为 2。
//
// 每次换页，务必暂停 1 秒以上。
func lanzou.GetFileList(shareURL string, password string, page int) ([]FileInfo, error)
```

## 蓝奏云分享合集

[蓝奏云分享合集](https://github.com/iuroc/lanzou-collect/blob/master/V1/%E6%95%B0%E6%8D%AE%E6%BA%90/%E6%A0%A1%E9%AA%8C%E6%88%90%E5%8A%9F%E6%95%B0%E6%8D%AE%E6%BA%90.txt)
