package conf

import (
	"github.com/spf13/viper"
	"strings"
)

type Conf struct {
	App      App      `mapstructure:"app"`
	HttpServ HttpServ `mapstructure:"http_serv"`
	Bot      Bot      `mapstructure:"bot"`
}

type App struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
}

type HttpServ struct {
	Addr         string `mapstructure:"addr"`
	Mode         string `mapstructure:"mode"`
	ReadTimeout  string `mapstructure:"read_timeout"`
	WriteTimeout string `mapstructure:"write_timeout"`
}

type Bot struct {
	Token   string `mapstructure:"token"`
	BaseUrl string `mapstructure:"base_url"`
}

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func New(file string) (conf *Conf, err error) {
	conf = &Conf{}
	viper.SetConfigFile(file)
	if err := viper.MergeInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(conf); err != nil {
		return nil, err
	}

	return conf, nil
}
