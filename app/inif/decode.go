package inif

import (
	"fmt"
	"reflect"
	"strconv"
)

//Unmarshal - сереализация данных
func Unmarshal(data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)

	//fmt.Println(reflect.TypeOf(v))
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("type %v должен быть pointer", reflect.TypeOf(v))
	}
	fmt.Println(rv.Elem().Kind())

	row := splitRow(data)

	if err := readRow(row, rv.Elem()); err != nil {
		return err
	}
	switch rv.Elem().Kind() {
	case reflect.Array:
		readArray(data, rv)
	case reflect.Struct:

	default:
		return fmt.Errorf("%v не возможно сериализовать", rv.Kind())
	}
	readArray(data, rv)
	return nil
}

func readStruct(data []byte, rv reflect.Value) {}

func readArray(arr []byte, rv reflect.Value) {
	offset := 0
	for i, d := range arr {
		if isNewRow(d) || i == len(arr)-1 {
			fmt.Printf("row: %s\n", string(arr[offset:i+1]))
			listField := splitRow(arr[offset : i+1])
			if len(listField) == 0 {
				continue
			}
			fmt.Println(listField)
			offset = i + 1
		}
	}
}

//for struct
func readRow(splitRow [][]byte, v reflect.Value) error {
	if len(splitRow) > v.NumField() {
		return fmt.Errorf("в %v полей меньше чем в []byte", v.Kind())
	}
	for i := 0; i < v.NumField(); i++ {
		if err := setSimpleValue(v.Field(i), splitRow[i]); err != nil {
			return err
		}
	}
	return nil
}

func setSimpleValue(v reflect.Value, data []byte) error {
	if !v.CanAddr() || !v.CanSet() {
		return fmt.Errorf("%v not possible to set value", v.Kind())
	}
	fmt.Println(v.Kind())
	switch v.Kind() {
	case reflect.Bool:
		x, err := strconv.ParseBool(string(data))
		if err != nil {
			return err
		}
		v.SetBool(x)
	case reflect.String:
		v.SetString(string(data))
	case reflect.Int:
		x, err := strconv.Atoi(string(data))
		if err != nil {
			return err
		}
		v.SetInt(int64(x))
	case reflect.Int8, reflect.Uint8:
		x, err := strconv.ParseInt(string(data), 10, 8)
		if err != nil {
			return err
		}
		v.SetInt(x)
	case reflect.Int32:
		x, err := strconv.ParseInt(string(data), 10, 32)
		if err != nil {
			return err
		}
		v.SetInt(x)
	case reflect.Int64:
		x, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(x)
	default:
		return fmt.Errorf("%v not simple type", v.Kind())
	}
	return nil
}

func splitRow(row []byte) [][]byte {
	listField := make([][]byte, 0)
	offset := 0
	for i, r := range row {
		//fmt.Println(string(r))
		if startComment(r) {
			if len(row[offset:i]) > 1 {
				listField = append(listField, row[offset:i])
				fmt.Printf("value: %s, length: %d\n", string(row[offset:i]), len(row[offset:i]))
				offset = i + 1

			}
			break
		}
		if isSpace(r) {
			if len(row[offset:i]) == 0 {
				offset = i + 1
				continue
			}
			listField = append(listField, row[offset:i])
			fmt.Printf("value: %s, length: %d\n", string(row[offset:i]), len(row[offset:i]))

			offset = i + 1
		}
		if i == len(row)-1 && len(row[offset:i]) > 1 {
			listField = append(listField, row[offset:i])
			fmt.Printf("value: %s, length: %d\n", string(row[offset:i]), len(row[offset:i]))
		}
	}
	return listField
}

func isSpace(v byte) bool {
	return v == ' ' || v == '\t'
}

func isNewRow(v byte) bool {
	return v == '\r' || v == '\n'
}

func startComment(v byte) bool {
	return v == '#'
}
