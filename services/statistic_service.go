package services

import (
	"goawaybot/store"
	"time"
)

const (
	ChallengeCoutnKey       = "challenge-count"
	ChallengePassedCountKey = "challenge-passed-count"
)

func InitStatistics() {
	store.Put(ChallengeCoutnKey, 0)
	store.Put(ChallengePassedCountKey, 0)
}
func NewChallenge() {
	store.Increment(ChallengeCoutnKey, 1)
}
func NewChallengePassed() {
	store.Increment(ChallengePassedCountKey, 1)
}
func GetSystemStatistics() map[string]interface{} {
	start, _ := time.Parse("2006 01-02 15:04:05", "2022 06-06 00:00:00")
	age := (time.Now().Unix() - start.Unix()) / 86400
	return map[string]interface{}{
		"runningtime":          age,
		"challengeCount":       string(store.Get(ChallengeCoutnKey, "0").([]uint8)),
		"challengePassedCount": string(store.Get(ChallengePassedCountKey, "0").([]uint8)),
	}
}
