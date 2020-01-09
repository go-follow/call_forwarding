package models

//Settings настройки для адресов и портов для переадресации
type Settings struct {
	ListnerIP   string
	ListnerPort int
	ForwardIP   string
	ForwardPort int
	Comment     string
}
