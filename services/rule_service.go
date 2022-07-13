package services

import (
	"errors"
	"fmt"
	"goawaybot/models"
	"goawaybot/store"
	"io/fs"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	RulePassed  = 1
	RuleExpired = 2
	RuleFailed  = 3
)

func CheckAnswer(challenge map[string]interface{}, answerStr string) error {
	ruleType := challenge["type"].(float64)
	switch ruleType {
	case models.TypeInput:
		input := challenge["rule"].(map[string]interface{})
		if strings.ToUpper(answerStr) == strings.ToUpper(input["answer"].(string)) {
			challenge["code"] = RulePassed
			store.Put(challenge["id"].(string), challenge, 600)
			return nil
		}
		challenge["attempts"] = challenge["attempts"].(float64) + 1
		store.Put(challenge["id"].(string), challenge, 600)
		return errors.New("fail")

	case models.TypeSlide:
		slide := challenge["rule"].(map[string]interface{})
		answerInt, _ := strconv.Atoi(answerStr)
		answer := float64(answerInt)
		if answer <= (slide["answer"].(float64)+slide["offset"].(float64)) && answer >= (slide["answer"].(float64)-slide["offset"].(float64)) {
			challenge["code"] = RulePassed
			store.Put(challenge["id"].(string), challenge, 600)
			return nil
		}
		challenge["attempts"] = challenge["attempts"].(float64) + 1
		store.Put(challenge["id"].(string), challenge, 600)
		return errors.New("fail")

	case models.TypeClick:
		click := challenge["rule"].(map[string]interface{})
		answers := strings.Split(answerStr, ",")
		x1, _ := strconv.Atoi(answers[0])
		y1, _ := strconv.Atoi(answers[1])
		x2, _ := strconv.Atoi(answers[2])
		y2, _ := strconv.Atoi(answers[3])
		ps := strings.Split(click["answer"].(string), "|")
		p1 := strings.Split(ps[0], ",")
		p2 := strings.Split(ps[1], ",")
		p1x, _ := strconv.Atoi(p1[0])
		p1y, _ := strconv.Atoi(p1[1])
		p2x, _ := strconv.Atoi(p2[0])
		p2y, _ := strconv.Atoi(p2[1])
		length := int(click["fontsize"].(float64))
		if x1 > p1x-length && x1 < p1x+length && y1 > p1y-length && y1 < p1y+length && x2 > p2x-length && x2 < p2x+length && y2 > p2y-length && y2 < p2y+length {
			challenge["code"] = RulePassed
			store.Put(challenge["id"].(string), challenge, 600)
			return nil
		}
		challenge["attempts"] = challenge["attempts"].(float64) + 1
		store.Put(challenge["id"].(string), challenge, 600)
		return errors.New("fail")
	}
	return nil
}
func GetImageFileNameForRule(rule string) string {
	var f []fs.DirEntry
	switch rule {
	case "slide":
		f, _ = ImagesFs.ReadDir("images/templates/slide")
	case "click":
		f, _ = ImagesFs.ReadDir("images/templates/click")
	}

	return fmt.Sprintf("%d.jpg", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(f)))
}
