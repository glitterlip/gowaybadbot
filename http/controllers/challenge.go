package controllers

import (
	"encoding/binary"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"goawaybot/http/respones"
	"goawaybot/models"
	"goawaybot/rules"
	"goawaybot/services"
	"goawaybot/store"
	"net"
	"net/http"
	"time"
)
import jsoniter "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Challenge(c echo.Context) error {
	t := c.QueryParam("type")
	challenge := &models.Challenge{
		Id:         uuid.New().String(),
		Code:       0,
		UserIp:     binary.BigEndian.Uint32(net.ParseIP(c.RealIP()).To4()),
		CreateTime: time.Now(),
		Secret:     "",
	}

	if t == "" {
		t = "1"
	}
	switch t {
	case "1":
		rule := rules.GetInputVerification()
		challenge.Type = models.TypeInput
		challenge.Rule = rule
	case "2":
		rule := rules.GetSlideVerification()
		challenge.Type = models.TypeSlide
		challenge.Rule = rule
	case "3":
		rule := rules.GetClickVerification()
		challenge.Type = models.TypeClick
		challenge.Rule = rule
	}
	store.Put(challenge.Id, challenge, 600)
	return c.Render(http.StatusOK, "challenge.html", map[string]interface{}{
		"challenge": challenge,
	})
}
func Check(c echo.Context) error {
	challenge := make(map[string]interface{})
	id := c.FormValue("challenge")
	bytes, err := store.Cache.Get([]byte(id))
	if err != nil {
		return respones.Error(c, services.RuleExpired, "expired")
	}

	json.Unmarshal(bytes, &challenge)

	//check ip

	//check if already successed
	if challenge["Code"].(float64) == float64(services.RulePassed) {
		return respones.Success(c, "success")
	}

	//check attempts
	if challenge["Attempts"].(float64) > 1 {
		return respones.Error(c, services.RuleExpired)
	}

	//check answer
	err = services.CheckAnswer(challenge, c.FormValue("answer"))

	if err != nil {
		return respones.Error(c, services.RuleFailed)
	}
	return respones.Success(c, "success")

}
func Verify(c echo.Context) error {
	id := c.QueryParam("challenge-id")
	bytes, err := store.Cache.Get([]byte(id))
	if err != nil {
		return respones.Error(c, services.RuleExpired, "not found")
	}

	challenge := make(map[string]interface{})

	json.Unmarshal(bytes, &challenge)
	if challenge["Code"].(float64) == float64(services.RulePassed) {
		return respones.Success(c, 0, "ok", challenge)
	}
	return respones.Error(c, services.RuleFailed, "not passed yet")
}

func Refresh(c echo.Context) error {
	oid := c.QueryParam("old-challenge")
	oldChallenge := make(map[string]interface{})

	bytes, err := store.Cache.Get([]byte(oid))
	if err != nil {
		return respones.Error(c, services.RuleExpired, "expired")
	}

	json.Unmarshal(bytes, &oldChallenge)

	challenge := &models.Challenge{
		Id:         uuid.New().String(),
		Code:       0,
		UserIp:     binary.BigEndian.Uint32(net.ParseIP(c.RealIP()).To4()),
		CreateTime: time.Now(),
		Secret:     "",
	}

	switch oldChallenge["Type"].(float64) {
	case float64(models.TypeInput):
		rule := rules.GetInputVerification()
		challenge.Type = models.TypeInput
		challenge.Rule = rule
	case float64(models.TypeSlide):
		rule := rules.GetSlideVerification()
		challenge.Type = models.TypeSlide
		challenge.Rule = rule
	case float64(models.TypeClick):
		rule := rules.GetClickVerification()
		challenge.Type = models.TypeClick
		challenge.Rule = rule
	}
	store.Put(challenge.Id, challenge, 600)

	challenge.Rule = challenge.Rule.(models.ToMapable).ToMapRule()
	return respones.Success(c, "ok", challenge)
}
