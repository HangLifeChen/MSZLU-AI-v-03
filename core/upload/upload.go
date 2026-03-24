package upload

import (
	"os"

	"github.com/mszlu521/thunder/upload"
)

var (
	AliyunOSSUpload *upload.AliyunOSSUpload
	QiniuUpload     *upload.QiniuUpload
)

func Init() {

	//初始化aliyun oss
	var err error
	AliyunOSSUpload, err = upload.InitAliyunOSSUpload(
		os.Getenv("ALIYUN_ACCESS_KEY_ID"),
		os.Getenv("ALIYUN_ACCESS_KEY_SECRET"),
		os.Getenv("ALIYUN_OSS_ENDPOINT"),
		os.Getenv("ALIYUN_OSS_BUCKET"))
	if err != nil {
		panic(err)
	}
	QiniuUpload, err = upload.InitQiniuUpload(
		os.Getenv("QINIU_REGION"),
		os.Getenv("QINIU_BUCKET"),
		os.Getenv("QINIU_ACCESS_KEY"),
		os.Getenv("QINIU_SECRET_KEY"))
	if err != nil {
		panic(err)
	}
}
