package utils

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"log"
	"os"
)

func ImageToBase64(img image.Image) (string, error) {
	buffer := bytes.NewBuffer(nil)
	saveImage("1.jpg", img)
	if err := png.Encode(buffer, img); err != nil {
		return "", err
	}
	ret := base64.StdEncoding.EncodeToString(buffer.Bytes())
	return ret, nil
}

func saveImage(path string, img image.Image) (err error) {
	// 需要保存的文件
	imgfile, err := os.Create(path)
	defer imgfile.Close()
	err = png.Encode(imgfile, img)
	if err != nil {
		log.Fatal(err)
	}
	return
}
