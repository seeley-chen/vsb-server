package idgen

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
)

/** 统一生成uuid */
func GenerateUuid() string {
	return uuid.New().String()
}

/** 统一生成字符串类型数字ID */
func GenerateStringId() string {
	n, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		return fmt.Sprintf("%d", time.Now().UnixMilli())
	}
	return fmt.Sprintf("%d%04d", time.Now().UnixMilli(), n.Int64())
}
