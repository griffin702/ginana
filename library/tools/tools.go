package tools

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	Tools = New()
)

type Tool struct {
}

const GoTime = "2006-01-02 15:04:05"

/**
 * 返回单例实例
 * @method New
 */
func New() (t *Tool) {
	var once sync.Once
	once.Do(func() { //只执行一次
		t = &Tool{}
	})
	return t
}

/**
 * string转int
 */
func (t *Tool) StrToInt(str string, def int) (int, error) {
	i, err := strconv.Atoi(str)
	if err != nil {
		return def, err
	} else {
		return i, nil
	}
}

/**
 * string转int64
 */
func (t *Tool) StrToInt64(str string, def int64) (int64, error) {
	i, err := strconv.Atoi(str)
	if err != nil {
		return def, err
	} else {
		return int64(i), nil
	}
}

/**
 * string转uint
 */
func (t *Tool) StrToUint(str string, def uint) (uint, error) {
	i, err := strconv.Atoi(str)
	if err != nil {
		return def, err
	} else {
		return uint(i), nil
	}
}

/**
 * int转string
 */
func (t *Tool) IntToStr(i int) string {
	return strconv.Itoa(i)
}

/**
 * float转string
 */
func (t *Tool) FloatToStr(f float64) string {
	return strconv.FormatFloat(f, 'f', 5, 64)
}

/**
 * struct转成json
 */
func (t *Tool) StructToStr(data interface{}) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	} else {
		return string(b), nil
	}
}

/**
 * 结构体转换成map对象，常用于解决gorm使用结构体update时忽略零值的问题
 */
func (t *Tool) StructToMap(obj interface{}) (result map[string]interface{}, err error) {
	k := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if k.Kind() != reflect.Ptr {
		err = fmt.Errorf("type must be a pointer")
		return
	}

	if k.Elem().Kind() != reflect.Struct {
		err = fmt.Errorf("element type must be a struct")
		return
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
			log.Println(r)
			return
		}
	}()
	result = make(map[string]interface{})
	for i := 0; i < k.Elem().NumField(); i++ {
		name := k.Elem().Field(i).Name
		field := v.Elem().Field(i)
		switch name {
		case "ID", "CreatedAt", "UpdatedAt", "DeletedAt":
			continue
		}
		switch field.Kind() {
		case reflect.Slice, reflect.Struct, reflect.Ptr:
			continue
		default:
			result[name] = field.Interface()
		}
	}
	return
}

/**
 * 结构体赋值到另一个结构体
 */
func (t *Tool) StructToStruct(source, input interface{}) (err error) {
	sourceVal := reflect.ValueOf(source)
	inputVal := reflect.ValueOf(input)
	if sourceVal.Kind() != reflect.Ptr {
		return fmt.Errorf("source must be a pointer")
	}

	if inputVal.Kind() != reflect.Struct {
		return fmt.Errorf("input must be a struct")
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
			log.Println(r)
			return
		}
	}()
	for i := 0; i < inputVal.NumField(); i++ {
		for j := 0; j < sourceVal.Elem().NumField(); j++ {
			// 零值处理
			if inputVal.Field(i).IsZero() {
				continue
			}
			if t.equalName(sourceVal.Type().Elem().Field(j), inputVal.Type().Field(i)) {
				_ = t.assignment(sourceVal.Elem().Field(i), inputVal.Field(j))
				break
			}
		}
	}
	return nil
}

func (t *Tool) assignment(s1, s2 reflect.Value) (err error) {
	//fmt.Println(s1.Type(), s2.Type())
	if s1.Kind() != s2.Kind() {
		err = fmt.Errorf("字段种类必须一致")
		return
	}
	if s1.Type() == s2.Type() {
		s1.Set(s2)
	}
	return
}

func (t *Tool) equalName(s1, s2 reflect.StructField) bool {
	return s1.Name == s2.Name
}

/**
 * 生成随机字符串
 */
func (t *Tool) GetRandomString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ~!@#$%^&*()+[]{}/<>;:=.,?"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

