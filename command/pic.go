package command

import (
	"errors"
	"github.com/CodFrm/iotqq-plugins/model"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"strings"
)

func RotatePic(command []string, pic *model.PicInfo) ([]image.Image, error) {
	if len(command) > 4 {
		return nil, errors.New("命令过多")
	}
	resp, err := http.Get(pic.Url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}
	if img.Bounds().Max.X > 2048 || img.Bounds().Max.Y > 2048 {
		return nil, errors.New("图片过大")
	}
	retImage := make([]image.Image, 0)
	for _, v := range command {
		command2 := strings.Split(v, "+")
		if len(command2) > 4 {
			continue
		}
		tmpimg := copyImg(img)
		for _, v := range command2 {
			switch v {
			case "镜像":
				tmpimg = mirror(tmpimg)
			case "垂直":
				tmpimg = rotate90(tmpimg)
			case "翻转":
				tmpimg = rotate180(tmpimg)
			default:
				continue
			}
		}
		retImage = append(retImage, tmpimg)
	}
	return retImage, nil
}

func copyImg(m image.Image) image.Image {
	new := image.NewRGBA(image.Rect(0, 0, m.Bounds().Dx(), m.Bounds().Dy()))
	for x := m.Bounds().Min.X; x < m.Bounds().Max.X; x++ {
		for y := m.Bounds().Min.Y; y < m.Bounds().Max.Y; y++ {
			new.Set(x, y, m.At(x, y))
		}
	}
	return new
}

func rotate90(m image.Image) image.Image {
	rotate90 := image.NewRGBA(image.Rect(0, 0, m.Bounds().Dy(), m.Bounds().Dx()))
	for x := m.Bounds().Min.Y; x < m.Bounds().Max.Y; x++ {
		for y := m.Bounds().Max.X - 1; y >= m.Bounds().Min.X; y-- {
			rotate90.Set(m.Bounds().Max.Y-x, y, m.At(y, x))
		}
	}
	return rotate90
}

func mirror(m image.Image) image.Image {
	mirror := image.NewRGBA(image.Rect(0, 0, m.Bounds().Dx(), m.Bounds().Dy()))
	for x := m.Bounds().Min.X; x < m.Bounds().Max.X; x++ {
		for y := m.Bounds().Min.Y; y < m.Bounds().Max.Y; y++ {
			mirror.Set(x, y, m.At(m.Bounds().Max.X-x, y))
		}
	}
	return mirror
}

func rotate180(m image.Image) image.Image {
	rotate180 := image.NewRGBA(image.Rect(0, 0, m.Bounds().Dx(), m.Bounds().Dy()))
	for x := m.Bounds().Min.X; x < m.Bounds().Max.X; x++ {
		for y := m.Bounds().Min.Y; y < m.Bounds().Max.Y; y++ {
			rotate180.Set(m.Bounds().Max.X-x, m.Bounds().Max.Y-y, m.At(x, y))
		}
	}
	return rotate180
}
