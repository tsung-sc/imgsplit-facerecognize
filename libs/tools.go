package libs

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	viapiutil "github.com/alibabacloud-go/viapi-utils/client"
)

func AddNum(num int, loop int) int {
	if num == loop {
		num = 1
	} else {
		num += 1
	}
	return num
}

func Pic2Base64(url string) string {
	ff, _ := os.Open(url)
	defer ff.Close()
	sourcebuffer := make([]byte, 500000)
	n, _ := ff.Read(sourcebuffer)
	//base64压缩
	sourcestring := base64.StdEncoding.EncodeToString(sourcebuffer[:n])
	return sourcestring
}

func DeCodeBase64(data string, savePath string) {
	dist, _ := base64.StdEncoding.DecodeString(data)
	//写入新文件
	f, _ := os.OpenFile(savePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer f.Close()
	f.Write(dist)
}

//科大讯飞鉴权认证
func assembleAuthUrl(hosturl string, apiKey, apiSecret string) (string, error) {
	ul, err := url.Parse(hosturl)
	if err != nil {
		return "", err
	}
	//签名时间
	date := time.Now().UTC().Format(time.RFC1123)
	//date = "Tue, 28 May 2019 09:10:42 MST"
	//参与签名的字段 host ,date, request-line
	signString := []string{"host: " + ul.Host, "date: " + date, "POST " + ul.Path + " HTTP/1.1"}
	//拼接签名字符串
	sgin := strings.Join(signString, "\n")
	// fmt.Println(sgin)
	//签名结果
	sha := HmacWithShaTobase64("hmac-sha256", sgin, apiSecret)
	// fmt.Println(sha)
	//构建请求参数 此时不需要urlencoding
	authUrl := fmt.Sprintf("hmac username=\"%s\", algorithm=\"%s\", headers=\"%s\", signature=\"%s\"", apiKey,
		"hmac-sha256", "host date request-line", sha)

	//将请求参数使用base64编码
	authorization := base64.StdEncoding.EncodeToString([]byte(authUrl))

	v := url.Values{}
	v.Add("host", ul.Host)
	v.Add("date", date)
	v.Add("authorization", authorization)
	//将编码后的字符串url encode后添加到url后面
	callurl := hosturl + "?" + v.Encode()
	return callurl, nil
}
func HmacWithShaTobase64(algorithm, data, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	encodeData := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}

//阿里云
func GetFilePath(imgPath string) (file1Path string) {
	accessKeyId := tea.String(AliKey)
	accessKeySecret := tea.String(AliKeySecret)
	fileUrl := tea.String(imgPath)
	fileLoadAddress, err := viapiutil.Upload(accessKeyId, accessKeySecret, fileUrl)
	if err != nil {
		return ""
	}
	return *fileLoadAddress
}

//百度
func GetAccessToken(id, secret string) (string, error) {
	url := "https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=" + id + "&client_secret=" + secret
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return "", err
	}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var result map[string]string
	json.Unmarshal(body, &result)
	if result != nil {
		if result["error"] != "" {
			return "", errors.New(result["error"])
		} else {
			return result["access_token"], nil
		}
	}
	return "", errors.New("鉴权认证返回的数据为空")
}
