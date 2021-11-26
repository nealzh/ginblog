package string_util

import (
	crand "crypto/rand"
	"github.com/google/uuid"
	"math"
	"math/big"
	mrand "math/rand"
	"time"
)

func HasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}

func HasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

//const GlobalGenBaseStr = "ABCDEFGHIGKLMNOPQRSTUVWXYZabcdefghigklmnopqrstuvwxyz0123456789"

const GlobalGenBaseStr = "abcdefghigklmnopqrstuvwxyz0123456789"

//var GlobalGenBaseStr = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GenRandomStr(rLen int, baseStr string) string {

	if baseStr == "" {
		baseStr = GlobalGenBaseStr
	}

	resultRunes := make([]rune, rLen)

	maxRandomIndex := len(baseStr)

	currentNanoTime := time.Now().UnixNano()

	currentRandNum, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))

	currentMRandSeed := currentNanoTime ^ currentRandNum.Int64()

	//fmt.Println(currentMRandSeed)
	mrand.Seed(currentMRandSeed)

	for i := 0; i < rLen; i++ {
		resultRunes[i] = rune(baseStr[mrand.Intn(maxRandomIndex)])
	}

	return string(resultRunes)
}

func GenUUIDStr(rLen int, baseStr string) string {
	return uuid.NewSHA1(uuid.NameSpaceX500, []byte(GenRandomStr(rLen, baseStr))).String()
}
