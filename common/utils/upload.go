package utils

import (
	"context"
	"core/upload"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mszlu521/thunder/logs"
)

func Upload(ctx context.Context, userID uuid.UUID, file *multipart.FileHeader) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// 检查文件大小（限制为5MB）
	if file.Size > 5*1024*1024 {
		return "", fmt.Errorf("文件大小不能超过5MB")
	}
	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	// 打开文件
	src, err := file.Open()
	if err != nil {
		logs.Errorf("打开文件失败: %v", err)
		return "", fmt.Errorf("打开文件失败")
	}
	// 生成文件名
	filename := fmt.Sprintf("avatar/%s/%s%s", userID.String(), file.Filename, ext)
	defer src.Close()
	// 上传到阿里云OSS
	err = upload.AliyunOSSUpload.Upload(ctx, src, filename)
	if err != nil {
		logs.Errorf("上传文件失败: %v", err)
		return "", fmt.Errorf("上传文件失败")
	}
	// 获取公开访问URL
	url := upload.AliyunOSSUpload.GetPublicUrl(filename)
	return url, nil
}
