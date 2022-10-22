package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alp = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alp)

	for i := 0; i < n; i++ {
		c := alp[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomAmount() int {
	return RandomInt(10, 1000)
}

func RandomEmail() string {
	return fmt.Sprintf("%v@gmail.com", RandomOwner())
}

func RandomCurrency() string {
	currencies := []string{USD, NGN, GHS}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
