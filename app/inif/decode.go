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

func readArray(arr []byte, rv reflect.Value) error {
	listRows := splitFile(arr)
	for _, r := range listRows {
		if len(bytes.Trim(r, " ")) == 0 {
			continue
		}
		arr := reflect.Indirect(rv)
		if !arr.CanSet() {
			return fmt.Errorf("non-set value for %v", rv)
		}
		sliceElemValue := reflect.Zero(rv.Type().Elem().Elem())
		newElem := reflect.ValueOf(&sliceElemValue)
		fmt.Println("rv.Type: ", rv.Type())
		fmt.Println("rv.Elem: ", rv.Elem().Type().Elem())
		fmt.Println("rv.Type.Elem: ", rv.Type().Elem())
		fmt.Println("rv.Type.Elem.Elem: ", rv.Type().Elem().Elem())

		fmt.Println(newElem.Elem().CanSet())
		if err := readRow(r, newElem.Elem()); err != nil {
			return err
		}

		arr.Set(reflect.Append(arr, newElem))
	}
	return nil
}

//for struct
func readRow(row []byte, v reflect.Value) error {
	listField := splitRow(row)
	fmt.Println("len v.NumField: ", (v.NumField()))
	// if len(listField) > v.NumField() {
	// 	return fmt.Errorf("the number of fields in the structure is %d, but should be %d",
	// 		v.NumField(), len(listField))
	// }
	for i := 0; i < v.NumField(); i++ {
		fmt.Println(v.Field(i))
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
			if len(row[offset:i]) == 0 {
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
