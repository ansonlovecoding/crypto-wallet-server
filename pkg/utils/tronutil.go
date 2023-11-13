package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/shopspring/decimal"
)

var base58Alphabets = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

// 1 TRX = 1000000 Sun
func TrxToSun(amount float64) *big.Int {
	amountDecimal := decimal.NewFromFloat(amount).Mul(decimal.NewFromInt(1000000))
	return amountDecimal.BigInt()
}

// 1 Sun = 0.000001 TRX
func SunToTrx(amount *big.Int) *decimal.Decimal {
	amountDecimal := decimal.NewFromBigInt(amount, 0).Div(decimal.NewFromInt(1000000))
	amountDecimal = amountDecimal.Truncate(6)
	return &amountDecimal
}

func Base58ToHexAddress(address string) string {
	return hex.EncodeToString(base58Decode([]byte(address)))
}

func HexToBase58Address(hexAddress string) (string, error) {
	hexAddress = strings.TrimPrefix(hexAddress, "0x")
	hexAddress = strings.TrimPrefix(hexAddress, "0X")
	addrByte, err := hex.DecodeString(hexAddress)
	if err != nil {
		return "", err
	}

	sha := sha256.New()
	sha.Write(addrByte)
	shaStr := sha.Sum(nil)

	sha2 := sha256.New()
	sha2.Write(shaStr)
	shaStr2 := sha2.Sum(nil)

	addrByte = append(addrByte, shaStr2[:4]...)
	return string(base58Encode(addrByte)), nil
}

func base58Encode(input []byte) []byte {
	x := big.NewInt(0).SetBytes(input)
	base := big.NewInt(58)
	zero := big.NewInt(0)
	mod := &big.Int{}
	var result []byte
	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		result = append(result, base58Alphabets[mod.Int64()])
	}
	reverseBytes(result)
	return result
}

func base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	for _, b := range input {
		charIndex := bytes.IndexByte(base58Alphabets, b)
		result.Mul(result, big.NewInt(58))
		result.Add(result, big.NewInt(int64(charIndex)))
	}
	decoded := result.Bytes()
	if input[0] == base58Alphabets[0] {
		decoded = append([]byte{0x00}, decoded...)
	}
	return decoded[:len(decoded)-4]
}

func reverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}
