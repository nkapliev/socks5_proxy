package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/nkapliev/go-socks5"
	"log"
	"fmt"
)

type AppConfig struct {
	ListenTransport      string `default:"tcp",split_words:"true"`
	ListenIp      	     string `default:"",split_words:"true"`
	ListenPort           string `default:"1080",split_words:"true"`
	Login                string `required:"true",split_words:"true"`
	Password             string `required:"true",split_words:"true"`
}

const CONFIG_ENV_PREFIX = "SOCKS5_PROXY"

func get_socks5_config(appConfig *AppConfig) *socks5.Config {
	cred := socks5.StaticCredentials{appConfig.Login: appConfig.Password}
	authenticator := socks5.UserPassAuthenticator{Credentials: cred}
	return &socks5.Config{AuthMethods: []socks5.Authenticator{authenticator}}
}

func main() {
	var appConfig AppConfig
	if err := envconfig.Process(CONFIG_ENV_PREFIX, &appConfig); err != nil {
		log.Fatal(err.Error())
	}

	format := "ListenTransport: %s\nListenIp: %s\nListenPort: %s\nLogin: %s\nPassword: %s\n"
	fmt.Printf(format, appConfig.ListenTransport, appConfig.ListenIp, appConfig.ListenPort, appConfig.Login, appConfig.Password)

	socks5Config := get_socks5_config(&appConfig)
	server, err := socks5.New(socks5Config)
	if err != nil {
		panic(err)
	}

	addr := appConfig.ListenIp + ":" + appConfig.ListenPort
	if err := server.ListenAndServe(appConfig.ListenTransport, addr); err != nil {
		panic(err)
	}
}
