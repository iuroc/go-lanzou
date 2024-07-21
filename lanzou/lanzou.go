// 蓝奏云直链解析程序
//
// 已实现功能：
//
// 1. 获取单个文件（可带密码）的直链
//
// 2. 获取文件夹（可带密码）内最新一个文件的直链
//
// 3. 获取文件夹（可带密码）内任意页码的文件列表
//
// 示例：
//
// package main
//
// import (
//
//	"fmt"
//	"github.com/iuroc/go-lanzou/lanzou"
//
// )
//
//	func main() {
//		shareURL := "https://www.lanzoui.com/imcSy2340ssb"
//		downloadURL, err := lanzou.GetDownloadURL(shareURL)
//		if err != nil {
//			fmt.Println("解析失败")
//		} else {
//			fmt.Println("解析成功：" + downloadURL)
//		}
//	}
package lanzou

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// 蓝奏云分享链接的开头部分，以域名结尾
const baseURL string = "https://iuroc.lanzoue.com"

// 获取单个文件的信息，包含直链
func GetDownloadInfo(shareURL string, password string) (*DownloadInfo, error) {
	fileId, err := getShareId(shareURL)
	if err != nil {
		return nil, err
	}
	html, err := SendRequest(RequestConfig{
		URL: baseURL + "/" + fileId,
		Headers: map[string]string{
			"User-Agent": "go-lanzou",
		},
	})
	if err != nil {
		return nil, err
	}

	fileInfo, _ := getFileInfoFromHTML(html)

	// 判断当前分享页是否需要密码
	if regexp.MustCompile(`class="passwdinput"`).MatchString(html) {
		if password == "" {
			fmt.Print("🔒 请输入文件访问密码：")
			fmt.Scan(&password)
		}
		downloadInfo, err := fetchWithPassword(html, password)
		downloadInfo.FileInfo.ShareId = fileId
		downloadInfo.FileInfo.Password = password
		downloadInfo.FileInfo.Name = fileInfo.Name
		downloadInfo.FileInfo.Size = fileInfo.Size
		downloadInfo.FileInfo.Date = fileInfo.Date
		return downloadInfo, err
	}

	iframeURLMatch := regexp.MustCompile(`src="(\/fn\?[^"]{20,})"`).FindStringSubmatch(html)
	if len(iframeURLMatch) == 0 {
		return nil, errors.New("[GetDownloadURL] 获取 iframeURL 失败")
	}
	iframeURL := baseURL + iframeURLMatch[1]
	downloadInfo, err := fetchIframe(iframeURL)
	if err != nil {
		return nil, err
	}

	downloadInfo.FileInfo.ShareId = fileId
	downloadInfo.FileInfo.Password = password
	downloadInfo.FileInfo.Name = fileInfo.Name
	downloadInfo.FileInfo.Size = fileInfo.Size
	downloadInfo.FileInfo.Date = fileInfo.Date
	return downloadInfo, err
}

// 获取页面中文件的一些信息，比如文件名、大小、时间
func getFileInfoFromHTML(html string) (FileInfo, error) {
	fileInfo := FileInfo{}
	nameMatch := regexp.MustCompile(`(?:padding: 56px 0px 20px 0px;">|id="filenajax">)(.*?)</div>`).FindStringSubmatch(html)
	if len(nameMatch) != 0 {
		fileInfo.Name = nameMatch[1]
	}
	sizeMatch := regexp.MustCompile(`(?:文件大小：</span>|<div class="n_filesize">大小：)(.*?)(?:<br>|</div>)`).FindStringSubmatch(html)
	if len(sizeMatch) != 0 {
		fileInfo.Size = sizeMatch[1]
	}
	dateMatch := regexp.MustCompile(`(?:上传时间：</span>|<div class="n_file_info"><span class="n_file_infos">)(.*?)(?:</span>|<br>)`).FindStringSubmatch(html)
	if len(dateMatch) != 0 {
		fileInfo.Date = dateMatch[1]
	}
	return fileInfo, nil
}

