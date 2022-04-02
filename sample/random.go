package sample

import (
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomStringFromSet(a ...string) string {
	n := len(a)
	if n == 0 {
		return ""
	}
	return a[rand.Intn(n)]
}

func RandomBool() bool {
	return rand.Intn(2) == 1
}

func RandomInt(min, max int) int {
	return min + rand.Int()%(max-min+1)
}

func RandomFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func RandomFloat32(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func RandomID() string {
	return uuid.New().String()
}

func RandomSummary() string {
	return RandomStringFromSet("refactoring done", "bug fix", "performace improvment")
}

func randomEmail() string {
	return gofakeit.Email()
}

func randomPassword() string {
	return gofakeit.Password(true, true, true, true, true, 8)
}
