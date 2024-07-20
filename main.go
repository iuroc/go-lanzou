package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Printf("[ https://github.com/iuroc/go-lanzou ]\n\n")
	for {
		func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf(
						"❌ [main] 解析失败，原因是：%s\n\n%s\n\n",
						err.(error).Error(),
						strings.Repeat("-", 50),
					)
				}
			}()
			fmt.Print("👉 请输入蓝奏云文件分享链接：")
			var shareURL string
			fmt.Scan(&shareURL)
			downURL, err := GetDownloadURL(shareURL)
			if err != nil {
				panic(err)
			}
			fmt.Printf(
				"🍉 文件直链解析成功\n%s\n\n%s\n\n",
				downURL,
				strings.Repeat("-", 50),
			)
		}()
	}
}