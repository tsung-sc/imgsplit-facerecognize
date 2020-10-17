package libs

//科大讯飞
type Input struct {
	EnCoding string `json:"encoding"`
	Status   int    `json:"status"`
	Image    string `json:"image"`
}

type Payload struct {
	Input1 Input `json:"input1"`
	Input2 Input `json:"input2"`
}

type FaceCompareResult struct {
	EnCoding string `json:"encoding"`
	Compress string `json:"compress"`
	Format   string `json:"format"`
}

type S67c9c78c struct {
	ServiceKind       string            `json:"service_kind"`
	FaceCompareResult FaceCompareResult `json:"face_compare_result"`
}

type Parameter struct {
	S67c9c78c S67c9c78c `json:"s67c9c78c"`
}

type Header struct {
	AppId  string `json:"app_id"`
	Status int    `json:"status"`
}

type BodyRequest struct {
	Header    Header    `json:"header"`
	Parameter Parameter `json:"parameter"`
	Payload   Payload   `json:"payload"`
}

type BodyReturn struct {
	Header  HeaderReturn  `json:"header"`
	Payload PayloadReturn `json:"payload"`
}
type FaceCompareResultReturn struct {
	Compress string `json:"compress"`
	Encoding string `json:"encoding"`
	Format   string `json:"format"`
	Text     string `json:"text"`
}
type PayloadReturn struct {
	FaceCompareResultReturn FaceCompareResultReturn `json:"face_compare_result"`
}
type HeaderReturn struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Sid     string `json:"sid"`
}
type Text struct {
	Ret   int     `json:"ret"`
	Score float64 `json:"score"`
}

func GetXunfeiBody(img1path, img2path string) BodyRequest {
	img1data := Pic2Base64(img1path)
	img2data := Pic2Base64(img2path)
	body := BodyRequest{}
	body.Header.AppId = Xunfeiappid
	body.Header.Status = 3
	body.Parameter.S67c9c78c.FaceCompareResult.EnCoding = "utf8"
	body.Parameter.S67c9c78c.FaceCompareResult.Compress = "raw"
	body.Parameter.S67c9c78c.FaceCompareResult.Format = "json"
	body.Parameter.S67c9c78c.ServiceKind = "face_compare"
	body.Payload.Input1.EnCoding = "jpg"
	body.Payload.Input1.Status = 3
	body.Payload.Input1.Image = img1data
	body.Payload.Input2.EnCoding = "jpg"
	body.Payload.Input2.Status = 3
	body.Payload.Input2.Image = img2data
	return body
}

//百度
type BaiduRequest struct {
	Image           string `json:"image"`
	ImageType       string `json:"image_type"`
	FaceType        string `json:"face_type"`
	QualityControl  string `json:"quality_control"`
	LivenessControl string `json:"liveness_control"`
}

func GetBaiduBody(img1path, img2path string) []BaiduRequest {
	img1data := Pic2Base64(img1path)
	img2data := Pic2Base64(img2path)
	body := make([]BaiduRequest, 2)
	body[0].Image = img1data
	body[0].ImageType = "BASE64"
	body[0].FaceType = "LIVE"
	body[0].QualityControl = "LOW"
	body[0].LivenessControl = "NONE"
	body[1].Image = img2data
	body[1].ImageType = "BASE64"
	body[1].FaceType = "LIVE"
	body[1].QualityControl = "LOW"
	body[1].LivenessControl = "NONE"
	return body
}
