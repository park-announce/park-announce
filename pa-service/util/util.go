package util

import (
	"math/rand"

	uuid "github.com/satori/go.uuid"
)

func GenerateRandomNumber(max int) int {
	return rand.Intn(max)
}

func GenerateGuid() string {
	return uuid.NewV4().String()
}
