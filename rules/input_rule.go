package rules

//inspired by https://learnku.com/articles/44827

import (
	"github.com/fogleman/gg"
	"goawaybot/services"
	"image"
	"image/color"
	"math/rand"
	"strings"
	"time"
)

const (
	Chars = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz23456789" //remove 0 1 L I O o
)

var (
	BackgroundColor = color.RGBA{
		R: 255, G: 255, B: 255, A: 1,
	}
	Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type InputRule struct {
	RuleImage image.Image `json:"-"`
	Answer    string
	Hint      string
}

func GetInputVerification() *InputRule {
	rule := &InputRule{}
	rule.Answer = GetInputRuleAnswer()
	rule.RuleImage = GetInputRuleImage(rule)
	return rule
}

func GetInputRuleAnswer() string {
	b := strings.Builder{}
	length := len(Chars)
	for i := 0; i < 4; i++ {
		index := Rand.Intn(length)
		b.WriteString(Chars[index : index+1])
	}
	return b.String()
}
func GetInputRuleImage(rule *InputRule) image.Image {
	ctx := gg.NewContext(200, 100)
	ctx.SetColor(BackgroundColor)
	ctx.Clear()
	MakeNoise(ctx)
	length := len(rule.Answer)
	fontSize := float64(ctx.Height() / 2)
	ctx.SetFontFace(services.GetRegularFont())
	for i := 0; i < length; i++ {
		r, g, b, _ := getRandColor(100)
		ctx.SetRGBA255(r, g, b, 255)
		fontPosX := float64(ctx.Width()/length*i) + fontSize*0.4
		writeChar(ctx, rule.Answer[i:i+1], fontPosX, float64(ctx.Height()/2))
	}

	return ctx.Image()

}
func MakeNoise(ctx *gg.Context) {
	width := ctx.Width()
	height := ctx.Height()
	for i := 0; i < 5; i++ {
		x1, y1 := getRandPos(width, height)
		x2, y2 := getRandPos(width, height)
		r, g, b, a := getRandColor(255)
		w := float64(Rand.Intn(3) + 3)
		ctx.SetRGBA255(r, g, b, a)
		ctx.SetLineWidth(w)
		ctx.DrawLine(x1, y1, x2, y2)
		ctx.Stroke()
	}
}

func getRandPos(width, height int) (x float64, y float64) {
	x = Rand.Float64() * float64(width)
	y = Rand.Float64() * float64(height)
	return x, y
}
func getRandColor(maxColor int) (r, g, b, a int) {
	r = int(uint8(Rand.Intn(maxColor)))
	g = int(uint8(Rand.Intn(maxColor)))
	b = int(uint8(Rand.Intn(maxColor)))
	a = int(uint8(Rand.Intn(255)))
	return r, g, b, a
}

func writeChar(ctx *gg.Context, text string, x, y float64) {
	xfloat := 5 - Rand.Float64()*10 + x
	yfloat := 5 - Rand.Float64()*10 + y
	radians := 40 - Rand.Float64()*80
	ctx.RotateAbout(gg.Radians(radians), x, y)
	ctx.DrawStringAnchored(text, xfloat, yfloat, 0.2, 0.5)
	ctx.RotateAbout(-1*gg.Radians(radians), x, y)
	ctx.Stroke()
}
func (rule *InputRule) ToMapRule() map[string]interface{} {
	res := make(map[string]interface{})
	res["Hint"] = rule.Hint
	res["RuleImage"] = services.ImgToBase64(rule.RuleImage)
	return res
}
