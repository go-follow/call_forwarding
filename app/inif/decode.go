package inif

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
)

//Unmarshal - сереализация данных
func Unmarshal(data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)
	fmt.Println("Type: ", reflect.TypeOf(v))
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("%v must pass a pointer, not a value, to Unmarshal", reflect.TypeOf(v))
	}

	switch rv.Elem().Kind() {
	case reflect.Array, reflect.Slice:
		return readArray(data, rv)			
	case reflect.Struct:
		return readRow(data, rv.Elem())
	default:
		return fmt.Errorf("not possible to serialize type %v", rv.Elem().Kind())
	}	
}

func readStruct(data []byte, rv reflect.Value) {}

func readArray(arr []byte, rv reflect.Value) error {
	listRows := splitFile(arr)
	for _, r := range listRows {
		if len(bytes.Trim(r, " ")) == 0 {						
			continue
		}
		arr := reflect.Indirect(rv)
		if !arr.CanSet() {
			return fmt.Errorf("%v not possible to set value", rv.Kind())
		}
		sliceElemValue := reflect.Zero(rv.Type().Elem().Elem())
		arr.Set(reflect.Append(arr, sliceElemValue))

		if err := readRow(r, arr.Index(arr.Len() - 1)); err != nil {
			return err
		}
	}
	return nil
}

//for struct
func readRow(row []byte, v reflect.Value) error {
	listField := splitRow(row)

	if len(listField) > v.NumField() {
		return fmt.Errorf("the number of fields in the structure is %d, but should be %d",
			v.NumField(), len(listField))
	}
	for i := 0; i < v.NumField(); i++ {
		if err := setSimpleValue(v.Field(i), listField[i]); err != nil {
			return err
		}
	}
	return nil
}

func setSimpleValue(v reflect.Value, data []byte) error {
	if !v.CanSet() {
		return fmt.Errorf("%v not possible to set value", v.Kind())
	}
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
		v.SetInt(int32(data))
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
	startComment := -1
	for i, d := range data {
		if isStartComment(d) && startComment < 0 {
			startComment = i //запоминаем старт для комментариев
			continue
		}
		if isNewRow(d) {
			if startComment > 0 {
				listField = append(listField, data[offset:startComment])
				// fmt.Printf("value: %s, length: %d\n", string(data[offset:startComment]), len(data[offset:startComment]))
				offset = i + 1
				startComment = -1 //сбрасывем старт для комментариев
				continue
			}

			if len(data[offset:i]) == 0 {
				offset = i + 1
				continue
			}
			listField = append(listField, data[offset:i])
			// fmt.Printf("value: %s, length: %d\n", string(data[offset:i]), len(data[offset:i]))
			offset = i + 1
		}
		if i == len(data)-1 && len(data[offset:i+1]) > 0 {
			if startComment > 0 {
				listField = append(listField, data[offset:startComment])
				// fmt.Printf("value: %s, length: %d\n", string(data[offset:startComment]), len(data[offset:startComment]))
				continue
			}
			listField = append(listField, data[offset:i+1])
			// fmt.Printf("value: %s, length: %d\n", string(data[offset:i + 1]), len(data[offset:i + 1]))
		}
	}
	return listField
}

func splitRow(row []byte) [][]byte {
	listField := make([][]byte, 0)
	offset := 0
	for i, r := range row {
		if isStartComment(r) {
			if len(row[offset:i]) > 1 {
				listField = append(listField, row[offset:i])
				offset = i + 1
			}
			break
		}
		if isSpace(r) {
			if len(row[offset:i]) == 0 || isNewRow(r) {
				offset = i + 1
				continue
			}
			listField = append(listField, row[offset:i])
			offset = i + 1
		}
		if i == len(row)-1 && len(row[offset:i+1]) > 0 {
			listField = append(listField, row[offset:i+1])
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

func isStartComment(v byte) bool {
	return v == '#'
}
