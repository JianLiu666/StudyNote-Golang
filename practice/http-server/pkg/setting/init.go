package setting

import (
	"log"
	"time"

	"gopkg.in/ini.v1"
)

var cfg *ini.File

func init() {
	var err error
	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini: %v", err)
	}

	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)
	mapTo("app", AppSetting)

	ServerSetting.ReadTimeout *= time.Second
	ServerSetting.WriteTimeout *= time.Second

}

func mapTo(section string, v any) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("cfg.MapTo %s err: %v", section, err)
	}
}
