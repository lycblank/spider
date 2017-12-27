package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Mail MailConfig `json:"mail"`
}

type MailConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	To       string `json:"to"`
}

var Conf Config

func Init(runmode string) {
	filename := "conf/dev.json"
	if runmode == "prod" {
		filename = "conf/prod.json"
	}
	datas, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("read config file failed. error:%s", err)
		os.Exit(1)
	}

	if err := json.Unmarshal(datas, &Conf); err != nil {
		fmt.Println("unmarshal json datas failed. error:%s", err)
		os.Exit(1)
	}
}
