package libs

//
import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"

	bda "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bda/v20200324"
	tencenterr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

//定义的服务器错误信息
var ImgSplitServerErr = ""

//分割图片的函数，需要传入服务器优先顺序：firstSplitImage，原始图片路径：urlPath
func ImageSplit(firstSplitImage int, urlPath string) (newImgBase64, serverInfo string, err error) {
	var finalErr error = errors.New("全部服务器出错")
	Loopcount += 1
	if Loopcount == 4 {
		return "", ImgSplitServerErr, finalErr
	}
	switch firstSplitImage {
	case 1:
		newImgBase64, err = TencentImgSplit(urlPath)
		ImgSplitServerErr = ImgSplitServerErr + "腾讯服务器错误：" + fmt.Sprintf("%s", err) + "  "
		break
	case 2:
		newImgBase64, err = BaiduImgSplit(urlPath)
		ImgSplitServerErr = ImgSplitServerErr + "百度服务器错误：" + fmt.Sprintf("%s", err) + "  "
		break
	case 3:
		newImgBase64, err = AliImgSplit(urlPath)
		ImgSplitServerErr = ImgSplitServerErr + "阿里服务器错误：" + fmt.Sprintf("%s", err) + "  "
		break
	default:
		return "", ImgSplitServerErr, errors.New("输入参数错误")
	}
	if err != nil {
		firstSplitImage = AddNum(firstSplitImage, 3)
		newImgBase64, _, err = ImageSplit(firstSplitImage, urlPath)
	}
	if Loopcount == 4 {
		return "", ImgSplitServerErr, finalErr
	}
	return newImgBase64, ImgSplitServerErr, nil
}

func TencentImgSplit(picPath string) (string, error) {
	credential := common.NewCredential(
		TencentKey,
		TencentKeySecret,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "bda.tencentcloudapi.com"
	client, err := bda.NewClient(credential, "ap-guangzhou", cpf)
	if err != nil {
		return "", err
	}

	request := bda.NewSegmentPortraitPicRequest()
	params := "{\"Image\":\"" + Pic2Base64(picPath) + "\"}"
	err = request.FromJsonString(params)
	if err != nil {
		return "", err
	}
	response, err := client.SegmentPortraitPic(request)
	// if _, ok := err.(*tencenterr.TencentCloudSDKError); ok {
	// 	fmt.Printf("An API error has returned: %s", err)
	// }
	if err != nil {
		return "", errors.New(err.(*tencenterr.TencentCloudSDKError).Message)
	}
	//以下为处理返回数据的代码
	rawdata := response.ToJsonString()
	var obj map[string]map[string]string
	json.Unmarshal([]byte(rawdata), &obj)
	return obj["Response"]["ResultImage"], nil
}

func BaiduImgSplit(picPath string) (string, error) {
	//先获取AccessToken
	AccessToken, err := GetAccessToken(BaiduImgSplitKey, BaiduImgSplitKeySecret)
	if err != nil {
		return "", err
	}
	url := "https://aip.baidubce.com/rest/2.0/image-classify/v1/body_seg?access_token=" + AccessToken
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("image", Pic2Base64(picPath))
	err1 := writer.Close()
	if err1 != nil {
		return "", err1
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	//处理返回的数据
	var obj map[string]interface{}
	json.Unmarshal(body, &obj)
	if obj["error_code"] != nil {
		errCodeStr := obj["error_msg"].(string)
		return "", errors.New(errCodeStr)
	}
	return obj["foreground"].(string), nil
}

func AliImgSplit(picPath string) (string, error) {
	//将图片上传至阿里服务器
	fileUrl := GetFilePath(picPath)
	client, err := sdk.NewClientWithAccessKey("cn-shanghai", AliKey, AliKeySecret)
	if err != nil {
		return "", err
	}

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "imageseg.cn-shanghai.aliyuncs.com"
	request.Version = "2019-12-30"
	request.ApiName = "SegmentBody"
	request.QueryParams["ImageURL"] = fileUrl
	request.QueryParams["RegionId"] = "cn-shanghai"

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		AliErr := fmt.Sprintf("%s", err)
		result := strings.Index(AliErr, "Message")
		AliErr = AliErr[result+9:]
		return "", errors.New(AliErr)
	}
	// fmt.Print(response.GetHttpContentString())

	//处理返回的数据
	var result map[string]map[string]string
	json.Unmarshal([]byte(response.GetHttpContentString()), &result)

	//通过Http获取处理完的图片并编码成BASA64
	resp, err1 := http.Get(result["Data"]["ImageURL"])
	body, err2 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		return "", err1
	}
	if err2 != nil {
		return "", err2
	}
	sourcestring := base64.StdEncoding.EncodeToString(body)
	return sourcestring, nil
}
