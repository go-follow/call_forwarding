package main

import (
	//"encoding/json"

	"flag"
	"io/ioutil"

	"github.com/go-follow/call_forwarding/app/config"
	"github.com/go-follow/call_forwarding/app/forward"
	"github.com/go-follow/call_forwarding/app/inif"
	"github.com/go-follow/call_forwarding/logger"
)

type Settings struct {
	ListnerIP   string
	ListnerPort byte
	ForwardIP   string
	ForwardPort int
}

func main() {
	path := getPath()
	c := make(chan int)
	data, err := ioutil.ReadFile("../config.conf")
	if err != nil {
		logger.Fatal(err)
	}
	s := Settings{}
	if err := inif.Unmarshal(data, &s); err != nil {
		logger.Fatal(err)
	}

	sList, err := config.ReadConfig(path)
	if err != nil {
		logger.Fatal("не удалось прочитать конфигурационный файл: ", err)
	}

	for _, s := range sList {
		f, err := forward.NewForward(s)
		if err != nil {
			logger.Error(err)
			continue
		}
		go f.StartListner()
	}
	logger.Info("call_forwarding успешно запущен")
	<-c
}

func getPath() string {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		return "config.conf"
	}
	return args[0]
}
