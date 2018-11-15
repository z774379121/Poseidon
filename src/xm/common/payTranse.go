package common

import (
	"github.com/smallnest/rpcx/log"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const RootName string = "xml"

const BufSize int = 32 * 1024

func GetMapFromStruct(inputstruct interface{}) map[string]string {
	v := reflect.ValueOf(inputstruct)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	m := make(map[string]string)
	if v.Kind() != reflect.Struct {
		log.Error(inputstruct)
		return nil
	}
	for i := 0; i < v.NumField(); i++ {
		//金额
		if v.Field(i).Kind() == reflect.Int {
			if v.Field(i).Interface().(int) == 0 {
				continue
			}
			m[v.Type().Field(i).Tag.Get(RootName)] = strconv.Itoa(v.Field(i).Interface().(int))
		} else if v.Field(i).Kind() == reflect.String {
			//空值不参与
			if v.Field(i).Interface().(string) == "" {
				continue
			}
			m[v.Type().Field(i).Tag.Get(RootName)] = v.Field(i).Interface().(string)
		}
	}
	return m
}

//生成随机字符串:长度要求小于32位
func RandString(length int) string {
	rand.Seed(time.Now().UnixNano())
	rs := make([]string, length)
	for start := 0; start < length; start++ {
		t := rand.Intn(3)
		if t == 0 {
			rs = append(rs, strconv.Itoa(rand.Intn(10)))
		} else if t == 1 {
			rs = append(rs, string(rand.Intn(26)+65))
		} else {
			rs = append(rs, string(rand.Intn(26)+97))
		}
	}
	return strings.Join(rs, "")
}
