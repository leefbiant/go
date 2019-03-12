package help

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
)

/**
 * 首字母大写
 */
func Ucfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

/**
 * 首字母小写
 */
func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

/**
 * 随机字符串
 */
func RandomString(strlen int64) string {
	str := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := make([]byte, 0, strlen)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := int64(0); i < strlen; i++ {
		result = append(result, str[r.Intn(len(str))])
	}
	return string(result)
}

/**
 * 获取当前路径
 * @param {[type]} ) (string, error [description]
 */
func GetCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return "", errors.New(`error: Can't find "/" or "\".`)
	}
	return string(path[0 : i+1]), nil
}

/**
 * 利用反射调用数据处理方法
 */
func ReflectInterface(any interface{}, name string, args ...interface{}) []reflect.Value {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}

	if v := reflect.ValueOf(any).MethodByName(name); v.String() == "<invalid Value>" {
		return nil
	} else {
		return v.Call(inputs)
	}
}

/**
 * 蛇形字符串转换为驼峰型
 */
func CamelCase(str string) string {
	temp := strings.Split(str, "_")
	var upperStr string
	for k, v := range temp {
		if k != 0 {
			v = Ucfirst(v)
		}
		upperStr += v
	}
	return upperStr
}

/**
 * 获取结构体中字段名
 */
func GetStructFieldName(structName interface{}) []string {
	t := reflect.TypeOf(structName)
	if t.Kind() != reflect.Struct { // this type is not struct
		return nil
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		result = append(result, t.Field(i).Name)
	}
	return result
}

/**
 * 获取结构体中的Tag名
 */
func GetStructFieldTag(structName interface{}) map[string]string {
	t := reflect.TypeOf(structName)
	if t.Kind() != reflect.Struct { // this type is not struct
		return nil
	}
	fieldNum := t.NumField()
	result := make(map[string]string)
	for i := 0; i < fieldNum; i++ {
		tagName := t.Field(i).Name
		tags := strings.Split(string(t.Field(i).Tag), "\"")
		if len(tags) > 1 {
			result[tagName] = tags[1]
		}
	}
	return result
}

/**
 * 将数据解析到指定结构体中
 */
func ConvToStruct(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)
	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}
	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}
	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	var err error
	if structFieldType != val.Type() {
		val, err = TypeConversion(fmt.Sprintf("%v", value), structFieldType.String())
		if err != nil {
			return err
		}
	}
	structFieldValue.Set(val)
	return nil
}

/**
 * string类型转换成其它类型
 * @param {[type]} value string  [description]
 * @param {[type]} ntype string) (reflect.Value, error [description]
 */
func TypeConversion(value string, ntype string) (reflect.Value, error) {
	var (
		res interface{}
		err error
	)

	switch ntype {
	case "string":
		res = value
		break
	case "int":
		res, err = strconv.Atoi(value)
		break
	case "int8":
		i, er := strconv.ParseInt(value, 10, 64)
		res = int8(i)
		err = er
		break
	case "int16":
		i, er := strconv.ParseInt(value, 10, 64)
		res = int16(i)
		err = er
		break
	case "int32":
		i, er := strconv.ParseInt(value, 10, 64)
		res = int32(i)
		err = er
		break
	case "int64":
		i, er := strconv.ParseInt(value, 10, 64)
		res = int64(i)
		err = er
		break
	case "float32":
		i, er := strconv.ParseFloat(value, 32)
		res = float32(i)
		err = er
		break
	case "float64":
		var newVal float64
		newVal, err = strconv.ParseFloat(value, 64)
		res = float64(newVal)
		break
	case "bool":
		res, err = strconv.ParseBool(value)
		break
	case "uint":
		res, err = strconv.Atoi(value)
		break
	case "uint8":
		i, er := strconv.ParseUint(value, 10, 64)
		res = uint8(i)
		err = er
		break
	case "uint16":
		i, er := strconv.ParseUint(value, 10, 64)
		res = uint16(i)
		err = er
		break
	case "uint32":
		i, er := strconv.ParseUint(value, 10, 64)
		res = uint32(i)
		err = er
		break
	case "uint64":
		var newValue uint64
		newValue, err = strconv.ParseUint(value, 10, 64)
		res = uint64(newValue)
		break
	case "[][]float64":
		var rs [][]float64
		if len(value) == 2 {
			return reflect.ValueOf(rs), nil
		}
		value = string([]rune(value)[2 : len(value)-2])
		tmp := strings.Split(value, "] [")
		for _, v := range tmp {
			l := strings.Split(v, " ")
			ret := make([]float64, 0, len(l))
			for _, s := range l {
				t, err := TypeConversion(s, "float64")
				if err != nil {
					return reflect.ValueOf(t), err
				}
				ret = append(ret, t.Float())
			}
			rs = append(rs, ret)
			ret = ret[:0]
		}
		res = rs
		break
	default:
		res = value
		err = errors.New("unknow type: " + ntype)
		break
	}
	return reflect.ValueOf(res), err
}

type Accuracy func() float64

func (this Accuracy) Equal(a, b float64) bool {
	return math.Abs(a-b) < this()
}

func (this Accuracy) Greater(a, b float64) bool {
	return math.Max(a, b) == a && math.Abs(a-b) > this()
}

func (this Accuracy) Smaller(a, b float64) bool {
	return math.Max(a, b) == b && math.Abs(a-b) > this()
}

func (this Accuracy) GreaterOrEqual(a, b float64) bool {
	return math.Max(a, b) == a || math.Abs(a-b) < this()
}

func (this Accuracy) SmallerOrEqual(a, b float64) bool {
	return math.Max(a, b) == b || math.Abs(a-b) < this()
}

// 数据库查询
func QueryFormMysql(Db *sql.DB, sqlstr string) (*[]map[string]string, error) {
	// log.Printf("sql:%s", sqlstr)
	rows, err := Db.Query(sqlstr)
	if err != nil {
		fmt.Println("QueryFormMysql Query err:", err)
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		fmt.Println("Columns err:", err)
		return nil, err
	}

	values := make([][]byte, len(columns))
	scans := make([]interface{}, len(columns))
	list := make([]map[string]string, 0)

	for k := range values {
		scans[k] = &values[k]
	}

	for rows.Next() {
		if err := rows.Scan(scans...); err != nil {
			fmt.Println("Scan err:", err)
			return nil, err
		}
		row := make(map[string]string)
		for i, col := range values {
			if col != nil {
				row[columns[i]] = string(col)
			}
		}
		list = append(list, row)
	}
	return &list, nil
}

func String2Float(value string) float64 {
	if s, err := strconv.ParseFloat(value, 64); err == nil {
		return s
	}
	fmt.Println("err value:", value)
	return 0.0
}

func String2Int(value string) int64 {
	if s, err := strconv.ParseInt(value, 10, 64); err == nil {
		return s
	}
	fmt.Println("err value:", value)
	return 0
}
