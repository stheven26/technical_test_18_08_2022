package config

import "github.com/tkanos/gonfig"

type Configuration struct {
	DB_USERNAME string
	DB_PASSWORD string
	DB_PORT     string
	DB_HOST     string
	DB_NAME     string
}

func ConfigDB() Configuration {
	config := Configuration{}
	gonfig.GetConf("./config/config.json", &config)
	return config
}