// 从分享链接中提取出标识字符串
func getShareId(urlOrId string) (string, error) {
	match := regexp.MustCompile(`^https?://.*?/([a-zA-Z0-9]+)`).FindStringSubmatch(urlOrId)
	if len(match) != 0 {
		return match[1], nil
	} else if regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(urlOrId) {
		return urlOrId, nil
	} else {
		return "", errors.New("[getShareId] 输入的格式错误，获取 shareId 失败")
	}
}

// 从 HTML 代码中根据 key 获取下面几种格式的 value
//
// 'key':123  =>  获得 "123"
//
// 'key':value  => 获得 "value"
//
// 'key':'str'  => 获得 "str"
func getValueKey(html string, key string) (string, error) {
	match := regexp.MustCompile(`'` + key + `':'?([^',]+)`).FindStringSubmatch(html)
	if len(match) == 0 {
		return "", errors.New("[getValueKey] 正则获取 " + key + " 失败")
	}
	return match[1], nil
}

// 自动获取 valueKey，然后根据 valueKey 获取值
func getValue(html string, key string) (string, error) {
	varName, err := getValueKey(html, key)
	if err != nil {
		return "", nil
	}
	match := regexp.MustCompile(`var ` + varName + ` = '(.*?)'`).FindStringSubmatch(html)
	if len(match) == 0 {
		return "", errors.New("[getValue] 正则获取 " + varName + "失败")
	}
	return match[1], nil
}

func fetchIframe(iframeURL string) (*DownloadInfo, error) {
	html, err := SendRequest(RequestConfig{
		URL: iframeURL,
	})

	if err != nil {
		return nil, err
	}
	signMatch := regexp.MustCompile(`'sign':'(.*?)'`).FindStringSubmatch(html)
	if len(signMatch) == 0 {
		return nil, errors.New("[fetchIframe] 获取 sign 失败")
	}
	sign := signMatch[1]
	signs, err := getValue(html, "signs")
	if err != nil {
		return nil, nil
	}
	websign, err := getValue(html, "websign")
	if err != nil {
		return nil, nil
	}
	websignkey, err := getValue(html, "websignkey")
	if err != nil {
		return nil, nil
	}
	params := url.Values{}
	params.Set("action", "downprocess")
	params.Set("signs", signs)
	params.Set("websign", websign)
	params.Set("websignkey", websignkey)
	params.Set("ves", "1")
	params.Set("sign", sign)
	postURLMatch := regexp.MustCompile(`url : '(.*?)'`).FindStringSubmatch(html)
	if len(postURLMatch) == 0 {
		return nil, errors.New("[fetchIframe] 获取 postURL 失败")
	}
	postURL := baseURL + postURLMatch[1]
	return ajaxm(postURL, params)
}

// 通过密码获取文件直链，需要先获取文件分享页的 HTML 源码
func fetchWithPassword(html string, password string) (*DownloadInfo, error) {
	sign, err := getValue(html, "sign")
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("action", "downprocess")
	params.Set("sign", sign)
	params.Set("p", password)
	params.Set("kd", "1")
	postURLMatch := regexp.MustCompile(`url : '(.*?)'`).FindStringSubmatch(html)
	if len(postURLMatch) == 0 {
		return nil, errors.New("[fetchWithPassword] 获取 postURL 失败")
	}
	postURL := baseURL + postURLMatch[1]
	downloadInfo, err := ajaxm(postURL, params)
	if err != nil {
		return nil, errors.New("[fetchWithPassword] 密码错误")
	}
	return downloadInfo, nil
}

// 获取文件夹最新的一个文件的信息，包含直链，urlOrId 是文件夹的分享链接或 ID，password 是访问密码
func GetLatestFile(shareURL string, password string) (*DownloadInfo, error) {
	fileList, err := GetFileList(shareURL, password, 0)
	if err != nil {
		return nil, err
	}
	if len(fileList) == 0 {
		return nil, errors.New("没有发现文件")
	}
	downloadInfo, err := GetDownloadInfo(fileList[0].ShareURL(), password)
	if err != nil {
		return nil, err
	}
	downloadInfo.FileInfo = fileList[0]
	return downloadInfo, nil
}

