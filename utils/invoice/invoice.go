package invoice

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateNumber() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	timestamp := time.Now().UnixMilli()
	randomNumber := r.Intn(10000)
	return fmt.Sprintf("INV-%d-%04d", timestamp, randomNumber)
}
