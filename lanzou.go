package main

import (
	"encoding/json"
	"errors"
	"net/url"
	"regexp"
	"strings"
)

var baseURL = "https://iuroc.lanzoue.com"

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

func GetDownloadURL(urlOrId string) (string, error) {
	fileId, err := getFileId(urlOrId)
	if err != nil {
		return "", err
	}
	html, err := SendRequest(RequestConfig{
		URL: baseURL + "/" + fileId,
		Headers: map[string]string{
			"User-Agent": "go-lanzou",
		},
	})
	if err != nil {
		return "", err
	}
	iframeURLMatch := regexp.MustCompile(`src="(\/fn\?[^"]{20,})"`).FindStringSubmatch(html)
	if len(iframeURLMatch) == 0 {
		return "", errors.New("[GetDownloadURL] 获取 iframeURL 失败")
	}
	iframeURL := baseURL + iframeURLMatch[1]
	return fetchIframe(iframeURL)
}

func fetchIframe(iframeURL string) (string, error) {
	html, err := SendRequest(RequestConfig{
		URL: iframeURL,
	})

	getValueKey := func(key string) (string, error) {
		match := regexp.MustCompile(`'` + key + `':([^,]+)`).FindStringSubmatch(html)
		if len(match) == 0 {
			return "", errors.New("[getValueKey] 正则获取 " + key + "失败")
		}
		return match[1], nil
	}
	getValue := func(key string) (string, error) {
		varName, err := getValueKey(key)
		if err != nil {
			return "", nil
		}
		match := regexp.MustCompile(`var ` + varName + ` = '(.*?)'`).FindStringSubmatch(html)
		if len(match) == 0 {
			return "", errors.New("[getValue] 正则获取 " + varName + "失败")
		}
		return match[1], nil
	}

	if err != nil {
		return "", err
	}
	signMatch := regexp.MustCompile(`'sign':'(.*?)'`).FindStringSubmatch(html)
	if len(signMatch) == 0 {
		return "", errors.New("[fetchIframe] 获取 sign 失败")
	}
	sign := signMatch[1]
	signs, err := getValue("signs")
	if err != nil {
		return "", nil
	}
	websign, err := getValue("websign")
	if err != nil {
		return "", nil
	}
	websignkey, err := getValue("websignkey")
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
	postURL := baseURL + postURLMatch[1]
	return ajaxm(postURL, params)
}

func ajaxm(postURL string, params url.Values) (string, error) {
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