// 获取文件夹中指定页码的文件列表
//
// page 的值务必从 0 开始，每次只允许增长 1，不可以直接从 0 变为 2。
//
// 每次换页，务必暂停 1 秒以上。
func GetFileList(shareURL string, password string, page int) ([]FileInfo, error) {
	shareId, err := getShareId(shareURL)
	if err != nil {
		return nil, err
	}
	html, err := SendRequest(RequestConfig{
		URL: baseURL + "/" + shareId,
		Headers: map[string]string{
			"User-Agent": "go-lanzou",
		},
	})
	if err != nil {
		return nil, err
	}
	postURLMatch := regexp.MustCompile(`url : '(.*?)'`).FindStringSubmatch(html)
	if len(postURLMatch) == 0 {
		return nil, errors.New("[fetchWithPassword] 获取 postURL 失败")
	}
	postURL := baseURL + postURLMatch[1]
	// lx
	lx, err := getValueKey(html, "lx")
	if err != nil {
		return nil, err
	}
	// uid
	uid, err := getValueKey(html, "uid")
	if err != nil {
		return nil, err
	}
	// up
	up, err := getValueKey(html, "up")
	if err != nil {
		return nil, err
	}
	// fid
	fid, err := getValueKey(html, "fid")
	if err != nil {
		return nil, err
	}
	// rep
	rep, err := getValueKey(html, "rep")
	if err != nil {
		return nil, err
	}
	// t
	t, err := getValue(html, "t")
	if err != nil {
		return nil, err
	}
	// k
	k, err := getValue(html, "k")
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Set("lx", lx)
	params.Set("fid", fid)
	params.Set("uid", uid)
	params.Set("pg", strconv.Itoa(page+1))
	params.Set("rep", rep)
	params.Set("t", t)
	params.Set("k", k)
	params.Set("up", up)
	params.Set("pwd", password)
	// ls
	ls, err := getValueKey(html, "ls")
	if err == nil {
		params.Set("ls", ls)
	}
	return ajaxList(postURL, params)
}

func ajaxm(postURL string, params url.Values) (*DownloadInfo, error) {
	body, err := SendRequest(RequestConfig{
		URL:    postURL,
		Method: "POST",
		Body:   strings.NewReader(params.Encode()),
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
			"Referer":      baseURL,
		},
	})
	if err != nil {
		return nil, err
	}
	var resData struct {
		Dom string `json:"dom"`
		Url string `json:"url"`
	}
	err2 := json.Unmarshal([]byte(body), &resData)
	if err2 != nil {
		return nil, err2
	}
	downloadURL := resData.Dom + "/file/" + resData.Url
	return &DownloadInfo{
		URL: downloadURL,
	}, nil
}

func ajaxList(postURL string, params url.Values) ([]FileInfo, error) {
	body, err := SendRequest(RequestConfig{
		URL:    postURL,
		Method: "POST",
		Body:   strings.NewReader(params.Encode()),
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
			"Referer":      baseURL,
		},
	})
	if err != nil {
		return nil, err
	}

	var resData struct {
		Zt   int    `json:"zt"`
		Info string `json:"info"`
		List []struct {
			ShareId string `json:"id"`
			Name    string `json:"name_all"`
			Size    string `json:"size"`
			Date    string `json:"time"`
		} `json:"text"`
	}

	err2 := json.Unmarshal([]byte(body), &resData)
	if err2 != nil {
		return nil, errors.New("JSON 解析失败，可能是页码越界或访问密码错误，接口提示：" + resData.Info)
	}

	fileList := make([]FileInfo, len(resData.List))
	for index, item := range resData.List {
		fileList[index] = FileInfo{
			ShareId: item.ShareId,
			Name:    item.Name,
			Size:    item.Size,
			Date:    item.Date,
		}
	}
	return fileList, nil
}

// 文件基础信息，不含直链
type FileInfo struct {
	// 从分享链接提取的标识字符串
	ShareId string
	// 文件名称
	Name string
	// 文件大小
	Size string
	// 上传日期
	Date string
	// 访问密码
	Password string
}

// 文件分享链接，根据 baseURL 和 ShareId 构建而成
func (f FileInfo) ShareURL() string {
	return baseURL + "/" + f.ShareId
}

// 文件基础信息，加上直链
type DownloadInfo struct {
	FileInfo
	URL string
}
