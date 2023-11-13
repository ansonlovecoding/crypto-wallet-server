package utils

import (
	"Share-Wallet/pkg/common/constant"
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/big"
	"math/rand"
	"net/http"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/shopspring/decimal"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
)

func GetKeysFromMap(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func FindMaxInMap(params map[string]int) (k string) {
	var v = 0
	for key, score := range params {
		if score > v {
			k = key
			v = score
		}
	}
	return k
}
func ArrayColumn(array interface{}, key string) (result map[string]interface{}, err error) {
	result = make(map[string]interface{})
	t := reflect.TypeOf(array)
	v := reflect.ValueOf(array)
	if t.Kind() != reflect.Slice {
		return nil, nil
	}
	if v.Len() == 0 {
		return nil, nil
	}

	for i := 0; i < v.Len(); i++ {
		indexv := v.Index(i)
		if indexv.Type().Kind() != reflect.Struct {
			return nil, nil
		}
		mapKeyInterface := indexv.FieldByName(key)
		if mapKeyInterface.Kind() == reflect.Invalid {
			return nil, nil
		}
		mapKeyString, err := interfaceToString(mapKeyInterface.Interface())
		if err != nil {
			return nil, err
		}
		result[mapKeyString] = indexv.Interface()
	}
	return result, err
}

func interfaceToString(v interface{}) (result string, err error) {
	switch reflect.TypeOf(v).Kind() {
	case reflect.Int64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		result = fmt.Sprintf("%v", v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result = fmt.Sprintf("%v", v)
	case reflect.String:
		result = v.(string)
	default:
		err = nil
	}
	return result, err
}
func ArrayKey(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// copy a by b  b->a
func CopyStructFields(a interface{}, b interface{}, fields ...string) (err error) {
	return copier.Copy(a, b)
}

func Wrap(err error, message string) error {
	return errors.Wrap(err, "==> "+printCallerNameAndLine()+message)
}

func WithMessage(err error, message string) error {
	return errors.WithMessage(err, "==> "+printCallerNameAndLine()+message)
}

func printCallerNameAndLine() string {
	pc, _, line, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name() + "()@" + strconv.Itoa(line) + ": "
}

func GetSelfFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return cleanUpFuncName(runtime.FuncForPC(pc).Name())
}
func cleanUpFuncName(funcName string) string {
	end := strings.LastIndex(funcName, ".")
	if end == -1 {
		return ""
	}
	return funcName[end+1:]
}

// Get the intersection of two slices
func Intersect(slice1, slice2 []uint32) []uint32 {
	m := make(map[uint32]bool)
	n := make([]uint32, 0)
	for _, v := range slice1 {
		m[v] = true
	}
	for _, v := range slice2 {
		flag, _ := m[v]
		if flag {
			n = append(n, v)
		}
	}
	return n
}

// Get the diff of two slices
func Difference(slice1, slice2 []uint32) []uint32 {
	m := make(map[uint32]bool)
	n := make([]uint32, 0)
	inter := Intersect(slice1, slice2)
	for _, v := range inter {
		m[v] = true
	}
	for _, v := range slice1 {
		if !m[v] {
			n = append(n, v)
		}
	}

	for _, v := range slice2 {
		if !m[v] {
			n = append(n, v)
		}
	}
	return n
}

// Get the intersection of two slices
func IntersectString(slice1, slice2 []string) []string {
	m := make(map[string]bool)
	n := make([]string, 0)
	for _, v := range slice1 {
		m[v] = true
	}
	for _, v := range slice2 {
		flag, _ := m[v]
		if flag {
			n = append(n, v)
		}
	}
	return n
}

// Get the diff of two slices
func DifferenceString(slice1, slice2 []string) []string {
	m := make(map[string]bool)
	n := make([]string, 0)
	inter := IntersectString(slice1, slice2)
	for _, v := range inter {
		m[v] = true
	}
	for _, v := range slice1 {
		if !m[v] {
			n = append(n, v)
		}
	}

	for _, v := range slice2 {
		if !m[v] {
			n = append(n, v)
		}
	}
	return n
}
func OperationIDGenerator() string {
	return strconv.FormatInt(time.Now().UnixNano()+int64(rand.Uint32()), 10)
}

func RemoveRepeatedStringInList(slc []string) []string {
	var result []string
	tempMap := map[string]byte{}
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l {
			result = append(result, e)
		}
	}
	return result
}

func GetSuperGroupTableName(groupID string) string {
	return constant.SuperGroupChatLogsTableNamePre + groupID
}
func GetErrSuperGroupTableName(groupID string) string {
	return constant.SuperGroupErrChatLogsTableNamePre + groupID
}

func ParseResponse(response *http.Response) (map[string]interface{}, error) {
	var result map[string]interface{}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err = json.Unmarshal(body, &result)
	}
	return result, err
}

func String2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func CompressStr(str string) string {
	if str == "" {
		return ""
	}
	//匹配一个或多个空白符的正则表达式
	reg := regexp.MustCompile("\\s+")
	return reg.ReplaceAllString(str, "")
}

