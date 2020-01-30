package inif

import "reflect"

import "fmt"

func Unmarshal(data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)
	fmt.Println("Kind: ", rv.Kind())
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("type %v должен быть pointer", reflect.TypeOf(v))
	}
	offset := 0
	for i, d := range data {
		// fmt.Printf("%d : %d-%c\n", i, d, d)
		if isNewRow(d) || i == len(data)-1 {
			fmt.Printf("row: %s\n", string(data[offset:i+1]))
			fillRow(data[offset:i+1], v)
			offset = i + 1
		}
	}
	return nil
}

func fillRow(row []byte, v interface{}) {
	offset := 0
	for i, r := range row {
		if startComment(r) {
			if len(row[offset:i+1]) > 1 {
				fmt.Println("field: ", string(row[offset:i])) //Тут нужно привести field к нужному полю
				offset = i + 1
			}
			break
		}
		if isSpace(r) {
			if len(row[offset:i+1]) == 1 && isSpace(row[offset : i+1][0]) {
				offset = i + 1
				continue
			}
			fmt.Println("field: ", string(row[offset:i+1])) //тут нужно привести field к нужному полю
			offset = i + 1
		}
	}
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
