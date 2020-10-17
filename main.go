package main

import (
	"fmt"
	"lianxi/rxfg/libs"
)

func main() {
	checkface := new(libs.FaceVerify)
	checkface.Faceimg1 = "D:/jaychou.jpg"
	checkface.Faceimg2 = "D:/jaychou2.jpg"
	// //科大讯飞测试
	// result, err := checkface.XunfeiCheckFace(checkface.Faceimg1, checkface.Faceimg2)
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(result)
	// }

	//百度测试
	// result, err := checkface.BaiduCheckFace(checkface.Faceimg1, checkface.Faceimg2)
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(result)
	// }

	//腾讯测试
	// result, err := checkface.TenCheckFace(checkface.Faceimg1, checkface.Faceimg2)
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(result)
	// }

	//阿里测试
	// result, err := checkface.AliCheckFace(checkface.Faceimg1, checkface.Faceimg2)
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(result)
	// }

	//总测试
	result, errinfo, err := checkface.CheckFace(2, checkface.Faceimg1, checkface.Faceimg2)
	if err != nil {
		fmt.Println(err)
		fmt.Println(errinfo)
	} else {
		fmt.Println(result)
		fmt.Println(errinfo)
	}

	// result := libs.GetAccessToken(libs.BaiduFaceVerifyKey, libs.BaiduFaceVerifyKeySecret)
	// fmt.Println(result)
	// _, err := libs.AliImgSplit("D:/jaychou.jpg")
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println("123")
	// }

	// _, errinfo, err := libs.ImageSplit(2, "D:/jaychou.jpg")
	// if err != nil {
	// 	fmt.Println(err)
	// 	fmt.Println(errinfo)
	// } else {
	// 	// fmt.Println(result)
	// 	fmt.Println(errinfo)
	// }
	// test := libs.Auth()
	// fmt.Println(test)
}
