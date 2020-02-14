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
		return readStruct(data, rv.Elem())
	default:
		return fmt.Errorf("not possible to serialize type %v", rv.Elem().Kind())
	}	
}

func readStruct(data []byte, rv reflect.Value) error{
	listRows := split(data, isNewRow)
	if len (listRows) > 1 {
		return fmt.Errorf("there should be only one entry for the structure type, but there are %d of them", len(listRows))
	}
	if err := readRow(listRows[0], rv); err != nil {
		return err
	}
	return nil
}

func readArray(arr []byte, rv reflect.Value) error {
	listRows := split(arr, isNewRow)
	for _, r := range listRows {		
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
		if err := setSimpleValue(listField[i], v.Field(i)); err != nil {
			return err
		}
	}
	return nil
}

func setSimpleValue(data []byte, v reflect.Value) error {
	if !v.CanSet() {
		return fmt.Errorf("%v not possible to set value", v.Kind())
	}
	s := string(data)
	switch v.Kind() {
	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return fmt.Errorf("failed to cast '%s' to type bool: %v", s, err)
		}
		v.SetBool(b)
		return nil		
	case reflect.String:
		v.SetString(s)
		return nil
	case reflect.Int:		
		d, err := strconv.ParseInt(s, 10, 0)
		if err != nil {
			return fmt.Errorf("failed to cast '%s' to type int: %v", s, err)
		}
		v.SetInt(d)
		return nil
	case reflect.Uint:
		d, err := strconv.ParseUint(string(data), 10, 0)
		if err != nil {
			return fmt.Errorf("failed to cast '%s' to type uint: %v", s, err)
		}
		v.SetUint(d)
		return nil
	case reflect.Int8:
		d, err := strconv.ParseInt(s, 10, 8)
		if err != nil {
			return fmt.Errorf("failed to cast '%s' to type int8: %v", s, err)
		}
		v.SetInt(d)
		return nil
	case reflect.Uint8:
		d, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			return fmt.Errorf("failed to cast '%s' to type uint8: %v", s, err)
		}
		v.SetUint(d)
		return nil
	case reflect.Int16:
		d, err := strconv.ParseInt(s, 10, 16)
		if err != nil {
			return fmt.Errorf("failed to cast '%s' to type int16: %v", s, err)
		}
		v.SetInt(d)
		return nil
	case reflect.Uint16:
		d, err := strconv.ParseUint(s, 10, 16)
		if err != nil {
			return fmt.Errorf("failed to cast '%s' to type uint16: %v", s, err)
		}
		v.SetUint(d)
		return nil
	case reflect.Int32:
		d, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return fmt.Errorf("failed to cast '%s' to type int32: %v", s, err)
		}
		v.SetInt(d)
		return nil
	case reflect.Uint32:
		d, err := strconv.ParseUint(string(data), 10, 32)
		if err != nil {
			return fmt.Errorf("failed to cast '%s' to type uint32: %v", s, err)
		}
		v.SetUint(d)
		return nil
	case reflect.Int64:
		d, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return fmt.Errorf("failed to cast '%s' to type int64: %v", s, err)
		}
		v.SetInt(d)
		return nil
	case reflect.Uint64:
		d, err := strconv.ParseUint(string(data), 10, 64)
		if err != nil {
			return fmt.Errorf("failed to cast '%s' to type uint64: %v", s, err)
		}
		v.SetUint(d)
		return nil
	case reflect.Float32:
		d, err := strconv.ParseFloat(string(data), 32)
		if err != nil {
			return fmt.Errorf("failed to cast '%s' to type float32: %v", s, err)
		}
		v.SetFloat(d)
		return nil
	case reflect.Float64:
		d, err := strconv.ParseFloat(string(data), 64)
		if err != nil {
			return fmt.Errorf("failed to cast '%s' to type float64: %v", s, err)
		}
		v.SetFloat(d)
		return nil
	default:
		return fmt.Errorf("%v not simple type", v.Kind())
	}	
}

func split(data []byte, isSplit func(b byte) bool) [][]byte {
	listField := make([][]byte, 0)
	offset := 0
	startComment := -1
	for i, d := range data {
		if isStartComment(d) && startComment < 0 {
			startComment = i //запоминаем старт для комментариев
			continue
		}
		if isSplit(d) {
			if startComment >= 0 {
				listField = append(listField, data[offset:startComment])
				//fmt.Printf("value: %s, length: %d\n", string(data[offset:startComment]), len(data[offset:startComment]))
				offset = i + 1
				startComment = -1 //сбрасывем старт для комментариев
				continue
			}

			if len(data[offset:i]) == 0 {
				offset = i + 1
				continue
			}
			listField = append(listField, data[offset:i])
			//fmt.Printf("value: %s, length: %d\n", string(data[offset:i]), len(data[offset:i]))
			offset = i + 1
		}
		if i == len(data)-1 && len(data[offset:i+1]) > 0 {
			if startComment > 0 {
				listField = append(listField, data[offset:startComment])
				//fmt.Printf("value: %s, length: %d\n", string(data[offset:startComment]), len(data[offset:startComment]))
				continue
			}
			listField = append(listField, data[offset:i+1])
			//fmt.Printf("value: %s, length: %d\n", string(data[offset:i + 1]), len(data[offset:i + 1]))
		}
	}
	return filterRows(listField)
}

func filterRows(listRows [][]byte) [][]byte {
	filter := make([][]byte, 0)
	for _, r := range listRows {
		if len(bytes.Trim(r, " ")) == 0 {
			continue
		}
		filter = append(filter, r)
	}
	return filter
}

func splitRow(row []byte) [][]byte {
	listField := make([][]byte, 0)
	offset := 0
	for i, r := range row {
		if isStartComment(r) {
			//для коментариев после строки
			if len(row[offset:i]) > 1 {
				listField = append(listField, row[offset:i])
				offset = i + 1
				break
			} else {
				continue
			}			
		}
		if isSpace(r) {
			if len(row[offset:i]) == 0 || isNewRow(row[offset:i][0]) {
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
