package rules

import (
	"fmt"
	"github.com/fogleman/gg"
	"goawaybot/services"
	"image"
	"strings"
)

const (
	WordCount = 5
)

var (
	Words = []rune("天地玄黄宇宙洪荒日月盈昃辰宿列张寒来暑往秋收冬藏闰余成岁律召调阳云腾致雨露结为霜金生丽水玉出昆岗剑号巨阙珠称夜光果珍李柰菜重芥姜海咸河淡鳞潜羽翔龙师火帝鸟官人皇始制文字乃服衣裳推位让国有虞陶唐吊民伐罪周发殷汤坐朝问道垂拱平章爱育黎首臣伏戎羌遐迩壹体率宾归王鸣凤在树白驹食场化被草木赖及万方盖此身发四大五常恭惟鞠养岂敢毁伤女慕贞洁男效才良知过必改得能莫忘罔谈彼短靡恃己长信使可覆器欲难量墨悲丝染诗赞羔羊景行维贤克念作圣德建名立形端表正空谷传声虚堂习听祸因恶积福缘善庆尺辟非宝寸阴是竞资父事君曰严与敬孝当竭力忠则尽命临深履薄夙兴温清似兰斯馨如松之盛川流不息渊澄取映容止若思言辞安定笃初诚美慎终宜令荣业所基籍甚无竟学优登仕摄职从政存以甘棠去而益咏乐殊贵贱礼别尊卑上和下睦夫唱妇随")
)

type ClickRule struct {
	RuleImage image.Image `json:"-"`
	Answer    string
	Hint      string
	FielName  string
	Fontsize  int
}

func GetClickVerification() *ClickRule {
	click := &ClickRule{
		Fontsize: 40,
	}
	click.FielName = services.GetImageFileNameForRule("click")
	answer := GetClickRuleAnswer()
	SetClickImage(click, answer)
	SetClickAnswer(click, answer)
	return click
}

func GetClickRuleAnswer() string {
	b := strings.Builder{}
	length := len(Words)
	for i := 0; i < 5; i++ {
		index := Rand.Intn(length)
		b.WriteString(string(Words[index : index+1]))
	}
	return b.String()
}

func SetClickImage(rule *ClickRule, str string) {
	imgFile, _ := services.ImagesFs.Open(fmt.Sprintf("images/templates/click/%s", rule.FielName))
	img, _, _ := image.Decode(imgFile)
	defer imgFile.Close()
	ctx := gg.NewContextForImage(img)
	fontSize := rule.Fontsize
	ctx.SetFontFace(services.GetChineseFont())
	for i := 0; i < WordCount; i++ {
		x := float64(fontSize) * (float64(i)*1.5 + 1)
		y := float64(ctx.Height() / 2)
		ctx.SetRGBA255(0, 0, 0, 255)
		rule.writeWord(ctx, string([]rune(str)[i:i+1]), x, y)
	}
	rule.RuleImage = ctx.Image()
}

func (rule *ClickRule) writeWord(ctx *gg.Context, text string, x, y float64) {
	xfloat := x
	yfloat := y + 80 - Rand.Float64()*160
	//draw surround lines for debug
	//a := []float64{xfloat - ClickFontSize, yfloat - ClickFontSize, xfloat + ClickFontSize, yfloat - ClickFontSize}
	//b := []float64{xfloat - ClickFontSize, yfloat - ClickFontSize, xfloat - ClickFontSize, yfloat + ClickFontSize}
	//c := []float64{xfloat + ClickFontSize, yfloat + ClickFontSize, xfloat + ClickFontSize, yfloat - ClickFontSize}
	//d := []float64{xfloat - ClickFontSize, yfloat + ClickFontSize, xfloat + ClickFontSize, yfloat + ClickFontSize}
	ctx.ShearAbout(0.5, 0.5, xfloat, yfloat)
	ctx.DrawStringAnchored(text, xfloat, yfloat, 0.5, 0.5)
	ctx.Stroke()
	ctx.Identity()
	//ctx.DrawLine(a[0], a[1], a[2], a[3])
	//ctx.DrawLine(b[0], b[1], b[2], b[3])
	//ctx.DrawLine(c[0], c[1], c[2], c[3])
	//ctx.DrawLine(d[0], d[1], d[2], d[3])
	//ctx.DrawPoint(xfloat, yfloat, 4)
	//use answer as a temporary holder
	rule.Answer = rule.Answer + (fmt.Sprintf("%.f,%.f|", xfloat, yfloat))
}
func SetClickAnswer(rule *ClickRule, str string) {
	firstIndex := Rand.Intn(WordCount)
	secondIndex := Rand.Intn(WordCount)
	for firstIndex == secondIndex {
		secondIndex = Rand.Intn(WordCount)
	}
	//random pick two word as answer
	rule.Hint = fmt.Sprintf("请按顺序点击[%s,%s]", string([]rune(str)[firstIndex:firstIndex+1]), string([]rune(str)[secondIndex:secondIndex+1]))
	temp := strings.Split(rule.Answer, "|")
	rule.Answer = fmt.Sprintf("%s|%s", temp[firstIndex], temp[secondIndex])
}
func (rule *ClickRule) ToMapRule() map[string]interface{} {
	res := make(map[string]interface{})
	res["Hint"] = rule.Hint
	res["RuleImage"] = services.ImgToBase64(rule.RuleImage)
	return res
}