//生成随机验证码
func (t *Tool) GetRandomInt(n int) string {
	const letterBytes = "0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

/**
 * 字符串截取
 */
func (t *Tool) SubString(str string, start, length int) string {
	if length == 0 {
		return ""
	}
	runeStr := []rune(str)
	lenStr := len(runeStr)

	if start < 0 {
		start = lenStr + start
	}
	if start > lenStr {
		start = lenStr
	}
	end := start + length
	if end > lenStr {
		end = lenStr
	}
	if length < 0 {
		end = lenStr + length
	}
	if start > end {
		start, end = end, start
	}
	return string(runeStr[start:end])
}

/**
 * TimeFormat
 */
func (t *Tool) TimeFormat(time *time.Time, format string) string {
	var datePatterns = []string{
		// year
		"Y", "2006", // A full numeric representation of a year, 4 digits   Examples: 1999 or 2003
		"y", "06", //A two digit representation of a year   Examples: 99 or 03

		// month
		"m", "01", // Numeric representation of a month, with leading zeros 01 through 12
		"n", "1", // Numeric representation of a month, without leading zeros   1 through 12
		"M", "Jan", // A short textual representation of a month, three letters Jan through Dec
		"F", "January", // A full textual representation of a month, such as January or March   January through December

		// day
		"d", "02", // Day of the month, 2 digits with leading zeros 01 to 31
		"j", "2", // Day of the month without leading zeros 1 to 31

		// week
		"D", "Mon", // A textual representation of a day, three letters Mon through Sun
		"l", "Monday", // A full textual representation of the day of the week  Sunday through Saturday

		// time
		"g", "3", // 12-hour format of an hour without leading zeros    1 through 12
		"G", "15", // 24-hour format of an hour without leading zeros   0 through 23
		"h", "03", // 12-hour format of an hour with leading zeros  01 through 12
		"H", "15", // 24-hour format of an hour with leading zeros  00 through 23

		"a", "pm", // Lowercase Ante meridiem and Post meridiem am or pm
		"A", "PM", // Uppercase Ante meridiem and Post meridiem AM or PM

		"i", "04", // Minutes with leading zeros    00 to 59
		"s", "05", // Seconds, with leading zeros   00 through 59

		// time zone
		"T", "MST",
		"P", "-07:00",
		"O", "-0700",

		// RFC 2822
		"r", "Mon, 02 Jan 2006 15:04:05 -0700",
	}
	replacer := strings.NewReplacer(datePatterns...)
	format = replacer.Replace(format)
	return time.Format(format)
}

func (t *Tool) StringFormatTime(timeLayout string) int64 {
	theTime, _ := time.Parse(GoTime, timeLayout)
	timeUnix := theTime.Unix()
	return timeUnix
}

/**
 * UTF82GB2312
 */
func (t *Tool) UTF82GB2312(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GB18030.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

/**
 * UTF82GBK
 */
func (t *Tool) UTF82GBK(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

/**
 * GBK2UTF8
 */
func (t *Tool) GBK2UTF8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

/**
 * RemoveRepeated
 */
func (t *Tool) RemoveRepeated(obj interface{}) (err error) {
	value := reflect.ValueOf(obj)
	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("obj must be a pointer")
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
			log.Println(r)
			return
		}
	}()
	l := value.Elem().Len()
	valueType := value.Elem().Type()
	result := reflect.New(valueType).Elem()
	isRepeat := false
	for i := 0; i < l; i++ {
		repeat := false
		for j := i + 1; j < l; j++ {
			if value.Elem().Index(i).Interface() == value.Elem().Index(j).Interface() {
				repeat = true
				isRepeat = true
				break
			}
		}
		if !repeat {
			result = reflect.Append(result, value.Elem().Index(i))
		}
	}
	if isRepeat {
		value.Elem().Set(result)
	}
	return
}

/**
 * PtrIsNil
 */
func (t *Tool) PtrIsNil(ptr interface{}) bool {
	vi := reflect.ValueOf(ptr)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}

// MD5HashString MD5哈希值
func (t *Tool) MD5HashString(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// BcryptHashGenerate 生成bctypt的哈希值
func (t *Tool) BcryptHashGenerate(s string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hash)
}

// BcryptHashCompare 对比bctypt的哈希值
func (t *Tool) BcryptHashCompare(current string, req string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(current), []byte(req))
	if err == nil {
		return true
	}
	return false
}

// MustUUID 创建UUID，如果发生错误则抛出panic
func (t *Tool) MustUUID() string {
	v, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	return v.String()
}

func (t *Tool) HideStar(str string) (result string) {
	if str == "" {
		return "***"
	}
	if strings.Contains(str, "@") {
		res := strings.Split(str, "@")
		if len(res[0]) < 3 {
			resString := "***"
			result = resString + "@" + res[1]
		} else {
			res2 := t.Substr2(str, 0, 3)
			resString := res2 + "***"
			result = resString + "@" + res[1]
		}
		return result
	} else {
		reg := `^\d{9}$`
		rgx := regexp.MustCompile(reg)
		mobileMatch := rgx.MatchString(str)
		if mobileMatch {
			result = t.Substr2(str, 0, 3) + "****" + t.Substr2(str, 7, 11)
		} else {

			nameRune := []rune(str)
			lens := len(nameRune)
			if lens <= 1 {
				result = "***"
			} else if lens == 2 {
				result = string(nameRune[:1]) + "*"
			} else if lens == 3 {
				result = string(nameRune[:1]) + "*" + string(nameRune[2:3])
			} else if lens == 4 {
				result = string(nameRune[:1]) + "**" + string(nameRune[lens-1:lens])
			} else if lens > 4 {
				result = string(nameRune[:2]) + "***" + string(nameRune[lens-2:lens])
			}
		}
		return
	}
}

