package services

import (
	"bytes"
	"embed"
	"encoding/base64"
	"fmt"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"image"
	"image/png"
	"math/rand"
	"os"
	"time"
)

var ImagesFs embed.FS

func GetChineseFont() font.Face {
	fontBytes, _ := ImagesFs.ReadFile("resources/fonts/chinese.ttf")
	f, _ := truetype.Parse(fontBytes)
	face := truetype.NewFace(f, &truetype.Options{
		Size: 50,
	})
	return face
}
func GetRegularFont() font.Face {
	fontBytes, _ := ImagesFs.ReadFile("resources/fonts/Hack-Regular.ttf")
	f, _ := truetype.Parse(fontBytes)
	face := truetype.NewFace(f, &truetype.Options{
		Size: 50,
	})
	return face
}
func ImgToBase64(img image.Image) string {
	r := bytes.Buffer{}
	err := png.Encode(&r, img)
	f, _ := os.Create(fmt.Sprintf("images/tmp/%d%d.png", time.Now().UnixNano(), rand.Intn(1000)))
	png.Encode(f, img)
	if err != nil {
		return ""
	}
	res := "data:image/png;base64," + base64.StdEncoding.EncodeToString(r.Bytes())
	return res
}
