package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/iuroc/go-lanzou/lanzou"
)

// 带密码文件夹：https://ww0.lanzouj.com/b03dwumkli，密码 2024
// 不带密码文件夹：https://oyp.lanzoue.com/b083tujwh
// 带密码文件：https://oyp.lanzoue.com/ilF46iudy0f，密码 1234
// 不带密码文件：https://oyp.lanzoue.com/iSQzC0kfd5wb

func main() {
	Test1()
	fmt.Println(strings.Repeat("🍑", 20), "Test1 测试通过", strings.Repeat("🍑", 20))
	Test2()
	fmt.Println(strings.Repeat("🍑", 20), "Test2 测试通过", strings.Repeat("🍑", 20))
	Test3()
	fmt.Println(strings.Repeat("🍑", 20), "Test3 测试通过", strings.Repeat("🍑", 20))
	Test4()
	fmt.Println(strings.Repeat("🍑", 20), "Test4 测试通过", strings.Repeat("🍑", 20))
	Test5()
	fmt.Println(strings.Repeat("🍑", 20), "Test5 测试通过", strings.Repeat("🍑", 20))
	Test6()
	fmt.Println(strings.Repeat("🍑", 20), "Test6 测试通过", strings.Repeat("🍑", 20))
	fmt.Println("测试通过")
}

// 测试：获取文件夹（带密码）内最新一个文件的直链
func Test1() {
	info, err := lanzou.GetLatestFile("https://ww0.lanzouj.com/b03dwumkli", "2024")
	if err != nil {
		log.Fatalln(err)
	} else {
		text, _ := json.Marshal(info)
		fmt.Println(string(text))
	}
}

// 测试：获取文件夹（不带密码）内最新一个文件的直链
func Test2() {
	info, err := lanzou.GetLatestFile("https://oyp.lanzoue.com/b083tujwh", "")
	if err != nil {
		log.Fatalln(err)
	} else {
		text, _ := json.Marshal(info)
		fmt.Println(string(text))
	}
}

// 测试：获取单个文件（带密码）的直链
func Test3() {
	info, err := lanzou.GetDownloadInfo("https://oyp.lanzoue.com/ilF46iudy0f", "1234")
	if err != nil {
		log.Fatalln(err)
	} else {
		text, _ := json.Marshal(info)
		fmt.Println(string(text))
	}
}

// 测试：获取单个文件（不带密码）的直链
func Test4() {
	info, err := lanzou.GetDownloadInfo("https://oyp.lanzoue.com/iSQzC0kfd5wb", "")
	if err != nil {
		log.Fatalln(err)
	} else {
		text, _ := json.Marshal(info)
		fmt.Println(string(text))
	}
}

// 测试：获取文件夹（带密码）内任意页码的文件列表
func Test5() {
	info, err := lanzou.GetFileList("https://ww0.lanzouj.com/b03dwumkli", "2024", 0)
	if err != nil {
		log.Fatalln(err)
	} else {
		text, _ := json.Marshal(info)
		fmt.Println(string(text))
	}
}

// 测试：获取文件夹（不带密码）内任意页码的文件列表
func Test6() {
	info, err := lanzou.GetFileList("https://oyp.lanzoue.com/b083tujwh", "", 0)
	if err != nil {
		log.Fatalln(err)
	} else {
		text, _ := json.Marshal(info)
		fmt.Println(string(text))
	}
}
