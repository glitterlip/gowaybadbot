package models

import "time"

const (
	TypeInput = 1
	TypeSlide = 2
	TypeClick = 3
)

type ToMapable interface {
	ToMapRule() map[string]interface{}
}
type Challenge struct {
	Id         string      `json:"id"`          //challenge id
	Type       int         `json:"type"`        //verification type
	Code       int         `json:"code"`        //verification result
	UserIp     string      `json:"user_ip"`     //challenger ip ipv6
	CreateTime time.Time   `json:"create_time"` //challenge time
	Secret     string      `json:"secret"`      //encrypt key reserve
	Rule       interface{} `json:"rule"`
	Attempts   int8        `json:"attempts"` //failed attempts
	Ticket     string      `json:"ticket"`   //verify ticket
}
