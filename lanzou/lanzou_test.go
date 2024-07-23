package lanzou

import (
	"encoding/json"
	"fmt"
	"testing"
)

// 带密码文件夹：https://ww0.lanzouj.com/b03dwumkli，密码 2024
// 不带密码文件夹：https://oyp.lanzoue.com/b083tujwh
// 带密码文件：https://oyp.lanzoue.com/ilF46iudy0f，密码 1234
// 不带密码文件：https://oyp.lanzoue.com/iSQzC0kfd5wb

// 测试：获取文件夹（带密码）内最新一个文件的直链
func TestGetLatestFile1(t *testing.T) {
	info, err := GetLatestFile("https://ww0.lanzouj.com/b03dwumkli", "2024")
	if err != nil {
		t.Error(err)
	} else {
		text, _ := json.Marshal(info)
		fmt.Println(string(text))
	}
}

// 测试：获取文件夹（不带密码）内最新一个文件的直链
func TestGetLatestFile2(t *testing.T) {
	info, err := GetLatestFile("https://oyp.lanzoue.com/b083tujwh", "")
	if err != nil {
		t.Error(err)
	} else {
		text, _ := json.Marshal(info)
		fmt.Println(string(text))
	}
}

// 测试：获取单个文件（带密码）的直链
func TestGetDownloadInfo1(t *testing.T) {
	info, err := GetDownloadInfo("https://oyp.lanzoue.com/ilF46iudy0f", "1234", false)
	if err != nil {
		t.Error(err)
	} else {
		text, _ := json.Marshal(info)
		fmt.Println(string(text))
	}
}

// 测试：获取单个文件（不带密码）的直链
func TestGetDownloadInfo2(t *testing.T) {
	info, err := GetDownloadInfo("https://oyp.lanzoue.com/iSQzC0kfd5wb", "", false)
	if err != nil {
		t.Error(err)
	} else {
		text, _ := json.Marshal(info)
		fmt.Println(string(text))
	}
}

// 测试：获取文件夹（带密码）内任意页码的文件列表
func TestGetFileList1(t *testing.T) {
	info, err := GetFileList("https://ww0.lanzouj.com/b03dwumkli", "2024", 0)
	if err != nil {
		t.Error(err)
	} else {
		text, _ := json.Marshal(info)
		fmt.Println(string(text))
	}
}

// 测试：获取文件夹（不带密码）内任意页码的文件列表
func TestGetFileList2(t *testing.T) {
	info, err := GetFileList("https://oyp.lanzoue.com/b083tujwh", "", 0)
	if err != nil {
		t.Error(err)
	} else {
		text, _ := json.Marshal(info)
		fmt.Println(string(text))
	}
}
