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
	Id         string    //challenge id
	Type       int       //verification type
	Code       int       //verification result
	UserIp     string    //challenger ip ipv6
	CreateTime time.Time //challenge time
	Secret     string    //encrypt key reserve
	Rule       interface{}
	Attempts   int8   //failed attempts
	Ticket     string //verify ticket
}
