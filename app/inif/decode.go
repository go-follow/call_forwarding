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
		return fmt.Errorf("%v must pass a pointer, not a value, to Unmarshal", reflect.TypeOf(v))
	}
	fmt.Println(rv.Elem().Kind())

	switch rv.Elem().Kind() {
	case reflect.Array, reflect.Slice:
		readArray(data, rv)
	case reflect.Struct:
		if err := readRow(data, rv); err != nil {
			return err
		}
	default:
		return fmt.Errorf("%v не возможно сериализовать", rv.Kind())
	}
	readArray(data, rv)
	return nil
}

func readStruct(data []byte, rv reflect.Value) {}

func readArray(arr []byte, rv reflect.Value) {
	rows := splitFile(arr)
	fmt.Println(rows)
}

//for struct
func readRow(row []byte, v reflect.Value) error {
	listField := splitRow(row)
	//fmt.Println("field in file: ", len(listField))
	//fmt.Println("fiels in struct: ", v.Elem().NumField())
	if len(listField) > v.Elem().NumField() {
		return fmt.Errorf("the number of fields in the structure is %d, but should be %d", 
		v.Elem().NumField(), len(listField))
	}
	for i := 0; i < v.Elem().NumField(); i++ {
		if err := setSimpleValue(v.Elem().Field(i), listField[i]); err != nil {
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

func splitFile(data []byte) [][]byte {
	listField := make([][]byte, 0)
	offset := 0
	for i, d := range data {
		if isNewRow(d) {
			if len(data[offset:i]) == 0 {
				offset = i + 1
				continue
			}
			listField = append(listField, data[offset:i])
			fmt.Printf("value: %s, length: %d\n", string(data[offset:i]), len(data[offset:i]))

			offset = i + 1
		}
		if i == len(data)-1 && len(data[offset:i]) > 1 {
			listField = append(listField, data[offset:i])
			fmt.Printf("value: %s, length: %d\n", string(data[offset:i]), len(data[offset:i]))
		}
	}
	return listField
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
