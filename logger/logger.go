package logger

import (
	"fmt"
	"log"
	"os"
)

//Info - вывод лога уровня INFO
func Info(v ...interface{}) {
	l := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	if err := l.Output(2, fmt.Sprintln(v...)); err != nil {
		fmt.Println("ERROR: не удалось вывести сообщение в лог. Сообщение:", fmt.Sprintf("%v", v...))
	}
}

//Fatal - вывод лога уровня FATAL. Завершение работы сервиса
func Fatal(v ...interface{}) {
	l := log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)
	if err := l.Output(2, fmt.Sprintln(v...)); err != nil {
		fmt.Println("ERROR: не удалось вывести сообщение в лог. Сообщение:", fmt.Sprintf("%v", v...))
	}
	os.Exit(1)
}

//Error - вывод лога уровня ERROR
func Error(v ...interface{}) {
	l := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	if err := l.Output(2, fmt.Sprintln(v...)); err != nil {
		fmt.Println("ERROR: не удалось вывести сообщение в лог. Сообщение:", fmt.Sprintf("%v", v...))
	}
}

//Infof - форматированный вывод лога уровня INFO
func Infof(mes string, v ...interface{}) {
	l := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	if err := l.Output(2, fmt.Sprintf(mes, v...)); err != nil {
		fmt.Println("ERROR: не удалось вывести сообщение в лог. Сообщение:", fmt.Sprintf(mes, v...))
	}
}

//Errorf - форматированный вывод лога уровня ERROR
func Errorf(mes string, v ...interface{}) {
	l := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	if err := l.Output(2, fmt.Sprintf(mes, v...)); err != nil {
		fmt.Println("ERROR: не удалось вывести сообщение в лог. Сообщение:", fmt.Sprintf(mes, v...))
	}
}