func RemoveDuplicatesAndEmpty(adders []string) []string {
	result := make([]string, 0, len(adders))
	temp := map[string]struct{}{}
	for _, item := range adders {
		if _, ok := temp[item]; !ok {
			if item != "" {
				temp[item] = struct{}{}
				result = append(result, item)
			}
		}
	}
	return result
}
func RandomSecretString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, n)
	rand.Seed(time.Now().Unix())
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
func EncryptAES(plainText string, key string) (string, error) {

	c, err := aes.NewCipher([]byte(key))

	if err != nil {
		fmt.Println(err)
	}

	gcm, err := cipher.NewGCM(c)

	if err != nil {
		fmt.Println(err)
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(crand.Reader, nonce); err != nil {
		fmt.Println(err)
	}

	enc := gcm.Seal(nonce, nonce, []byte(plainText), nil)
	return string(enc), nil

}

func DecryptAES(enc string, key string) (string, error) {

	encbyte := []byte(enc)
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Println(err)
	}
	gcm, err := cipher.NewGCM(c)
	nonceSize := gcm.NonceSize()
	if len(encbyte) < nonceSize {
		fmt.Println(err)
	}
	nonce, ciphertext := encbyte[:nonceSize], encbyte[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println(err)
	}
	Decrypted := string(plaintext)
	return Decrypted, nil
}

func GetCoinName(coinType uint8) string {
	switch coinType {
	case constant.BTCCoin:
		return "BTC"
	case constant.ETHCoin:
		return "ETH"
	case constant.USDTERC20:
		return "USDT-ERC20"
	case constant.TRX:
		return "TRX"
	case constant.USDTTRC20:
		return "USDT-TRC20"
	default:
		return ""
	}
}

func GetCoinType(coinType string) int {
	switch coinType {
	case "BTC":
		return constant.BTCCoin
	case "ETH":
		return constant.ETHCoin
	case "USDT-ERC20":
		return constant.USDTERC20
	case "TRX":
		return constant.TRX
	case "USDT-TRC20":
		return constant.USDTTRC20
	default:
		return 0
	}
}

func GetFundLogStateToString(state int8) string {
	switch state {
	case constant.FundlogFailed:
		return "failed"
	case constant.FundlogSuccess:
		return "success"
	case constant.FundlogPending:
		return "pending"
	default:
		return ""
	}
}

// 保留两位小数，舍弃尾数，无进位运算
// 主要逻辑就是先乘，trunc之后再除回去，就达到了保留N位小数的效果
func FormatFloat(num float64, decimal int) (float64, error) {
	// 默认乘1
	d := float64(1)
	if decimal > 0 {
		// 10的N次方
		d = math.Pow10(decimal)
	}
	// math.trunc作用就是返回浮点数的整数部分
	// 再除回去，小数点后无效的0也就不存在了
	res := strconv.FormatFloat(math.Trunc(num*d)/d, 'f', -1, 64)
	return strconv.ParseFloat(res, 64)
}

func ConvertUSDTBalance(balance string) string {
	balanceDecimal, _ := decimal.NewFromString(balance)
	//keep 6 decimal
	balanceDecimal = balanceDecimal.Div(decimal.NewFromInt(1000000)).Truncate(6)
	//tmpBalance, _ := strconv.ParseFloat(balance, 64)
	//
	//tmpBalance = (tmpBalance / 1000000)
	//usdtBalance, _ := FormatFloat(tmpBalance, 6)
	return balanceDecimal.String()
}

func ConvertFloatUSDTToBigInt(amount float64) *big.Int {
	amountDecimal := decimal.NewFromFloat(amount).Mul(decimal.NewFromInt(1000000))
	return amountDecimal.BigInt()
}

func ConvertBigUSDT2Float(amount *big.Int) float64 {
	//6位小数
	tmpBalance := float64(amount.Uint64())
	tmpBalance = (tmpBalance / 1000000)
	usdtBalance, _ := FormatFloat(tmpBalance, 6)
	return usdtBalance
}

// prec: numbers of decimal
func Float64WithoutRound(amount string, prec int) string {
	amountDecimal, _ := decimal.NewFromString(amount)
	amountDecimal = amountDecimal.Truncate(int32(prec))
	amountStr := amountDecimal.String()
	return amountStr
}
func Wei2Eth_str(amount *big.Int) string {
	compact_amount := big.NewInt(0)
	reminder := big.NewInt(0)
	divisor := big.NewInt(1e18)
	compact_amount.QuoRem(amount, divisor, reminder)
	return fmt.Sprintf("%v.%018s", compact_amount.String(), reminder.String())
}

// convert hex string to decimal
func HexToDecimal(hex string) (*decimal.Decimal, error) {
	hex = strings.TrimPrefix(hex, "0x")
	hex = strings.TrimPrefix(hex, "0X")
	amountInt, err := strconv.ParseUint(hex, 16, 64)
	if err != nil {
		return nil, err
	}
	amountDecimal := decimal.NewFromInt(int64(amountInt))
	return &amountDecimal, nil
}

func WrapErrorWithCode(info constant.ErrInfo) error {
	return errors.Wrap(errors.New(info.ErrMsg), fmt.Sprintf("%d", info.ErrCode))
}
