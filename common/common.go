package common

import (
	"errors"
	"log"
	"reflect"
	"strconv"
	"time"
)

/*将数据库中返回记录生成的map[string]string，根据datamodels中的struct里面各个属性的类型进行转换并映射回结构体实例*/

func DataToStructByTagSql(data map[string]string, obj interface{}) {
	// ValueOf 返回一个新值，初始化为存储在接口 i 中的具体值; Elem 返回接口 v 包含的值或指针 v 指向的值
	objValue := reflect.ValueOf(obj).Elem()
	for i := 0; i < objValue.NumField(); i++ {
		// 获取sql中的值(对应该结构体属性的值)
		sqlVal := data[objValue.Type().Field(i).Tag.Get("sql")]
		// 先获取结构体中每个属性的名称和类型
		name := objValue.Type().Field(i).Name
		structFieldType := objValue.Field(i).Type()

		//获取变量类型，也可以直接写"string类型"
		val := reflect.ValueOf(sqlVal)
		var err error
		// 判断当前结构体属性的类型与sql传过来经过处理的map是否一致
		if structFieldType != val.Type() {
			// 类型转换
			val, err = TypeConversion(sqlVal, structFieldType.Name())
			if err != nil {
				log.Println(err)
			}
		}

		// 设置类性值
		objValue.FieldByName(name).Set(val)
	}
}

// 类型转换
func TypeConversion(value string, ntype string) (reflect.Value, error) {
	if ntype == "string" {
		return reflect.ValueOf(value), nil
	} else if ntype == "time.Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "int" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	} else if ntype == "int8" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	} else if ntype == "int32" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int64(i)), err
	} else if ntype == "int64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "float32" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	} else if ntype == "float64" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	}

	//else if .......增加其他一些类型的转换

	return reflect.ValueOf(value), errors.New("未知的类型：" + ntype)
}
