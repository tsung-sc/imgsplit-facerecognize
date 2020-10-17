package libs

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tencenterr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	iai "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/iai/v20200303"
)

type FaceVerify struct {
	Faceimg1 string
	Faceimg2 string
}

//定义的服务器错误信息
var FaceVerifyServerInfo string = ""

// 输入服务器使用优先顺序、对比图片1的原始路径、对比图片2的原始路径，输出结果和错误信息
// 返回人脸对比结果、各服务器错误信息、函数执行错误信息。
func (f *FaceVerify) CheckFace(firstSplitImage int, url1Path, url2Path string) (result float64, serverInfo string, err error) {
	//当所有服务都无法验证时返回此错误
	var finalErr error = errors.New("全部服务器出错")
	//每当一个服务器无法验证时此参数便加一，直到所有服务器都尝试了一遍
	Loopcount += 1
	if Loopcount == 5 {
		return 0, FaceVerifyServerInfo, finalErr
	}
	//根据传入服务器的优先顺序选择处理图片的服务器
	switch firstSplitImage {
	case 1:
		result, err = f.TenCheckFace(url1Path, url2Path)
		//当此服务器报错时，将报错信息添加到FaceVerifyServerInfo上，程序最终处理完毕后将返回此信息
		FaceVerifyServerInfo = FaceVerifyServerInfo + "腾讯服务器错误：" + fmt.Sprintf("%s", err) + "  "
		break
	case 2:
		result, err = f.BaiduCheckFace(url1Path, url2Path)
		FaceVerifyServerInfo = FaceVerifyServerInfo + "百度服务器错误：" + fmt.Sprintf("%s", err) + "  "
		break
	case 3:
		result, err = f.AliCheckFace(url1Path, url2Path)
		FaceVerifyServerInfo = FaceVerifyServerInfo + "阿里服务器错误：" + fmt.Sprintf("%s", err) + "  "
		break
	case 4:
		result, err = f.XunfeiCheckFace(url1Path, url2Path)
		FaceVerifyServerInfo = FaceVerifyServerInfo + "讯飞服务器错误：" + fmt.Sprintf("%s", err) + "  "
	default:
		return 0, FaceVerifyServerInfo, errors.New("输入参数错误")
	}
	//每当一个服务器无法处理便将优先顺序加一，交给下一个服务器进行处理
	if err != nil {
		firstSplitImage = AddNum(firstSplitImage, 4)
		result, _, err = f.CheckFace(firstSplitImage, url1Path, url2Path)
	}
	//判断传回来的参数，如果所有服务器都无法处理，将0和前面定义好的错误以及各个服务器报错信息返回主程序
	if Loopcount == 5 {
		return 0, FaceVerifyServerInfo, finalErr
	}
	//如果任一服务器处理成功，便将返回的数值传回主程序或者传入递归的上一级程序，最终传回到主程序
	return result, FaceVerifyServerInfo, nil
}

func (f *FaceVerify) AliCheckFace(img1path, img2path string) (float64, error) {
	//阿里云需要将图片通过此步骤上传生成URL地址才能处理
	fileUrl1 := GetFilePath(img1path)
	fileUrl2 := GetFilePath(img2path)
	//以下为阿里云定义好的程序
	client, err := sdk.NewClientWithAccessKey("cn-shanghai", AliKey, AliKeySecret)
	if err != nil {
		return 0, err
	}
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "facebody.cn-shanghai.aliyuncs.com"
	request.Version = "2019-12-30"
	request.ApiName = "CompareFace"
	request.QueryParams["RegionId"] = "cn-shanghai"
	request.QueryParams["ImageURLA"] = fileUrl1
	request.QueryParams["ImageURLB"] = fileUrl2

	response, err := client.ProcessCommonRequest(request)
	//判断是否有错误信息，如果有的话返回给主函数此错误信息
	if err != nil {
		AliErr := fmt.Sprintf("%s", err)
		result := strings.Index(AliErr, "Message")
		AliErr = AliErr[result+9:]
		return 0, errors.New(AliErr)
	}
	//定义一个RESULT接收处理返回的结果
	var result map[string]map[string]float64
	json.Unmarshal([]byte(response.GetHttpContentString()), &result)
	if result != nil {
		return result["Data"]["Confidence"], nil
	}
	return 0, errors.New("接收到的数据为空")
}

