package lanzou

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

const BASE_URL string = "https://iuroc.lanzoue.com"

func GetDownloadURL(urlOrId string, password string) (string, error) {
	fileId, err := getFileId(urlOrId)
	if err != nil {
		return "", err
	}
	html, err := SendRequest(RequestConfig{
		URL: BASE_URL + "/" + fileId,
		Headers: map[string]string{
			"User-Agent": "go-lanzou",
		},
	})
	if err != nil {
		return "", err
	}

	// 判断当前分享页是否需要密码
	if regexp.MustCompile(`class="passwdinput"`).MatchString(html) {
		if password == "" {
			fmt.Print("🔒 请输入文件访问密码：")
			fmt.Scan(&password)
		}
		return fetchWithPassword(html, password)
	}

	iframeURLMatch := regexp.MustCompile(`src="(\/fn\?[^"]{20,})"`).FindStringSubmatch(html)
	if len(iframeURLMatch) == 0 {
		return "", errors.New("[GetDownloadURL] 获取 iframeURL 失败")
	}
	iframeURL := BASE_URL + iframeURLMatch[1]
	return fetchIframe(iframeURL)
}

func getFileId(urlOrId string) (string, error) {
	match := regexp.MustCompile(`^https?://.*/([^/]+)`).FindStringSubmatch(urlOrId)
	if len(match) != 0 {
		return match[1], nil
	} else if regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(urlOrId) {
		return urlOrId, nil
	} else {
		return "", errors.New("输入的格式错误，获取 fileId 失败")
	}
}

func getValueKey(html string, key string) (string, error) {
	match := regexp.MustCompile(`'` + key + `':([^,]+)`).FindStringSubmatch(html)
	if len(match) == 0 {
		return "", errors.New("[getValueKey] 正则获取 " + key + "失败")
	}
	return match[1], nil
}
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

func fetchIframe(iframeURL string) (string, error) {
	html, err := SendRequest(RequestConfig{
		URL: iframeURL,
	})

	if err != nil {
		return "", err
	}
	signMatch := regexp.MustCompile(`'sign':'(.*?)'`).FindStringSubmatch(html)
	if len(signMatch) == 0 {
		return "", errors.New("[fetchIframe] 获取 sign 失败")
	}
	sign := signMatch[1]
	signs, err := getValue(html, "signs")
	if err != nil {
		return "", nil
	}
	websign, err := getValue(html, "websign")
	if err != nil {
		return "", nil
	}
	websignkey, err := getValue(html, "websignkey")
	if err != nil {
		return "", nil
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
		return "", errors.New("[fetchIframe] 获取 postURL 失败")
	}
	postURL := BASE_URL + postURLMatch[1]
	return ajaxm(postURL, params)
}

func ajaxm(postURL string, params url.Values) (string, error) {
	body, err := SendRequest(RequestConfig{
		URL:    postURL,
		Method: "POST",
		Body:   strings.NewReader(params.Encode()),
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
			"Referer":      BASE_URL,
		},
	})
	if err != nil {
		return "", err
	}
	var downloadInfo struct {
		Dom string `json:"dom"`
		Url string `json:"url"`
	}
	err2 := json.Unmarshal([]byte(body), &downloadInfo)
	if err2 != nil {
		return "", err2
	}
	downloadURL := downloadInfo.Dom + "/file/" + downloadInfo.Url
	return downloadURL, nil
}

func fetchWithPassword(html string, password string) (string, error) {
	sign, err := getValue(html, "sign")
	if err != nil {
		return "", err
	}
	params := url.Values{}
	params.Set("action", "downprocess")
	params.Set("sign", sign)
	params.Set("p", password)
	params.Set("kd", "1")
	postURLMatch := regexp.MustCompile(`url : '(.*?)'`).FindStringSubmatch(html)
	if len(postURLMatch) == 0 {
		return "", errors.New("[fetchIframe] 获取 postURL 失败")
	}
	postURL := BASE_URL + postURLMatch[1]
	downloadURL, err := ajaxm(postURL, params)
	if err != nil {
		return "", errors.New("密码错误")
	}
	return downloadURL, nil
}
