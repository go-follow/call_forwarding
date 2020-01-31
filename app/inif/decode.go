package inif

import "reflect"
import "fmt"

//Unmarshal - сереализация данных
func Unmarshal(data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)
	
	//fmt.Println(reflect.TypeOf(v))
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("type %v должен быть pointer", reflect.TypeOf(v))
	}
	fmt.Println(reflect.TypeOf(v).String())
	switch rv.Kind() {
	case reflect.Ptr:
		//fmt.Println(rv.Addr())
	case reflect.Array:
		readArray(data, rv)
	case reflect.Struct:
		splitRow(data, rv)
	default:
		return fmt.Errorf("%v не возможно сериализовать", rv.Kind())
	}
	readArray(data, rv)
	return nil
}

func readArray(arr []byte, rv reflect.Value) {
	offset := 0
	for i, d := range arr {
		if isNewRow(d) || i == len(arr)-1 {			
			fmt.Printf("row: %s\n", string(arr[offset:i+1]))
			listField := splitRow(arr[offset:i+1], rv)
			if len(listField) == 0 {
				continue
			}
			fmt.Println(listField)
			offset = i + 1
		}
	}
}

func readRow(splitRow[][]byte, rv reflect.Value) error {
	if len(splitRow) > rv.NumField() {
		return fmt.Errorf("в %v полей меньше чем в []byte", rv.Kind())
	}
	for i := 0; i < rv.NumField(); i++ {
		
	}
}

func splitRow(row []byte, rv reflect.Value) [][]byte {
	listField := make([][]byte, 0)
	offset := 0	
	for i, r := range row {
		if startComment(r) {
			if len(row[offset:i+1]) > 1 {
				listField = append(listField, row[offset:i])
				offset = i + 1
			}
			break
		}
		if isSpace(r) {
			if len(row[offset:i+1]) == 1 && isSpace(row[offset : i+1][0]) {
				offset = i + 1
				continue
			}
			listField = append(listField, row[offset:i+1])			
			offset = i + 1
		}
		if i == len(row) - 1 && len(row[offset:i]) > 1 {
			listField = append(listField, row[offset:i])
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
