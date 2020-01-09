package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/go-follow/call_forwarding/app/models"
	"gopkg.in/ini.v1"
)

//ReadConfig чтение конфигурационного файла ini, если файла не существует он создается со значениями по умолчанию и сообщением выхода из программы
func ReadConfig(pathConf string) ([]*models.Settings, error) {
	data, err := ioutil.ReadFile(pathConf)
	if err != nil {
		return nil, err
	}
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("not exist configuration file config.conf: %v", err)
	}
	arrFile := strings.Split(string(data), "\n")
	if len(arrFile) == 0 {
		return nil, fmt.Errorf("empty config.conf")
	}
	listSettings := make([]*models.Settings, 0)

	for _, a := range arrFile {
		if strings.Trim(a, " ") == "" {
			continue
		}
		main, comment := mainSplitComment(a)
		if main == "" {
			continue
		}
		arrMain := strings.Split(main, " ")
		if len(arrMain) < 4 {
			return nil, fmt.Errorf("not valid settings for row- %s. Example: '172.22.2.60 7371 192.168.41.26 3050'", main)
		}
		listerIP := strings.Trim(arrMain[0], " ")
		forwardIP := strings.Trim(arrMain[2], " ")
		portListn := strings.Trim(arrMain[1], " ")
		portForw := strings.Trim(arrMain[3], " ")
		portListn = strings.Trim(portListn, "\r")
		portForw = strings.Trim(portForw, "\r")
		portListner, err := strconv.Atoi(portListn)
		if err != nil {
			return nil, fmt.Errorf("port for listner not integer: %v", err)
		}
		portForward, err := strconv.Atoi(portForw)
		if err != nil {
			return nil, fmt.Errorf("port for forward not integer: %v", err)
		}
		s := &models.Settings{
			ListnerIP:   listerIP,
			ListnerPort: portListner,
			ForwardIP:   forwardIP,
			ForwardPort: portForward,
			Comment:     comment,
		}
		listSettings = append(listSettings, s)
	}
	return listSettings, nil
}

//ReadConfIni подключение, использующее библиотеку gopkg.in/ini.v1
func ReadConfIni() ([]*models.Settings, error) {
	data, err := ioutil.ReadFile("config.ini")
	if err != nil {
		return nil, err
	}

	cfg, err := ini.Load(data)
	if err != nil {
		return nil, err
	}
	fmt.Println(cfg)
	return nil, nil
}

func mainSplitComment(row string) (string, string) {
	arr := strings.Split(row, "#")
	if len(arr) <= 0 {
		return "", ""
	}
	if len(arr) == 1 {
		return arr[0], ""
	}
	if len(arr) == 2 {
		return arr[0], arr[1]
	}
	if len(arr) > 2 {
		for i := 1; i < len(arr); i++ {
			if strings.Trim(arr[i], " ") != "" && strings.Trim(arr[i], " ") != "#" {
				return arr[0], arr[i]
			}
		}
	}
	return "", ""
}
