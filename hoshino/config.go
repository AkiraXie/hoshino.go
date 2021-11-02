package hoshino

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	NickName       []string `yaml:"nickname"`
	CommandPrefix  string   `yaml:"command_prefix"`
	SuperUsers     []string `yaml:"super_users"`
	ExpireDuration int      `yaml:"expire_duration"`
	AccessToken    string   `yaml:"access_token"`
	Host           string   `yaml:"host"`
	Port           int      `yaml:"port"`
}

var HsoConfig Config

func init() {
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic("config.yaml打开失败或不存在！")
	}
	err = yaml.Unmarshal(data, &HsoConfig)
	if err == nil {
		log.Info("已成功加载配置: ", HsoConfig)
	} else {
		panic(err)
	}

}
