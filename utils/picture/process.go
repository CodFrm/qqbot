package picture

import (
	"image"
	"image/color"

	"github.com/nfnt/resize"
)

// CopyImg 复制一张新图片
func CopyImg(m image.Image) image.Image {
	new := image.NewRGBA(image.Rect(0, 0, m.Bounds().Dx(), m.Bounds().Dy()))
	for x := m.Bounds().Min.X; x < m.Bounds().Max.X; x++ {
		for y := m.Bounds().Min.Y; y < m.Bounds().Max.Y; y++ {
			new.Set(x, y, m.At(x, y))
		}
	}
	return new
}

// Mirror 镜像
func Mirror(m image.Image) image.Image {
	mirror := image.NewRGBA(image.Rect(0, 0, m.Bounds().Dx(), m.Bounds().Dy()))
	for x := m.Bounds().Min.X; x < m.Bounds().Max.X; x++ {
		for y := m.Bounds().Min.Y; y < m.Bounds().Max.Y; y++ {
			mirror.Set(x, y, m.At(m.Bounds().Max.X-x-1, y))
		}
	}
	return mirror
}

// Rotate90 旋转90°
func Rotate90(m image.Image) image.Image {
	rotate90 := image.NewRGBA(image.Rect(0, 0, m.Bounds().Dy(), m.Bounds().Dx()))
	for x := m.Bounds().Min.Y; x < m.Bounds().Max.Y; x++ {
		for y := m.Bounds().Max.X - 1; y >= m.Bounds().Min.X; y-- {
			rotate90.Set(m.Bounds().Max.Y-x-1, y, m.At(y, x))
		}
	}
	return rotate90
}

// Rotate180 旋转180°
func Rotate180(m image.Image) image.Image {
	rotate180 := image.NewRGBA(image.Rect(0, 0, m.Bounds().Dx(), m.Bounds().Dy()))
	for x := m.Bounds().Min.X; x < m.Bounds().Max.X; x++ {
		for y := m.Bounds().Min.Y; y < m.Bounds().Max.Y; y++ {
			rotate180.Set(m.Bounds().Max.X-x-1, m.Bounds().Max.Y-y-1, m.At(x, y))
		}
	}
	return rotate180
}

// Narrow 缩放
func Narrow(m image.Image, scale float32) image.Image {
	return resize.Resize(uint(float32(m.Bounds().Dx())*scale), uint(float32(m.Bounds().Dy())*scale), m, resize.Lanczos3)
}

// ColorReverse 颜色反转
func ColorReverse(m image.Image) image.Image {
	reverse := image.NewRGBA(image.Rect(0, 0, m.Bounds().Dx(), m.Bounds().Dy()))
	for x := m.Bounds().Min.X; x < m.Bounds().Max.X; x++ {
		for y := m.Bounds().Min.Y; y < m.Bounds().Max.Y; y++ {
			colorRgb := m.At(x, y)
			r, g, b, a := colorRgb.RGBA()
			reverse.Set(x, y, color.RGBA{uint8(255 - (r >> 8)), uint8(255 - (g >> 8)), uint8(255 - (b >> 8)), uint8(a >> 8)})
		}
	}
	return reverse
}

// Gray 灰白
func Gray(m image.Image) image.Image {
	gray := image.NewRGBA(image.Rect(0, 0, m.Bounds().Dx(), m.Bounds().Dy()))
	for x := m.Bounds().Min.X; x < m.Bounds().Max.X; x++ {
		for y := m.Bounds().Min.Y; y < m.Bounds().Max.Y; y++ {
			colorRgb := m.At(x, y)
			_, g, _, a := colorRgb.RGBA()
			g_uint8 := uint8(g >> 8)
			a_uint8 := uint8(a >> 8)
			gray.Set(x, y, color.RGBA{g_uint8, g_uint8, g_uint8, a_uint8})
		}
	}
	return gray
}
