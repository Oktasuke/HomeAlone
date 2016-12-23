package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/golang/glog"
	"net/http"
	"io/ioutil"
	"os"
	"encoding/csv"
	"io"
)

const (
	REGION string = "us-east-1"
	ENVY_LIMIT float64 = 100.0
)

func main() {
	// Rekognition clientの生成
	se, err := session.NewSession()
	if err != nil {
		glog.Fatal("failed to create session,", err)
	}
	svc :=rekognition.New(se,aws.NewConfig().WithRegion(REGION))

	// テスト用CSVの読み込み
	// フォーマット(幸せ写真セルフ判定値(0,1),"URL")
	fp, err := os.Open("resources/test.csv")
	if err != nil {
		glog.Fatal(err)
	}
	defer fp.Close()
	reader := csv.NewReader(fp)
	reader.Comma = ','
	elementNum, correctNum := 0.0, 0.0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			glog.Fatal(err)
		}
		var result = IsHeartBreak(record[1], svc) // 幸せフィルタの実行

		// 正答チェック
		if record[0] == fmt.Sprint(bool2int(result)){
			correctNum++
		}
		elementNum++
	}
	fmt.Printf("テスト画像数	: %f \n", elementNum)
	fmt.Printf("正答数	: %f \n", correctNum)
	fmt.Printf("正答率 	: %f \n", correctNum / elementNum)
}

/**
 * 幸せフィルタ
 */
func IsHeartBreak(url string, rec *rekognition.Rekognition) bool{
	re := executeDetectFaces(url, rec)
	// ２人であるか？
	if len(re.FaceDetails) == 2 {
		// 男女ペアであるか？
		if *re.FaceDetails[0].Gender.Value != *re.FaceDetails[1].Gender.Value {
			// HAPPY度の合計が閾値を超えているか?
			var gnh float64 = 0
			for i := 0; i < len(re.FaceDetails); i++ {
				for j := 0; j < len(re.FaceDetails[i].Emotions); j++{
					if(*re.FaceDetails[i].Emotions[j].Type == "HAPPY"){
						gnh += *re.FaceDetails[i].Emotions[j].Confidence
					}
				}
			}
			if ENVY_LIMIT < gnh {
				return true
			}
		}
	}
	return false
}

/**
 * Facial Analysisの実行
 */
func executeDetectFaces(url string, rec *rekognition.Rekognition) *rekognition.DetectFacesOutput{
	params := &rekognition.DetectFacesInput{
		Image: &rekognition.Image{
			Bytes: fetchByteImage(url),
		},
		Attributes: []*string{
			aws.String("ALL"),
		},
	}

	res, err := rec.DetectFaces(params)
	if err != nil {
		glog.Error(err.Error())
	}
	return res
}

/**
 * Detect Labelsの実行
 */
func executeDetectLabels(url string, svc *rekognition.Rekognition) *rekognition.DetectLabelsOutput{
	params := &rekognition.DetectLabelsInput{
		Image: &rekognition.Image{
			Bytes: fetchByteImage(url),
		},
	}
	res, err := svc.DetectLabels(params)

	if err != nil {
		glog.Error(err.Error())
		return
	}

	return res
}

/**
 * URLから画像を取得する
 */
func fetchByteImage(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		glog.Error(err)
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		glog.Error(err)
	}
	return bytes
}
/**
 * 正答判定用
 */
func bool2int(b bool) int{
	if b {
		return 1
	}
	return 0
}