package rules

import (
	"fmt"
	"github.com/fogleman/gg"
	"goawaybot/services"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"math/rand"
	"os"
	"time"
)

const (
	SlideWidth  = 320
	SlideHeight = 240
)

var (
	SlideHoleColor = color.RGBA{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	}
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}
type SlideRule struct {
	BackgroundImage image.Image `json:"-"`
	SealImage       image.Image `json:"-"`
	Answer          int
	Offset          int
	FielName        string
}

func GetSlideVerification() *SlideRule {

	slide := &SlideRule{}
	slide.Answer = 60 + rand.New(rand.NewSource(time.Now().UnixNano())).Intn(SlideWidth-50-60) //[60,270)
	slide.Offset = 5
	slide.FielName = services.GetImageFileNameForRule("slide")
	SetSlideImages(slide)
	return slide
}

func SetSlideImages(rule *SlideRule) {
	backgroundPath := fmt.Sprintf("images/caches/slide/%s_background_%d.png", rule.FielName, rule.Answer)
	sealPath := fmt.Sprintf("images/caches/slide/%s_seal_%d.png", rule.FielName, rule.Answer)
	//check cache
	backgroundImage, err := os.Open(backgroundPath)
	defer backgroundImage.Close()
	if err == nil {
		rule.BackgroundImage, _, _ = image.Decode(backgroundImage)
		sealImage, _ := os.Open(sealPath)
		rule.SealImage, _, _ = image.Decode(sealImage)
		defer sealImage.Close()
		return
	}
	imgFile, err := os.Open(fmt.Sprintf("images/templates/slide/%s", rule.FielName))
	img, _, _ := image.Decode(imgFile)
	defer imgFile.Close()

	//create seal
	sealImg := img.(SubImager).SubImage(image.Rect(rule.Answer, 50, rule.Answer+50, 100))
	rule.SealImage = sealImg

	//create background
	ctx := gg.NewContextForImage(img)
	ctx.DrawRectangle(float64(rule.Answer), 50, 50, 50)
	ctx.SetRGB(255, 255, 255)
	ctx.Fill()
	rule.BackgroundImage = ctx.Image()

	//cache
	sealReader, _ := os.Create(sealPath)
	backgroundReader, _ := os.Create(backgroundPath)
	png.Encode(sealReader, sealImg)
	png.Encode(backgroundReader, rule.BackgroundImage)

}
func (rule *SlideRule) ToMapRule() map[string]interface{} {
	res := make(map[string]interface{})
	res["BackgroundImage"] = services.ImgToBase64(rule.BackgroundImage)
	res["SealImage"] = services.ImgToBase64(rule.SealImage)
	res["Offset"] = rule.Offset
	res["FielName"] = rule.FielName

	return res
}
