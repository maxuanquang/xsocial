package configs

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func parseConfig(cfgPath string) (*Config, error) {
	// Read the YAML file
	yamlFile, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return &Config{}, err
	}

	// Parse the YAML data into a struct
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return &Config{}, err
	}

	return &config, nil
}

type Config struct {
	MySQL               MySQLConfig               `yaml:"mysql"`
	Redis               RedisConfig               `yaml:"redis"`
	AuthenticateAndPost AuthenticateAndPostConfig `yaml:"authenticate_and_post_config"`
	Newsfeed            NewsfeedConfig            `yaml:"newsfeed_config"`
	WebConfig           WebConfig                 `yaml:"web_config"`
}

type MySQLConfig struct {
	DSN                       string `yaml:"dsn"`
	DefaultStringSize         int    `yaml:"defaultStringSize"`
	DisableDatetimePrecision  bool   `yaml:"disableDatetimePrecision"`
	DontSupportRenameIndex    bool   `yaml:"dontSupportRenameIndex"`
	SkipInitializeWithVersion bool   `yaml:"skipInitializeWithVersion"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
}

type AuthenticateAndPostConfig struct {
	Port  int         `yaml:"port"`
	MySQL MySQLConfig `yaml:"mysql"`
	Redis RedisConfig `yaml:"redis"`
}

type NewsfeedConfig struct {
	Port  int         `yaml:"port"`
	MySQL MySQLConfig `yaml:"mysql"`
	Redis RedisConfig `yaml:"redis"`
}

type WebConfig struct {
	Port                int        `yaml:"port"`
	APIVersions         []string   `yaml:"api_version"`
	AuthenticateAndPost HostConfig `yaml:"authenticate_and_post"`
	Newsfeed            HostConfig `yaml:"newsfeed"`
}

type HostConfig struct {
	Hosts []string `yaml:"hosts"`
}