func (f *FaceVerify) BaiduCheckFace(img1path, img2path string) (float64, error) {
	//先获取AccessToken
	AccessToken, err := GetAccessToken(BaiduFaceVerifyKey, BaiduFaceVerifyKeySecret)
	if err != nil {
		return 0, err
	}
	url := "https://aip.baidubce.com/rest/2.0/face/v3/match?access_token=" + AccessToken
	method := "POST"
	//获得传给百度服务器的REQUESTBODY
	bodyraw := GetBaiduBody(img1path, img2path)
	//绑定为JSON格式
	finalbody, _ := json.Marshal(bodyraw)
	payload := strings.NewReader(string(finalbody))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return 0, err
	}
	req.Header.Add("content-type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	//定义RESULT接收处理错误的返回结果，RESULT1接收返回数据
	var result map[string]interface{}
	var result1 map[string]map[string]float64
	json.Unmarshal(body, &result)
	if result != nil {
		if result["error_msg"].(string) != "SUCCESS" {
			return 0, errors.New((result["error_msg"]).(string))
		} else {
			json.Unmarshal(body, &result1)
			return result1["result"]["score"], nil
		}
	}
	return 0, errors.New("接收到的数据为空")
}

func (f *FaceVerify) XunfeiCheckFace(img1path, img2path string) (float64, error) {
	//讯飞需要经过前面认证服务器才能处理数据，assembleAuthUrl函数处理返回需要POST数据的URL地址
	url, err := assembleAuthUrl(XunfeihostUrl, XunfeiapiKey, XunfeiapiSecret)
	if err != nil {
		return 0, err
	}
	method := "POST"
	bodyraw := GetXunfeiBody(img1path, img2path)
	finalbody, _ := json.Marshal(bodyraw)
	payload := strings.NewReader(string(finalbody))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return 0, err
	}
	req.Header.Add("content-type", "application/json")

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	//先处理看是否鉴权认证出错
	var resErr map[string]string
	json.Unmarshal(body, &resErr)
	if resErr["message"] != "" {
		return 0, errors.New(resErr["message"])
	}
	//以下程序为处理返回的数据
	var bodyreturn BodyReturn
	json.Unmarshal(body, &bodyreturn)
	if bodyreturn.Header.Message != "success" {
		return 0, errors.New(bodyreturn.Header.Message)
	}
	resultRaw := bodyreturn.Payload.FaceCompareResultReturn.Text
	result, _ := base64.StdEncoding.DecodeString(resultRaw)
	var text Text
	json.Unmarshal(result, &text)
	return text.Score * 100, nil
}

func (f *FaceVerify) TenCheckFace(img1path, img2path string) (float64, error) {
	//先将传入的图片通过Pic2Base64函数处理成BASE64编码，再传入腾讯的服务器进行处理
	img1data := Pic2Base64(img1path)
	img2data := Pic2Base64(img2path)
	credential := common.NewCredential(
		TencentKey,
		TencentKeySecret,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "iai.tencentcloudapi.com"
	client, _ := iai.NewClient(credential, "ap-guangzhou", cpf)

	request := iai.NewCompareFaceRequest()

	params := "{\"ImageA\":\"" + img1data + "\",\"ImageB\":\"" + img2data + "\",\"FaceModelVersion\":\"3.0\",\"QualityControl\":3}"
	err := request.FromJsonString(params)
	if err != nil {
		return 0, err
	}
	response, err := client.CompareFace(request)
	// if _, ok := err.(*tencenterr.TencentCloudSDKError); ok {
	// 	fmt.Printf("An API error has returned: %s", err)
	// 	return
	// }
	if err != nil {
		return 0, errors.New(err.(*tencenterr.TencentCloudSDKError).Message)
	}
	rawdata := response.ToJsonString()
	//定义一个RESULT接收处理返回的结果
	var result map[string]map[string]float64
	json.Unmarshal([]byte(rawdata), &result)
	return result["Response"]["Score"], nil
}