func (t *Tool) Substr2(str string, start int, end int) string {
	rs := []rune(str)
	return string(rs[start:end])
}

func (t *Tool) CheckPasswordLevel(ps string) error {
	if len(ps) < 8 {
		return fmt.Errorf("password len is < 8")
	}
	num := `[0-9]{1}`
	az := `[A-Za-z]{1}`
	//symbol := `[!@#~$%^&*()+|_]{1}`
	if b, err := regexp.MatchString(num, ps); !b || err != nil {
		return fmt.Errorf("password need num :%v", err)
	}
	if b, err := regexp.MatchString(az, ps); !b || err != nil {
		return fmt.Errorf("password need A_Z :%v", err)
	}
	//if b, err := regexp.MatchString(symbol, ps); !b || err != nil {
	//	return fmt.Errorf("password need symbol :%v", err)
	//}
	return nil
}

func (t *Tool) CheckUserName(uname string) (err error) {
	zz := `^\w{6,30}$`
	if b, err := regexp.MatchString(zz, uname); !b || err != nil {
		return fmt.Errorf("用户名不合法:%v", err)
	}
	zz = `^\d+$`
	if b, err := regexp.MatchString(zz, uname); b || err != nil {
		return fmt.Errorf("用户名不合法:%v", err)
	}
	return
}

func (t *Tool) CheckNickName(nikename string) (err error) {
	re := regexp.MustCompile("[\u0020-\u002F]|[\u003A-\u0040]|[\u005B-\u0060]|[\u00A0-\u00BF]")
	if err := re.MatchString(nikename); err == true {
		return fmt.Errorf("用户昵称不合法:%v", err)
	}
	return
}
func (t *Tool) RemoteIp(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get("ip"); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get("XForwardedFor"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}

	return remoteAddr
}

//email verify
func (t *Tool) VerifyEmailFormat(email string) bool {
	//pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`

	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

//mobile verify
func (t *Tool) VerifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

func (t *Tool) EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}

//金额转大写
func (t *Tool) AmountToCN(pMoney float64, pRound bool) string {
	var numberUpper = []string{"壹", "贰", "叁", "肆", "伍", "陆", "柒", "捌", "玖", "零"}
	var unit = []string{"分", "角", "圆", "拾", "佰", "仟", "万", "拾", "佰", "仟", "亿", "拾", "佰", "仟"}
	var regex = [][]string{
		{"零拾", "零"}, {"零佰", "零"}, {"零仟", "零"}, {"零零零", "零"}, {"零零", "零"},
		{"零角零分", "整"}, {"零分", "整"}, {"零角", "零"}, {"零亿零万零元", "亿元"},
		{"亿零万零元", "亿元"}, {"零亿零万", "亿"}, {"零万零元", "万元"}, {"万零元", "万元"},
		{"零亿", "亿"}, {"零万", "万"}, {"拾零圆", "拾元"}, {"零圆", "元"}, {"零零", "零"}}
	str, digitUpper, unitLen, round := "", "", 0, 0
	if pMoney == 0 {
		return "零"
	}
	if pMoney < 0 {
		str = "负"
		pMoney = math.Abs(pMoney)
	}
	if pRound {
		round = 2
	} else {
		round = 1
	}
	digitByte := []byte(strconv.FormatFloat(pMoney, 'f', round+1, 64)) //注意币种四舍五入
	unitLen = len(digitByte) - round

	for _, v := range digitByte {
		if unitLen >= 1 && v != 46 {
			s, _ := strconv.ParseInt(string(v), 10, 0)
			if s != 0 {
				digitUpper = numberUpper[s-1]

			} else {
				digitUpper = "零"
			}
			str = str + digitUpper + unit[unitLen-1]
			unitLen = unitLen - 1
		}
	}
	for i := range regex {
		reg := regexp.MustCompile(regex[i][0])
		str = reg.ReplaceAllString(str, regex[i][1])
	}
	if string(str[0:3]) == "元" {
		str = string(str[3:])
	}
	if string(str[0:3]) == "零" {
		str = string(str[3:])
	}
	return str
}
