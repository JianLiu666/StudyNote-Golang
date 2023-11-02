package setting

import "time"

var ServerSetting = &Server{}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var DatabaseSetting = &Database{}

type Database struct {
	Type     string
	User     string
	Password string
	Host     string
	DBName   string
}

var AppSetting = &Application{}

type Application struct {
	PageSize int
}
